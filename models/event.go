package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Event struct {
	ID          int
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      uuid.UUID
}

var events = []Event{}

func (e Event) Save() {
	//later add it to database
	events = append(events, e)
}

func GetAllEvents() []Event {
	return events
}
