package web

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	"testing"
)

type DBTestSuite struct {
	db     *sql.DB
	conf   config.Config
	testDb *testDatabase
	suite.Suite
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (suite *DBTestSuite) SetupSuite() {
	suite.testDb = &testDatabase{}
	suite.db, suite.conf = suite.testDb.setup("", suite.Assert())
}

func (suite *DBTestSuite) TearDownSuite() {
	suite.testDb.tearDown(suite.Assert())
}

func (suite *DBTestSuite) TestGenerateDBDriver() {

	config1 := loadConfigForTest("fixtures/testconfig_01.yml")
	db1, err1 := NewDatabase(config1)
	suite.Nil(db1)
	suite.NotNil(err1)

	config2 := loadConfigForTest("fixtures/testconfig_02.yml")
	db2, err2 := NewDatabase(config2)
	suite.NotNil(db2)
	suite.Nil(err2)
	suite.Nil(db2.Ping())
	suite.Nil(db2.Close())
}

func (suite *DBTestSuite) TestMigrations() {

	upErr := SetupDatabaseSchema(suite.db, suite.conf)
	suite.Nil(upErr)
	logError(upErr)

	downErr := TearDownDatabaseSchema(suite.db, suite.conf)
	suite.Nil(downErr)
	logError(downErr)
}
