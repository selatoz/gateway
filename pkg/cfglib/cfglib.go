package cfglib

import (
	"os"
	"fmt"
	"path/filepath"
	"strconv"

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
		AppPort:          strToInt(os.Getenv("APP_PORT")),
		GinMode:				os.Getenv("GIN_MODE"),

		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           strToInt(os.Getenv("DB_PORT")),
		DBName:           os.Getenv("DB_NAME"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),

		TokenExpAccess:  strToFloat32(os.Getenv("TOKEN_EXP_ACCESS")),
		TokenExpRefresh: strToFloat32(os.Getenv("TOKEN_EXP_REFRESH")),
  	}

	// Set the app mode
	if DefaultConf.GinMode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// Handle return
	return nil
}

func strToInt(str string) (int) {
	n, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		panic(fmt.Errorf("could not convert string '%s' to int: %s", str, err))
	}

	return int(n)
}

func strToFloat32(str string) (float32) {
	n, err := strconv.ParseFloat(str, 32)
	if err != nil {
		panic(fmt.Errorf("could not convert string '%s' to float64: %s", str, err))
	}

	return float32(n)
}