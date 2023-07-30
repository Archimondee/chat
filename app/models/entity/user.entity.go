package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}
