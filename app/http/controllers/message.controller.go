package controllers

import (
	"chat/app/interfaces"
	"chat/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageController struct {
	messageRepository interfaces.MessageRepository
	ctx               context.Context
}

func NewMessageController(ctx context.Context, messageRepository interfaces.MessageRepository) MessageController {
	return MessageController{
		messageRepository: messageRepository,
		ctx:               ctx,
	}
}

func (mc *MessageController) ReadMessage(ctx *gin.Context) {
	sender := ctx.Query("sender")
	recipient := ctx.Query("recipient")
	roomId := ctx.Query("room_id")

	if roomId != "" {
		data, err := mc.messageRepository.ReadRoomMessage(roomId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, utils.ResponseData("success", "success", data))
		return
	} else {
		data, err := mc.messageRepository.ReadMessage(sender, recipient)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, utils.ResponseData("success", "success", data))
		return
	}

}
