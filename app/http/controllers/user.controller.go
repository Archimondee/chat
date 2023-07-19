package controllers

import (
	"chat/app/interfaces"
	"chat/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userRepository interfaces.UserRepository
	ctx            context.Context
}

func NewUserController(userRepository interfaces.UserRepository, ctx context.Context) UserController {
	return UserController{
		userRepository: userRepository,
		ctx:            ctx,
	}
}

func (uc *UserController) GetAllUser(ctx *gin.Context) {
	users, err := uc.userRepository.FindUserAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "error", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "login success", users))

	return
}
