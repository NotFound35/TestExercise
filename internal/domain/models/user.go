package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            string `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Age           int    `json:"age"`
	RecordingDate int64  `json:"recording_date"` // Unix timestamp
}

func NewUser(firstName, lastName string, age int) *User {
	return &User{ // структура с полями возвращается
		ID:            uuid.New().String(), // Генерация UUID
		FirstName:     firstName,
		LastName:      lastName,
		Age:           age,
		RecordingDate: time.Now().Unix(),
	}
}

// создает новый объект в памяти, но вообще не работает с бд
