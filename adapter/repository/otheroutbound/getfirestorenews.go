package otheroutbound

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
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

func (repos GetFirestoreNewsRepository) converter(videoID string, updateAt string, participants map[string]bool) (domain.FirestoreNews, utilerror.IError) {
	var participantsList []string
	for k := range participants {
		participantsList = append(participantsList, k)
	}

	updateAtTime, err := wrappedbasics.NewWrappedTimeFromLocal(updateAt, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		return domain.FirestoreNews{}, err.WrapError()
	}

	return domain.NewFirestoreNews(videoID, updateAtTime, participantsList), nil
}

func (repos GetFirestoreNewsRepository) parseFirestoreData(m map[string]interface{}) (videoID string, updateAt string, participants map[string]bool, err utilerror.IError) {
	videoID, err = utility.ConvertFromInterfaceType[string](m[repos.documentNameVideoID])
	if err != nil {
		err = err.WrapError("videoId type error")
		return
	}
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
) ([]domain.FirestoreNews, utilerror.IError) {
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

		videoId, updateAt, participants, err := repos.parseFirestoreData(dataMap)
		if err != nil {
			logger.Error(err)
			continue
		}

		data, err := repos.converter(videoId, updateAt, participants)
		if err != nil {
			logger.Error(err.WrapError("convert"))
			continue
		}
		responseData = append(responseData, data)
	}
	return responseData, nil
}
