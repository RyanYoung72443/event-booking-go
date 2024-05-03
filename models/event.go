package models

import (
	"errors"
	"time"

	"example.com/event-booking/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?,?,?,?,?)`
	stmt, err := db.DB.Prepare(query) // Prepare is great for looped and batched requests using the same query
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(event.ID)
	return err
}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?,?)"
	db.DB.Prepare(query)
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e Event) CheckIfUserIsRegistered(userId int64) error {
	query := "SELECT id FROM registrations WHERE user_id = ? AND event_id = ?"
	row := db.DB.QueryRow(query, userId, e.ID)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil
	}
	return errors.New("user registered")
}

func (e Event) CheckRegistrations() ([]string, error) {
	registrationQuery := "SELECT user_id FROM registrations WHERE event_id = ?"
	regRows, err := db.DB.Query(registrationQuery, e.ID)
	if err != nil {
		return nil, err
	}
	defer regRows.Close()

	var registeredUserIds []int64

	for regRows.Next() {
		var user_id int64
		err := regRows.Scan(&user_id)
		if err != nil {
			return nil, err
		}

		registeredUserIds = append(registeredUserIds, user_id)
	}

	usersQuery := "SELECT email FROM users WHERE id = ?"
	stmt, err := db.DB.Prepare(usersQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var registeredUsers []string

	for id := range registeredUserIds {
		var user_email string
		err := stmt.QueryRow(id).Scan(&user_email)
		if err != nil {
			continue
		}

		registeredUsers = append(registeredUsers, user_email)
	}

	return registeredUsers, nil
}
