package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	events "github.com/sandro-samy/event-booking/models"
)

func GetEvents(c *gin.Context) {
	events, err := events.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events.",
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

func GetEventByID(c *gin.Context) {
	event, err := events.GetEventById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to get event with id " + c.Param("id"),
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"event": event,
	})
}

func CreateEvent(c *gin.Context) {
	var event events.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
			"error": err,
		})
	}

	event.UserID = "1" // Replace with the actual user ID from the authenticated user

	if err := event.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create event.",
			"error": err,
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created.",
		"event":   event,
	})
}

func UpdateEvent(c *gin.Context) {
	var event events.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
			"error": err,
		})
	}
	event.ID = c.Param("id")
	updatedEvent, err := event.Update()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update Event with id " + c.Param("id"),
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event with id " + c.Param("id") + " Updated Successfully",
		"event":   updatedEvent,
	})
}

func DeleteEvent(c *gin.Context) {
	err := events.Delete(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete Event with id " + c.Param("id"),
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event with id " + c.Param("id") + " Deleted Successfully",
	})
}
