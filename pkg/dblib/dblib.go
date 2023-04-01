package dblib

import (
	"fmt"
	"time"
	
	"gorm.io/gorm"
	"gorm.io/driver/postgres"

	"github.com/selatoz/gateway/pkg/cfglib"
)

func NewDatabase(config *cfglib.Config) (*gorm.DB, error) {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	// The commented line is for mysql
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DefaultConf.DBUser, config.DefaultConf.DBPassword, config.DefaultConf.DBHost, config.DefaultConf.DBPort, config.DefaultConf.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Set up connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}