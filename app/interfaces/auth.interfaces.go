package interfaces

import (
	"chat/app/models/entity"
	"chat/app/models/request"
)

type AuthRepository interface {
	SignupUser(input *request.UserCreateRequest) (*entity.User, error)
	SigninUser(input *request.UserSigninRequest) (*string, error)
}
