package entity

import (
	"github.com/google/uuid"
	"time"
)

type Room struct {
	Id        uint      `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
