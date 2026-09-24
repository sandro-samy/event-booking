package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sandro-samy/event-booking/utils"
)

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	userID, err := utils.VerifyToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	c.Set("userID", userID)
	c.Next()
}
