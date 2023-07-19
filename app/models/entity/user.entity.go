package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uint      `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}
