package entity

import (
	"github.com/google/uuid"
	"time"
)

type Participant struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Uuid      uuid.UUID `json:"uuid"`
	UserId    uuid.UUID `json:"user_id"`
	RoomId    uuid.UUID `json:"room_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

func (Participant) TableName() string {
	return "participants"
}
