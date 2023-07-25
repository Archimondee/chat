package request

import "github.com/google/uuid"

type ParticipantCreateRequest struct {
	Uuid   uuid.UUID `json:"uuid"`
	UserId string    `json:"user_id"`
	RoomId string    `json:"room_id"`
}

func (ParticipantCreateRequest) TableName() string {
	return "participants"
}
