package interfaces

import "chat/app/models/entity"

type ParticipantRepository interface {
	GetAllParticipant() ([]*entity.Participant, error)
}
