package models

import "time"

type Event struct {
	ID          int
	Name        string
	Description string
	Location    string
	DateTime    time.Time
	UserID      int
}

var Events = []Event{}

func GetAllEvents() []Event {
	return Events
}