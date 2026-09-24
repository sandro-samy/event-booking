package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sandro-samy/event-booking/middleware"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", GetEvents)
	
	// Auth routes
	authGroup := server.Group("/events") 
	authGroup.Use(middleware.Auth)
	authGroup.GET(":id", GetEventByID)
	authGroup.POST("", CreateEvent)
	authGroup.PUT("/:id", UpdateEvent)
	authGroup.DELETE("/:id", DeleteEvent)

	server.POST("/auth/register", Register)
	server.POST("/auth/login", Login)
}