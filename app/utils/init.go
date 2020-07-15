package utils

// Config global Configuration
var Config *ConfigReader

func init() {
	Config = NewConfigReader("config.json")
}
