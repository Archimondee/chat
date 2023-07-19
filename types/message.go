package types

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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
