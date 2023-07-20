package repository

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"context"
	"fmt"
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

func (u UserRepositoryImpl) FindUserByUuid(uuid string) (*entity.User, error) {
	var user *entity.User
	result := u.DB.First(&user, "uuid = ?", uuid)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserRepositoryImpl) FindUserById(id uint) (*entity.User, error) {
	var user *entity.User
	result := u.DB.First(&user, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UserRepositoryImpl) FindUserAll() ([]*entity.User, error) {
	var user []*entity.User
	result := u.DB.Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

func (u UserRepositoryImpl) UpdateStatus(status string, uuid string) {
	user, err := u.FindUserByUuid(uuid)
	if err != nil {
		fmt.Println("Error", err)
	}

	user.Status = status
	u.DB.Save(&user)
}
