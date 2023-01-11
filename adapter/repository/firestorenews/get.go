package firestorenews

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/firestorenews"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository struct {
	firestore                abstruct.Firestore
	collectionName           string
	compareble               string
	operator                 string
	documentNameVideoID      string
	documentNameVideoStatus  string
	documentNameParticipants string
	documentNameUpdateAt     string
}

func NewGetRepository(
	firestore abstruct.Firestore,
	collectionName, compareble, operator string,
	documentNameVideoID,
	documentNameVideoStatus,
	documentNameParticipants,
	documentNameUpdateAt string,
) firestorenews.GetRepository {
	return GetRepository{
		firestore:                firestore,
		collectionName:           collectionName,
		compareble:               compareble,
		operator:                 operator,
		documentNameParticipants: documentNameParticipants,
		documentNameUpdateAt:     documentNameUpdateAt,
		documentNameVideoID:      documentNameVideoID,
		documentNameVideoStatus:  documentNameVideoStatus,
	}
}

func (repos GetRepository) Get(
	converter firestorenews.GetDataConverter,
	targetTime utility.WrappedTime,
) ([]domain.FirestoreNews, utility.IError) {
	parser := func(m map[string]interface{}) (videoID string, videoStatus int, updateAt string, participants map[string]bool, err utility.IError) {
		videoID, err = utility.ConvertFromInterfaceType[string](m[repos.documentNameVideoID])
		if err != nil {
			err = err.WrapError("videoId type error")
			return
		}
		i64vs, err := utility.ConvertFromInterfaceType[int64](m[repos.documentNameVideoStatus])
		if err != nil {
			err = err.WrapError("videoStatus type error")
			return
		}
		videoStatus = int(i64vs)
		updateAt, err = utility.ConvertFromInterfaceType[string](m[repos.documentNameUpdateAt])
		if err != nil {
			err = err.WrapError("updateAt type error")
			return
		}
		partMapInterface, err := utility.ConvertFromInterfaceType[map[string]interface{}](m[repos.documentNameParticipants])
		if err != nil {
			err = err.WrapError("participants type error")
			return
		}

		// participants map convert to map[string]bool from map[string]interface{}
		participants = make(map[string]bool)
		for key := range partMapInterface {
			participants[key] = true
		}
		return
	}

	t := targetTime.ToLocalFormatString()
	iter := repos.firestore.Collection(repos.collectionName).Where(repos.compareble, repos.operator, t).Documents(repos.firestore.GetContext())
	var responseData []domain.FirestoreNews
	for {
		ok, snapshot, err := iter.Next()
		if !ok {
			break
		}
		if err != nil {
			return nil, err.WrapError("no item")
		}
		dataMap := snapshot.Data()
		videoId, videoStatus, updateAt, participants, err := parser(dataMap)
		if err != nil {
			utility.LogError(err)
			continue
		}
		data, err := converter(videoId, videoStatus, updateAt, participants)
		if err != nil {
			utility.LogError(err.WrapError("convert"))
			continue
		}
		responseData = append(responseData, data)
	}
	return responseData, nil
}
