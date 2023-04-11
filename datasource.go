package web

import (
	"database/sql"

	"github.com/golang/protobuf/proto"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewDataSourceRepository(db *sql.DB, logger log.Logger) *DataSourceRepository {
	return &DataSourceRepository{logger: logger, db: db}
}

func (repo *DataSourceRepository) Get(id string) (*DataSource, error) {

	datasource := DataSource{DataSourceConfig: make(map[string]string)}
	dataSourceRow := repo.db.QueryRow("SELECT * FROM datasources WHERE id = ?", id)
	err := dataSourceRow.Scan(&datasource.Id, &datasource.Type, &datasource.Event, &datasource.Name)
	if err != nil {
		return nil, err
	}

	dsConfigRows, err := repo.db.Query("SELECT key, value FROM datasource_config WHERE datasource_id = ?", id)
	if err != nil {
		return nil, err
	}

	for dsConfigRows.Next() {
		var key, value string
		err := dsConfigRows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		datasource.DataSourceConfig[key] = value

	}
	return &datasource, nil
}

func (repo *DataSourceRepository) Set(datasource DataSource) (*DataSource, error) {

	if len(datasource.Id) == 0 {
		datasource.Id = utils.NewId()
	}
	_, err := repo.db.Exec("INSERT INTO datasources(id, type, event_type, name) VALUES(?, ?, ?, ?)", datasource.Id, datasource.Type, datasource.Event, datasource.Name)
	return &datasource, err
}

func (repo *DataSourceRepository) Delete(id string) error {

	_, err := repo.db.Exec("DELETE FROM datasources WHERE id = ?", id)
	return err
}

func (repo *DataSourceRepository) List() ([]DataSource, error) {

	datasources := []DataSource{}
	dataSourceRows, err := repo.db.Query("SELECT id FROM datasources")
	if err != nil {
		return datasources, err
	}

	for dataSourceRows.Next() {
		var id string
		err := dataSourceRows.Scan(&id)
		if err != nil {
			return datasources, err
		}
		datasource, err := repo.Get(id)
		if err != nil {
			return datasources, err
		}
		datasources = append(datasources, *datasource)
	}
	return datasources, nil
}

// AsJson converts given protobuf message to its JSON representation.
func asJson(message proto.Message) ([]byte, error) {
	return protojson.Marshal(proto.MessageV2(message))
}
