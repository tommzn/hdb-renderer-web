package web

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	config "github.com/tommzn/go-config"
)

// NewDatabase creates a new database driver from passed config.
// Atm sqlite3 with local file is supported, only.
func NewDatabase(conf config.Config) (*sql.DB, error) {

	dbType := conf.Get("db.type", nil)
	if dbType == nil {
		return nil, errors.New("No database defined.")
	}
	switch *dbType {
	case "sqlite3":
		dbFile := conf.Get("db.sqlite3.file", nil)
		if dbFile == nil {
			return nil, errors.New("No sqlite3 file defined.")
		}
		return sql.Open("sqlite3", *dbFile)

	default:
		return nil, fmt.Errorf("Invaliddatabase: %s\n", *dbType)
	}
}

// SetupDatabaseSChema ensures required data structures.
func SetupDatabaseSchema(db *sql.DB, conf config.Config) error {

	m, err := newDbMigration(db, conf)
	if err != nil {
		return err
	}

	return m.Up()
}

// TearDownDatabaseSchema deletes alls data structures.
func TearDownDatabaseSchema(db *sql.DB, conf config.Config) error {

	m, err := newDbMigration(db, conf)
	if err != nil {
		return err
	}

	return m.Down()
}

// NewDbMigration create a new DB migrate instnace.
func newDbMigration(db *sql.DB, conf config.Config) (*migrate.Migrate, error) {

	dbType := conf.Get("db.type", nil)
	if dbType == nil {
		return nil, errors.New("No database defined.")
	}

	sourceDriver, sourceType, err := newMigrationSource(conf)
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance(*sourceType, *sourceDriver, *dbType, driver)
}

// NewMigrationSource create a new source of migration scripts from given config.
// Atm a file source is supported, only.
func newMigrationSource(conf config.Config) (*source.Driver, *string, error) {

	sourceType := conf.Get("db.migration.source_type", nil)
	if sourceType == nil {
		return nil, nil, errors.New("No migrations source type defined.")
	}

	sourceUrl := conf.Get("db.migration.url", nil)
	if sourceUrl == nil {
		return nil, nil, errors.New("No migrations source url defined.")
	}

	switch *sourceType {
	case "file":
		fileSource, err := (&file.File{}).Open(*sourceUrl)
		if err != nil {
			return nil, nil, err
		}
		return &fileSource, sourceType, nil
	default:
		return nil, nil, errors.New("Unsupported migrations source type: " + *sourceType)
	}
}
