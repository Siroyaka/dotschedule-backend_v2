package domain

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type FirestoreNews struct {
	UpdateAt     wrappedbasics.IWrappedTime
	VideoID      string
	Participants []string
}

func NewFirestoreNews(videoID string, updateAt wrappedbasics.IWrappedTime, participants []string) FirestoreNews {
	return FirestoreNews{
		UpdateAt:     updateAt,
		VideoID:      videoID,
		Participants: participants,
	}
}
