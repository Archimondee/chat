package repository

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"chat/app/models/request"
	"chat/config"
	"chat/utils"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type AuthRepositoryImpl struct {
	DB             *gorm.DB
	ctx            context.Context
	userRepository interfaces.UserRepository
}

func NewAuthRepository(ctx context.Context, db *gorm.DB, userRepository interfaces.UserRepository) interfaces.AuthRepository {
	return &AuthRepositoryImpl{
		ctx:            ctx,
		DB:             db,
		userRepository: userRepository,
	}
}

func (a AuthRepositoryImpl) SignupUser(user *request.UserCreateRequest) (*entity.User, error) {
	user.Email = strings.ToLower(user.Email)
	user.Name = user.Name
	user.Uuid = uuid.New()
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	result := a.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	var newUser *entity.User
	res := a.DB.First(&newUser, "email = ?", user.Email)
	if res.Error != nil {
		return nil, result.Error
	}

	return newUser, nil
}

func (a AuthRepositoryImpl) SigninUser(user *request.UserSigninRequest) (*entity.User, *string, error) {

	userSignIn, err := a.userRepository.FindUserByEmail(user.Email)

	if err != nil {
		return nil, nil, err
	}
	errPassword := utils.VerifyPassword(userSignIn.Password, user.Password)
	if errPassword != nil {
		return nil, nil, errPassword
	}

	config, _ := config.LoadConfig(".")
	access_token, err := utils.CreateToken(userSignIn, config.AccessTokenPrivateKey)

	return userSignIn, &access_token, nil
}
