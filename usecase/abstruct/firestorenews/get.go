package firestorenews

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository interface {
	Get(GetDataConverter, utility.WrappedTime) ([]domain.FirestoreNews, utility.IError)
}

type GetDataConverter func(videoID string, videoStatus int, updateAt string, participants map[string]bool) (domain.FirestoreNews, utility.IError)
