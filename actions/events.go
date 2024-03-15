package actions

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/daalvand/go-api/models"
	"github.com/daalvand/go-api/repositories"
	"github.com/gin-gonic/gin"
)

type EventActions struct {
	repo repositories.EventRepository
}

func NewEventActions(repo repositories.EventRepository) *EventActions {
	return &EventActions{repo: repo}
}

func (ea *EventActions) GetEvents(context *gin.Context) {
	events, err := ea.repo.FindAll()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Try again later."})
		return
	}
	context.JSON(http.StatusOK, events)
}

func (ea *EventActions) GetEvent(context *gin.Context) {
	eventID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := ea.repo.FindByID(uint(eventID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		logrus.Error("Error on fetching Event", err)
		return
	}

	context.JSON(http.StatusOK, event)
}

func (ea *EventActions) CreateEvent(context *gin.Context) {
	var event models.Event
	if err := context.ShouldBindJSON(&event); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	userID := context.GetUint("userId")
	event.UserID = userID

	if err := ea.repo.Create(&event); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event. Try again later."})
		logrus.Error("Error on creating Event", err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

func (ea *EventActions) UpdateEvent(context *gin.Context) {
	eventID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	userID := context.GetUint("userId")

	event, err := ea.repo.FindByID(uint(eventID))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event."})
		return
	}

	var updatedEvent models.Event
	if err := context.ShouldBindJSON(&updatedEvent); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	updatedEvent.ID = uint(eventID)
	if err := ea.repo.Update(&updatedEvent); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		logrus.Error("Error on updating Event", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!"})
}

func (ea *EventActions) DeleteEvent(context *gin.Context) {
	eventID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	userID := context.GetUint("userId")
	event, err := ea.repo.FindByID(uint(eventID))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete event."})
		return
	}

	if err := ea.repo.Delete(event); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event."})
		logrus.Error("Error on deleting Event", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}

func (ea *EventActions) RegisterForEvent(context *gin.Context) {
	userId := context.GetUint("userId")
	eventId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := ea.repo.FindByID(uint(eventId))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		logrus.Error("Error on fetching Event", err)
		return
	}

	err = ea.repo.Register(event, userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		logrus.Error("Error on registering Event", err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registered!"})
}

func (ea *EventActions) CancelRegistration(context *gin.Context) {
	userId := context.GetUint("userId")
	eventId, err := strconv.ParseUint(context.Param("id"), 10, 64)

	var event models.Event
	event.ID = uint(eventId)

	err = ea.repo.CancelRegistration(&event, userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration."})
		logrus.Error("Error on registering Event", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}
