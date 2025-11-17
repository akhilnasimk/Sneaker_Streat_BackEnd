package main

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/migrations"
	"github.com/akhilnasimk/SS_backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()         // concecting db
	migrations.RunMigrations() // running the automigrations

	//setting up the server
	baseRoute := gin.Default()
	baseRoute.Use(middlewares.CORSMiddleware())
	baseRoute.Use(middlewares.RateLimitMiddleware())

	routes.SetupRoutes(baseRoute)

	baseRoute.Run(":8080")
}
