package database

import (
	"fmt"
	"time"
	
	"gorm.io/gorm"
	"gorm.io/driver/postgres"

	"selatoz/pkg/cfglib"
)

func NewDatabase() (*gorm.DB, error) {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfglib.DefaultConf.DBHost, cfglib.DefaultConf.DBPort, cfglib.DefaultConf.DBUser, cfglib.DefaultConf.DBPassword, cfglib.DefaultConf.DBName)
	// The commented line is for mysql
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfglib.DefaultConf.DBUser, cfglib.DefaultConf.DBPassword, cfglib.DefaultConf.DBHost, cfglib.DefaultConf.DBPort, cfglib.DefaultConf.DBName)
	
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