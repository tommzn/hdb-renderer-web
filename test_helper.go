package web

import (
	"database/sql"
	"fmt"
	syslog "log"
	"os"
	"time"

	config "github.com/tommzn/go-config"

	"github.com/stretchr/testify/assert"
	log "github.com/tommzn/go-log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loadConfigForTest(fileName string) config.Config {

	configLoader := config.NewFileConfigSource(&fileName)
	config, _ := configLoader.Load()
	return config
}

func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func logValue(v interface{}) {
	fmt.Println(v)
}

func logError(err error) {
	if err != nil {
		syslog.Println(err.Error())
	}
}

func asTimeStamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

type testDatabase struct {
	db   *sql.DB
	conf config.Config
}

func (testDB *testDatabase) setup(configFile string, assert *assert.Assertions) (*sql.DB, config.Config) {

	if len(configFile) == 0 {
		configFile = "fixtures/testdbconfig.yml"
	}

	testDB.conf = loadConfigForTest(configFile)
	db, err := NewDatabase(testDB.conf)
	assert.NotNil(db)
	assert.Nil(err)
	testDB.db = db
	return testDB.db, testDB.conf
}

func (testDB *testDatabase) tearDown(assert *assert.Assertions) {

	if testDB.conf == nil {
		return
	}

	dbType := testDB.conf.Get("db.type", nil)
	if dbType == nil || *dbType != "sqlite3" {
		return
	}

	dbFile := testDB.conf.Get("db.sqlite3.file", nil)
	if dbFile == nil {
		return
	}
	assert.Nil(os.Remove(*dbFile))

	if testDB.db != nil {

		assert.Nil(testDB.db.Close())
	}
}

func (testDB *testDatabase) up(assert *assert.Assertions) {
	upErr := SetupDatabaseSchema(testDB.db, testDB.conf)
	assert.Nil(upErr)
	logError(upErr)

}

func (testDB *testDatabase) down(assert *assert.Assertions) {
	downErr := TearDownDatabaseSchema(testDB.db, testDB.conf)
	assert.Nil(downErr)
	logError(downErr)
}
