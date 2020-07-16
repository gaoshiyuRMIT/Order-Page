package utils

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"strings"

	/*
		Package pq is a pure Go Postgres driver for the database/sql package.
		In most cases clients will use the database/sql package instead of using this package directly
	*/
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConfigReader config reader
type ConfigReader struct {
	ConfigValue map[string]map[string]string
	PostgresDB *sql.DB
}

// NewConfigReader constructor
func NewConfigReader(filepath string) (*ConfigReader, error) {
	reader := &ConfigReader{}
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Cannot open config file. %w", err)
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot read from config file. %w", err)
	}
	err = json.Unmarshal(bytes, &reader.ConfigValue)
	if err != nil {
		return nil, fmt.Errorf("Parsing JSON from config file failed. %w", err)
	}
	reader.PostgresDB, err = reader.GetPostgresDB()
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// GetPostgresDB gets a postgres db instance
func (cfg ConfigReader) GetPostgresDB() (*sql.DB, error) {
	connStr := cfg.ConfigValue["db_connection"]["postgres"]
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to Postgres database. %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, nil
}

func (cfg ConfigReader) GetMongoDB()  (*mongo.Client, *mongo.Database, error) {
	connStr := cfg.ConfigValue["db_connection"]["mongodb"]
	components := strings.Split(connStr, "/")
	dbName := components[len(components) - 1]
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot connect to MongoDB. %w", err)
	}
	return client, client.Database(dbName), nil
}

func (cfg ConfigReader) GetAPIPort() string {
	portStr := cfg.ConfigValue["api"]["port"]
	return portStr
}