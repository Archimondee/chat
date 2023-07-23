package repository

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepositoryImpl struct {
	ctx context.Context
	DB  *gorm.DB
}

func NewMessageRepositoryImpl(ctx context.Context, db *gorm.DB) interfaces.MessageRepository {
	return &MessageRepositoryImpl{
		DB:  db,
		ctx: ctx,
	}
}

func (m MessageRepositoryImpl) CreateMessage(message interfaces.Message) error {
	sender, _ := uuid.Parse(message.Sender)
	recipient, _ := uuid.Parse(message.Recipient)

	data := &entity.Message{
		Uuid:      message.Uuid,
		Sender:    sender,
		Recipient: recipient,
		Text:      message.Message,
		Status:    message.Status,
	}
	result := m.DB.Table("messages").Create(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m MessageRepositoryImpl) UpdateMessage(message interfaces.Message) error {
	sender, _ := uuid.Parse(message.Sender)
	recipient, _ := uuid.Parse(message.Recipient)

	data := &entity.Message{
		Uuid:      message.Uuid,
		Sender:    sender,
		Recipient: recipient,
		Text:      message.Message,
		Status:    message.Status,
	}
	result := m.DB.Table("messages").Where("uuid = ?", data.Uuid).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m MessageRepositoryImpl) ReadMessage(sender string, recipient string) ([]*entity.Message, error) {
	var data []*entity.Message

	result := m.DB.Table("messages").Find(&data, "sender = ? and recipient = ? or sender = ? and recipient = ?", sender, recipient, recipient, sender)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}
