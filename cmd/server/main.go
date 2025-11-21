package main

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/migrations"
	"github.com/akhilnasimk/SS_backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()         // concecting db
	migrations.RunMigrations() // running the automigrations
	config.InitCloudinary()    //cloudinery initialization 

	//setting up the server
	baseRoute := gin.Default()

	routes.SetupRoutes(baseRoute)

	baseRoute.Run(":8080")
}
