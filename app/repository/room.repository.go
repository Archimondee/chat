package repository

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"chat/app/models/request"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepositoryImpl struct {
	ctx context.Context
	DB  *gorm.DB
}

func NewRoomRepositoryImpl(ctx context.Context, db *gorm.DB) interfaces.RoomRepository {
	return &RoomRepositoryImpl{
		DB:  db,
		ctx: ctx,
	}
}

func (r RoomRepositoryImpl) CreateRoom(room *request.RoomCreateRequest) (*entity.Room, error) {
	uuid := uuid.New()
	room.Uuid = uuid
	result := r.DB.Table("rooms").Create(&room)
	if result.Error != nil {
		return nil, result.Error
	}

	var newRoom *entity.Room
	res := r.DB.Table("rooms").First(&newRoom, "uuid = ?", uuid)
	if res.Error != nil {
		return nil, res.Error
	}

	return newRoom, nil
}

func (r RoomRepositoryImpl) GetAllRoom() ([]*entity.RoomUser, error) {
	var rooms []*entity.RoomUser
	result := r.DB.Preload("Participants").Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}

	return rooms, nil
}

func (r RoomRepositoryImpl) JoinRoom(participant *request.ParticipantCreateRequest) (*entity.Participant, error) {
	ok, err := r.CheckParticipant(participant.UserId, participant.RoomId)
	if err != nil {
		fmt.Println(err)
	}

	if !ok {
		id := uuid.New()
		//participant.Uuid = id
		userId, err := uuid.Parse(participant.UserId)
		if err != nil {
			fmt.Println("not uuid")
		}
		roomId, err := uuid.Parse(participant.RoomId)
		if err != nil {
			fmt.Println("not uuid")
		}
		data := &entity.Participant{
			Uuid:   id,
			UserId: userId,
			RoomId: roomId,
		}
		result := r.DB.Table("participants").Create(data)
		if result.Error != nil {
			return nil, result.Error
		}

		var newRoom *entity.Participant
		res := r.DB.Table("participants").First(&newRoom, "uuid = ?", id)
		if res.Error != nil {
			return nil, res.Error
		}
		return newRoom, nil
	}
	return nil, errors.New("joined")
}

func (r RoomRepositoryImpl) CheckParticipant(UserId string, RoomId string) (bool, error) {
	var data *entity.Participant
	result := r.DB.Table("participants").First(&data, "user_id = ? AND room_id = ?", UserId, RoomId)
	if result.Error != nil {
		return false, nil
	}
	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (r RoomRepositoryImpl) FindRoomById(RoomId string) (*entity.RoomUser, error) {
	var room *entity.RoomUser
	result := r.DB.Preload("Participants").Find(&room, "uuid = ? ", RoomId)
	if result.Error != nil {
		return nil, result.Error
	}

	return room, nil
}
