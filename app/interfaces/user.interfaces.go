package interfaces

import "chat/app/models/entity"

type UserRepository interface {
	FindUserByEmail(string) (*entity.User, error)
	FindUserById(string) (*entity.User, error)
}
