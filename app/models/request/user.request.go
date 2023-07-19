package request

import (
	"github.com/google/uuid"
)

type UserSigninRequest struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}

type UserCreateRequest struct {
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

func (UserCreateRequest) TableName() string {
	return "users"
}
