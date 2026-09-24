package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sandro-samy/event-booking/models"
	"github.com/sandro-samy/event-booking/utils"
)

func Register(c *gin.Context) {

	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "valid email and password are required",
			"error":   err,
		})

		return
	}

	if err := newUser.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something Went Wrong",
			"error":   err,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registered Successfully",
	})
}

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "valid email and password are required",
			"error":   err,
		})

		return
	}

	if err := user.ValidateCredentials(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "valid email and password are required",
			"error":   err,
		})
		
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successfully",
		"token":   token,
	})
}
