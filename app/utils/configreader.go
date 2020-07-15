package utils

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	/*
		Package pq is a pure Go Postgres driver for the database/sql package.
		In most cases clients will use the database/sql package instead of using this package directly
	*/
	_ "github.com/lib/pq"
)

// ConfigReader config reader
type ConfigReader struct {
	ConfigValue map[string]map[string]string
}

// NewConfigReader constructor
func NewConfigReader(filepath string) *ConfigReader {
	reader := &ConfigReader{}
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &reader.ConfigValue)
	if err != nil {
		log.Fatal(err)
	}
	return reader
}

// GetPostgresDB gets a postgres db instance
func (cfg ConfigReader) GetPostgresDB() *sql.DB {
	connStr := cfg.ConfigValue["db_connection"]["postgres"]
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
