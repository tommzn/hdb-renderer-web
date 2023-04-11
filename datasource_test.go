package web

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	utils "github.com/tommzn/go-utils"
	events "github.com/tommzn/hdb-events-go"
	"testing"
	"time"
)

type DatasourceTestSuite struct {
	db     *sql.DB
	testDb *testDatabase
	suite.Suite
}

func TestDatasourceTestSuite(t *testing.T) {
	suite.Run(t, new(DatasourceTestSuite))
}

type DatasourceUtilsTestSuite struct {
	suite.Suite
}

func TestDatasourceUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(DatasourceUtilsTestSuite))
}

func (suite *DatasourceTestSuite) SetupSuite() {
	suite.testDb = &testDatabase{}
	suite.db, _ = suite.testDb.setup("", suite.Assert())
	suite.testDb.up(suite.Assert())
}

func (suite *DatasourceTestSuite) TearDownSuite() {
	suite.testDb.down(suite.Assert())
	suite.testDb.tearDown(suite.Assert())
}
func (suite *DatasourceTestSuite) datasourceRepositoryForTest() *DataSourceRepository {
	return NewDataSourceRepository(suite.db, loggerForTest())
}

func (suite *DatasourceTestSuite) TestCRUDActions() {

	repo := suite.datasourceRepositoryForTest()

	datasource1 := DataSource{Type: DATASOURCE_KAFKA, Event: EVENTTYPE_WEATHER, Name: "Weather DS 1", DataSourceConfig: make(map[string]string)}
	datasource1_1, err1_1 := repo.Set(datasource1)
	suite.Nil(err1_1)
	suite.True(utils.IsId(datasource1_1.Id))

	datasource2 := DataSource{Type: DATASOURCE_KAFKA, Event: EVENTTYPE_WEATHER, Name: "Weather DS 2", DataSourceConfig: make(map[string]string)}
	datasource2_1, err2_1 := repo.Set(datasource2)
	suite.Nil(err2_1)
	suite.True(utils.IsId(datasource2_1.Id))

	suite.NotEqual(datasource1_1.Id, datasource2_1.Id)

	datasource1_2, err1_2 := repo.Get(datasource1_1.Id)
	suite.Nil(err1_2)
	suite.Equal(datasource1_1, datasource1_2)

	datasources, err := repo.List()
	suite.Nil(err)
	suite.Len(datasources, 2)

	suite.Nil(repo.Delete(datasource2_1.Id))

	datasources2, err2 := repo.List()
	suite.Nil(err2)
	suite.Len(datasources2, 1)
}

func (suite *DatasourceUtilsTestSuite) TestConvertToJson() {

	message := exchangeRateForTest()
	json, err := asJson(message)
	suite.Nil(err)
	fmt.Printf("J: %s\n", json)
}

func exchangeRateForTest() *events.ExchangeRates {
	return &events.ExchangeRates{
		Rates: []*events.ExchangeRate{
			&events.ExchangeRate{
				FromCurrency: "USD",
				ToCurrency:   "EUR",
				Rate:         1.23445,
				Timestamp:    asTimeStamp(time.Now()),
			},
		},
	}
}
