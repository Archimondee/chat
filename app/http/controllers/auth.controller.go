package controllers

import (
	"chat/app/interfaces"
	"chat/app/models/request"
	"chat/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	authRepository interfaces.AuthRepository
	ctx            context.Context
}

func NewAuthController(authRepository interfaces.AuthRepository, ctx context.Context) AuthController {
	return AuthController{authRepository, ctx}
}

func (ac *AuthController) SigninUser(ctx *gin.Context) {
	var user request.UserSigninRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "something error", err))
		return
	}

	token, err := ac.authRepository.SigninUser(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "wrong email or password", nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "login success", gin.H{
		"token": token,
	}))
}

func (ac *AuthController) SignupUser(ctx *gin.Context) {
	var user request.UserCreateRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "something error", err))
	}

	newUser, err := ac.authRepository.SignupUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "error", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "register success", newUser))
	return
}
