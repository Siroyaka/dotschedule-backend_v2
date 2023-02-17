package domain

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type SeedSchedule struct {
	id             string
	platformType   string
	status         string
	participants   []string
	visibleStatus  int
	completeStatus int
	publishedAt    wrappedbasics.WrappedTime
}

func NewSeedSchedule(id, platformType, status string, publishedAt wrappedbasics.WrappedTime) SeedSchedule {
	return SeedSchedule{
		id:             id,
		platformType:   platformType,
		status:         status,
		publishedAt:    publishedAt,
		participants:   []string{},
		visibleStatus:  0,
		completeStatus: 0,
	}
}

func NewSeedScheduleWithParticipants(id, platformType, status string, publishedAt wrappedbasics.WrappedTime, participants []string) SeedSchedule {
	return SeedSchedule{
		id:             id,
		platformType:   platformType,
		status:         status,
		participants:   participants,
		publishedAt:    publishedAt,
		visibleStatus:  0,
		completeStatus: 0,
	}
}

func (s SeedSchedule) GetPublishedAt() wrappedbasics.WrappedTime {
	return s.publishedAt
}

func (s SeedSchedule) GetID() string {
	return s.id
}

func (s SeedSchedule) GetPlatformType() string {
	return s.platformType
}

func (s SeedSchedule) GetStatus() string {
	return s.status
}

func (s SeedSchedule) GetParticipants() []string {
	return s.participants
}

func (s SeedSchedule) GetVisibleStatus() int {
	return s.visibleStatus
}

func (s SeedSchedule) GetCompleteStatus() int {
	return s.completeStatus
}
