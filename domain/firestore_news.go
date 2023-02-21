package domain

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type FirestoreNews struct {
	UpdateAt     wrappedbasics.IWrappedTime
	VideoID      string
	VideoStatus  int
	Participants []string
}

func NewFirestoreNews(videoID string, videoStatus int, updateAt wrappedbasics.IWrappedTime, participants []string) FirestoreNews {
	return FirestoreNews{
		UpdateAt:     updateAt,
		VideoID:      videoID,
		VideoStatus:  videoStatus,
		Participants: participants,
	}
}
