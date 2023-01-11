package domain

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type FirestoreNews struct {
	UpdateAt     utility.WrappedTime
	VideoID      string
	VideoStatus  int
	Participants []string
}

func NewFirestoreNews(videoID string, videoStatus int, updateAt utility.WrappedTime, participants []string) FirestoreNews {
	return FirestoreNews{
		UpdateAt:     updateAt,
		VideoID:      videoID,
		VideoStatus:  videoStatus,
		Participants: participants,
	}
}
