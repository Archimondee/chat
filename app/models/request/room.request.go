package request

import "github.com/google/uuid"

type RoomCreateRequest struct {
	Uuid uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

func (RoomCreateRequest) TableName() string {
	return "rooms"
}
