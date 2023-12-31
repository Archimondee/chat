package interfaces

import (
	"chat/app/models/entity"
	"chat/app/models/request"
)

type RoomRepository interface {
	CreateRoom(room *request.RoomCreateRequest) (*entity.Room, error)
	GetAllRoom() ([]*entity.RoomUser, error)
	FindRoomById(RoomId string) (*entity.RoomUser, error)
	JoinRoom(participant *request.ParticipantCreateRequest) (*entity.Participant, error)
	CheckParticipant(UserId string, RoomId string) (bool, error)
}
