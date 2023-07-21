package entity

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id        uint      `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	RoomId    uuid.UUID `json:"room_id"`
	Sender    uuid.UUID `json:"sender"`
	Recipient uuid.UUID `json:"participant"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}
