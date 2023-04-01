package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	// "github.com/selatoz/gateway/pkg/cfglib"
	// "github.com/selatoz/gateway/pkg/dblib"
	
	"selatoz/pkg/cfglib"
	"selatoz/pkg/dblib"
	"selatoz/routes"
)

// var db = make(map[string]string)

func main() {
	// Initialize configuration
	err := cfglib.Load();
	if err != nil {
		panic(fmt.Errorf("failed to initialize configuration: %w", err))
	}

	// Initialize the database
	db, err := dblib.NewDatabase(cfglib.DefaultConf)
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %w", err))
	}

	// Initialize the routes
	router := gin.Default()
	routes.NewRoutes(router, db)

	// Listen and Server in 0.0.0.0:8080
	router.Run(fmt.Sprintf(":%d", cfglib.DefaultConf.AppPort))
}
