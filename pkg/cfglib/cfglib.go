package cfglib

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Config represents the structure of the configuration.
type Config struct {
	AppName 		string
	AppSecret	string
	AppDebug   	bool
	AppPort    	int

	GinMode		string

	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	TokenExpAccess 	float32
	TokenExpRefresh 	float32
}

// Variable to store the default config, can be imported and used in other packages
var DefaultConf *Config

func Load() (error) {
	// Build the file path
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = godotenv.Load(filepath.Join(wd, "config", ".env"))
	if err != nil {
		return err
	}

	// Create configuration structures
	DefaultConf = &Config{
		AppName:          os.Getenv("APP_NAME"),
		AppSecret:        os.Getenv("APP_SECRET"),
		AppDebug:         os.Getenv("APP_DEBUG") == "true",
		AppPort:          8080,
		GinMode:				os.Getenv("GIN_MODE"),

		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           5432,
		DBName:           os.Getenv("DB_NAME"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),

		TokenExpAccess:  0.01, 	// In hours
		TokenExpRefresh: 168,	// In hours
  	}

	// Set the app mode
	if DefaultConf.GinMode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// Handle return
	return nil
}