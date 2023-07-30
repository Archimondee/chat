package repository

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"context"
	"gorm.io/gorm"
)

type ParticipantRepositoryImpl struct {
	ctx context.Context
	DB  *gorm.DB
}

func NewParticipantRepositoryImpl(ctx context.Context, db *gorm.DB) interfaces.ParticipantRepository {
	return &ParticipantRepositoryImpl{
		DB:  db,
		ctx: ctx,
	}
}

func (p ParticipantRepositoryImpl) GetAllParticipant() ([]*entity.Participant, error) {
	var participants []*entity.Participant
	result := p.DB.Find(&participants)

	if result.Error != nil {
		return nil, result.Error
	}

	return participants, nil
}
