package entity

import (
	"github.com/google/uuid"
	"time"
)

type Participant struct {
	Id        uint      `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	RoomId    uuid.UUID `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
