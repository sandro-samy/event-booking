package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sandro-samy/event-booking/middleware"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", GetEvents)
	server.GET("/events/:id", GetEventByID)

	// Auth routes
	authenticated := server.Group("/events")

	authenticated.Use(middleware.Auth)
	authenticated.POST("", CreateEvent)
	authenticated.PUT("/:id", UpdateEvent)
	authenticated.DELETE("/:id", DeleteEvent)

	server.POST("/auth/register", Register)
	server.POST("/auth/login", Login)
}