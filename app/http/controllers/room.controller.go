package controllers

import (
	"chat/app/interfaces"
	"chat/app/models/request"
	"chat/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoomController struct {
	roomRepository interfaces.RoomRepository
	ctx            context.Context
}

func NewRoomController(roomRepository interfaces.RoomRepository, ctx context.Context) RoomController {
	return RoomController{
		roomRepository: roomRepository,
		ctx:            ctx,
	}
}

func (rc *RoomController) CreateRoom(ctx *gin.Context) {
	var room request.RoomCreateRequest

	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "something error", err))
		return
	}

	createRoom, err := rc.roomRepository.CreateRoom(&room)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "error", nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success", createRoom))
	return
}

func (rc *RoomController) GetAllRoom(ctx *gin.Context) {
	rooms, err := rc.roomRepository.GetAllRoom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "error", err))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success", rooms))

	return
}

func (rc *RoomController) JoinRoom(ctx *gin.Context) {
	var participant request.ParticipantCreateRequest

	if err := ctx.ShouldBindJSON(&participant); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "something error", err))
		return
	}

	createParticipant, err := rc.roomRepository.JoinRoom(&participant)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success", createParticipant))
	return
}

func (rc *RoomController) CheckRoom(ctx *gin.Context) {
	var participant request.ParticipantCreateRequest

	if err := ctx.ShouldBindJSON(&participant); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "something error", err))
		return
	}

	joined, err := rc.roomRepository.CheckParticipant(participant.UserId, participant.RoomId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success", gin.H{"joined": joined}))
	return
}
