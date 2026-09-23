package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	users "github.com/sandro-samy/event-booking/models"
)

func Register(c *gin.Context) {

	var newUser users.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "valid email and password are required",
			"error":   err,
		})
	}

	if err := newUser.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something Went Wrong",
			"error":   err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registered Successfully",
	})
}
