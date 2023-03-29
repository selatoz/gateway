package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"selatoz/config"
	"selatoz/database"
	"selatoz/routes"
)

// var db = make(map[string]string)

func main() {
	// Initialize configuration
	conf, err := config.Init("default.yaml");
	if err != nil {
		panic(fmt.Errorf("failed to initialize configuration: %w", err))
	}

	// Initialize the database
	db, err := database.NewDatabase(conf)
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %w", err))
	}

	// Initialize the routes
	router := gin.Default()
	routes.Init(router, db)

	// Listen and Server in 0.0.0.0:8080
	router.Run(fmt.Sprintf(":%d", conf.AppPort))
}
