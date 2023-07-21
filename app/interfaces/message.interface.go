package interfaces

import (
	"github.com/google/uuid"
)

type Message struct {
	Uuid      uuid.UUID `json:"uuid"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Recipient string    `json:"recipient"`
	//Target  *Room       `json:"target"`
	Sender string `json:"sender"`
	Status string `json:"status"`
}

type MessageRepository interface {
	CreateMessage(message Message) error
}
