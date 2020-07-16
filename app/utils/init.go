package utils

import "fmt"


// Config global Configuration
var Config *ConfigReader

func init() {
	var err error
	Config, err = NewConfigReader("../config.json")
	if err != nil {
		panic(fmt.Errorf("Package initialization failed. %w", err))
	}
}
