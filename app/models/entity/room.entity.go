package entity

import (
	"github.com/google/uuid"
	"time"
)

type Room struct {
	Id        uint      `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type RoomUser struct {
	Id           uint          `json:"id" gorm:"primaryKey"`
	Uuid         uuid.UUID     `json:"uuid"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	CreatedAt    time.Time     `json:"-"`
	UpdatedAt    time.Time     `json:"-"`
	DeletedAt    time.Time     `json:"-"`
	Participants []Participant `json:"participants" gorm:"foreignKey:RoomId;references:Uuid"`
}

func (RoomUser) TableName() string {
	return "rooms"
}
