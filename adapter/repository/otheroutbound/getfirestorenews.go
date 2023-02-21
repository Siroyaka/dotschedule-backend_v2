package otheroutbound

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetFirestoreNewsRepository struct {
	firestore                abstruct.Firestore
	collectionName           string
	compareble               string
	operator                 string
	documentNameVideoID      string
	documentNameVideoStatus  string
	documentNameParticipants string
	documentNameUpdateAt     string
}

func NewGetFirestoreNewsRepository(
	firestore abstruct.Firestore,
	collectionName, compareble, operator string,
	documentNameVideoID,
	documentNameVideoStatus,
	documentNameParticipants,
	documentNameUpdateAt string,
) GetFirestoreNewsRepository {
	return GetFirestoreNewsRepository{
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

func (repos GetFirestoreNewsRepository) converter(videoID string, videoStatus int, updateAt string, participants map[string]bool) (domain.FirestoreNews, utility.IError) {
	var participantsList []string
	for k := range participants {
		participantsList = append(participantsList, k)
	}

	// update at convert to wrapped time. need common.
	updateAtTime, err := wrappedbasics.NewWrappedTimeFromLocal(updateAt, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		return domain.FirestoreNews{}, err.WrapError()
	}

	return domain.NewFirestoreNews(videoID, videoStatus, updateAtTime, participantsList), nil
}

func (repos GetFirestoreNewsRepository) parseFirestoreData(m map[string]interface{}) (videoID string, videoStatus int, updateAt string, participants map[string]bool, err utility.IError) {
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

func (repos GetFirestoreNewsRepository) Execute(
	targetTime wrappedbasics.IWrappedTime,
) ([]domain.FirestoreNews, utility.IError) {
	t := targetTime.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
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

		videoId, videoStatus, updateAt, participants, err := repos.parseFirestoreData(dataMap)
		if err != nil {
			utility.LogError(err)
			continue
		}

		data, err := repos.converter(videoId, videoStatus, updateAt, participants)
		if err != nil {
			utility.LogError(err.WrapError("convert"))
			continue
		}
		responseData = append(responseData, data)
	}
	return responseData, nil
}
