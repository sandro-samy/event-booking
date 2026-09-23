package main

import (
	"github.com/gin-gonic/gin"
	DB "github.com/sandro-samy/event-booking/db"
	"github.com/sandro-samy/event-booking/routes"
)

func main() {
	DB.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)



	server.Run(":8080")
}


