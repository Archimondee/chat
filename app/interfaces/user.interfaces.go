package interfaces

import "chat/app/models/entity"

type UserRepository interface {
	FindUserByEmail(string) (*entity.User, error)
	FindUserById(uint) (*entity.User, error)
	FindUserByUuid(string) (*entity.User, error)
	FindUserAll() ([]*entity.User, error)
	UpdateStatus(status string, uuid string)
}
