package models

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Age           int       `json:"age"`
	RecordingDate int64     `json:"recording_date"`
	IsDeleted     bool      `json:"is_deleted"`
}
