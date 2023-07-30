package interfaces

import (
	"chat/app/models/entity"
	"github.com/google/uuid"
)

type Message struct {
	Uuid      uuid.UUID `json:"uuid"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Recipient string    `json:"recipient"`
	RoomId    string    `json:"room_id"`
	//Target  *Room       `json:"target"`
	Sender string `json:"sender"`
	Status string `json:"status"`
}

type MessageRepository interface {
	CreateMessage(message Message) error
	UpdateMessage(message Message) error
	ReadMessage(sender string, recipient string) ([]*entity.Message, error)
}
