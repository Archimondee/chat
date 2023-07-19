package repository

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"context"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	ctx context.Context
	DB  *gorm.DB
}

func NewUserRepositoryImpl(ctx context.Context, db *gorm.DB) interfaces.UserRepository {
	return &UserRepositoryImpl{
		DB:  db,
		ctx: ctx,
	}
}

func (u UserRepositoryImpl) FindUserByEmail(email string) (*entity.User, error) {
	var user *entity.User
	result := u.DB.First(&user, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserRepositoryImpl) FindUserById(uuid string) (*entity.User, error) {
	var user *entity.User
	result := u.DB.First(&user, "uuid = ?", uuid)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
