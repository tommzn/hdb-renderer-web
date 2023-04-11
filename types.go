package web

import (
	"database/sql"
	"net/http"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-renderer-core"
)

// DataSourceType defines rechnical source of events, e.g. database, message broker, etc.
type DataSourceType string

// EventType defines events provided by a datasource, e.g. weather or indoor climate data
type EventType string

const (
	// DATASOURCE_KAFKA. event conumed from topics in Kafka.
	DATASOURCE_KAFKA DataSourceType = "kafka"
)

const (

	// EVENTTYPE_WEATHER, weather data, current weather and hourly/daily forecasts
	EVENTTYPE_WEATHER EventType = "weather"

	// EVENTTYPE_AWSBILLING, billing report form AWS for current expenses.
	EVENTTYPE_AWSBILLING EventType = "awsbilling"

	// EVENTTYPE_INDOORCLIMATE. climate data (temperature/humidity) collected in different rooms.
	EVENTTYPE_INDOORCLIMATE EventType = "indoorclimate"
)

type DataSource struct {
	Id               string
	Type             DataSourceType
	Event            EventType
	Name             string
	DataSourceConfig map[string]string
}

type ConfigRenderer struct {
	logger log.Logger
	db     *sql.DB
}

type BashbaorcRenderer struct {
	logger      log.Logger
	db          *sql.DB
	datasources map[EventType]core.DataSource
}

type DataSourceRepository struct {
	logger log.Logger
	db     *sql.DB
}

type webServer struct {
	port       string
	conf       config.Config
	logger     log.Logger
	httpServer *http.Server
}

type osFlags struct {
	configFile   string
	configSource string
	k8s          bool
}
