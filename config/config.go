package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	AppName 		string
	AppSecret	string
	AppDebug   	bool
	AppPort    	int

	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
}

// Init initializes the application configuration using the provided file name.
func Init(fileName string) (*Config, error) {
	// Initialize viper
	// viper.AddConfigPath(".")
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(filepath.Join(wd, "config", fileName))
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Create configuration struct
	conf := &Config{
		AppName:    viper.GetString("app.name"),
		AppSecret:	viper.GetString("app.secret"),
		AppDebug:   viper.GetBool("app.debug"),
		AppPort:    viper.GetInt("app.port"),
		DBHost:     viper.GetString("database.host"),
		DBPort:     viper.GetInt("database.port"),
		DBName:     viper.GetString("database.name"),
		DBUser:     viper.GetString("database.user"),
		DBPassword: viper.GetString("database.password"),
	}

	return conf, nil
}