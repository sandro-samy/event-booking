package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sandro-samy/event-booking/models"
)

func GetEvents(c *gin.Context) {
	events, err := models.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch events.",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

func GetEventByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "valid id is required",
			"error":   err,
		})
		return
	}

	userID := c.GetInt64("userID")
	event, err := models.GetEventByID(id)

	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "event with id " + c.Param("id") + " not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to get event with id " + c.Param("id"),
			"error":   err,
		})
		return
	}

	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "forbidden access",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": event,
	})
}

func CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "could not parse request data.",
			"error":   err,
		})
		return
	}

	userID := c.GetInt64("userID")
	event.UserID = userID // Replace with the actual user ID from the authenticated user

	if err := event.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create event.",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "event created.",
		"event":   event,
	})
}

func UpdateEvent(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please provide a valid ID",
			"error":   err,
		})
		return
	}

	_, err = models.GetEventByID(eventID)

	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "event with id " + c.Param("id") + " not found",
		})
		return
	}

	var eventUpdates models.Event
	if err := c.ShouldBindJSON(&eventUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "could not parse request data.",
			"error":   err,
		})
		return
	}


	eventUpdates.ID = eventID
	updatedEvent, err := eventUpdates.Update()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update Event with id " + c.Param("id"),
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "event with id " + c.Param("id") + " updated Successfully",
		"event":   updatedEvent,
	})
}

func DeleteEvent(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "valid event is required!",
			"error":   err,
		})
		return
	}

	_, err = models.GetEventByID(eventID)

	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event with id " + c.Param("id") + " not found",
		})
		return
	}

	err = models.Delete(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete event with id " + c.Param("id"),
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "event with id " + c.Param("id") + " deleted successfully",
	})
}
