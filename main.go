package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sandro-samy/event-booking/db"
	"github.com/sandro-samy/event-booking/routes"
	"github.com/sandro-samy/event-booking/utils"
)

func main() {
	_ = godotenv.Load()

	if err := utils.InitJWT(); err != nil {
		log.Fatal(err)
	}

	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
