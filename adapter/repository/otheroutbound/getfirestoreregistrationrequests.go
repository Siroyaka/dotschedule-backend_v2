package otheroutbound

import (
	"fmt"
	"time"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetFirestoreRegistrationRequestRepository struct {
	firestore                abstruct.Firestore
	collectionName           string
	compareble               string
	operator                 string
	documentNameID           string
	documentNameURL          string
	documentNamePlatform     string
	documentNameParticipants string
	documentNameUpdateAt     string
	documentNameStartDate    string
	documentNameStreamerID   string
	documentNameStreamerName string
	documentNameTitle        string
}

func NewGetFirestoreRegistrationRequestRepository(
	firestore abstruct.Firestore,
	collectionName, compareble, operator string,
	documentNameID,
	documentNameURL,
	documentNamePlatform,
	documentNameParticipants,
	documentNameUpdateAt,
	documentNameStartDate,
	documentNameStreamerID,
	documentNameStreamerName,
	documentNameTitle string,
) GetFirestoreRegistrationRequestRepository {
	return GetFirestoreRegistrationRequestRepository{
		firestore:                firestore,
		collectionName:           collectionName,
		compareble:               compareble,
		operator:                 operator,
		documentNameID:           documentNameID,
		documentNameURL:          documentNameURL,
		documentNamePlatform:     documentNamePlatform,
		documentNameParticipants: documentNameParticipants,
		documentNameUpdateAt:     documentNameUpdateAt,
		documentNameStartDate:    documentNameStartDate,
		documentNameStreamerID:   documentNameStreamerID,
		documentNameStreamerName: documentNameStreamerName,
		documentNameTitle:        documentNameTitle,
	}
}

func (repos GetFirestoreRegistrationRequestRepository) convert(m map[string]interface{}) (domain.FirestoreRegistrationRequest, utilerror.IError) {
	emptyRegistrationRequest := domain.FirestoreRegistrationRequest{}

	id, notexists, err := utility.PickupMapInterfaceData[string](m, repos.documentNameID)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameID))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("id does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameID)
	}

	participantsData, notexists, err := utility.PickupMapInterfaceData[[]interface{}](m, repos.documentNameParticipants)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameParticipants))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("participants does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameParticipants)
	}
	var participants []string
	for _, item := range participantsData {
		convertData, err := utility.ConvertFromInterfaceType[string](item)
		if err != nil {
			return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameParticipants))
		}
		participants = append(participants, convertData)
	}

	platform, notexists, err := utility.PickupMapInterfaceData[string](m, repos.documentNamePlatform)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNamePlatform))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("platform does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNamePlatform)
	}

	url, notexists, err := utility.PickupMapInterfaceData[string](m, repos.documentNameURL)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameURL))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("url does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameURL)
	}

	updateAtTime, notexists, err := utility.PickupMapInterfaceData[time.Time](m, repos.documentNameUpdateAt)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameUpdateAt))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("updateAt does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameUpdateAt)
	}

	updateAt, err := wrappedbasics.NewWrappedTimeFromUTC(updateAtTime.UTC().Format(string(wrappedbasics.WrappedTimeProps.DateTimeFormat())), wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		return emptyRegistrationRequest, err.WrapError("updateAt parse error")
	}

	if platform == "YOUTUBE" {
		return domain.NewFirestoreRegistrationRequestYoutube(
			id,
			participants,
			url,
			platform,
			updateAt,
		), nil
	}

	title, notexists, err := utility.PickupMapInterfaceData[string](m, repos.documentNameTitle)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameTitle))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("title does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameTitle)
	}

	streamerName, notexists, err := utility.PickupMapInterfaceData[string](m, repos.documentNameStreamerName)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameStreamerName))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("streamerName does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameStreamerName)
	}

	streamerID, notexists, err := utility.PickupMapInterfaceData[string](m, repos.documentNameStreamerID)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameStreamerID))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("streamerID does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameStreamerID)
	}

	startDateTime, notexists, err := utility.PickupMapInterfaceData[time.Time](m, repos.documentNameStartDate)
	if err != nil {
		return emptyRegistrationRequest, err.WrapError(fmt.Sprintf("objectName: %s", repos.documentNameStartDate))
	}
	if notexists {
		return emptyRegistrationRequest, utilerror.New("startDate does not exists in firestore data", utilerror.ERR_FIRESTORE_DATA_NOTEXISTS, repos.documentNameStartDate)
	}

	startDate, err := wrappedbasics.NewWrappedTimeFromUTC(startDateTime.UTC().Format(string(wrappedbasics.WrappedTimeProps.DateTimeFormat())), wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		return emptyRegistrationRequest, err.WrapError("startDate parse error")
	}

	return domain.NewFirestoreRegistrationRequestOthers(
		id,
		participants,
		url,
		platform,
		updateAt,
		startDate,
		streamerID,
		streamerName,
		title,
	), nil
}

func (repos GetFirestoreRegistrationRequestRepository) Execute(
	targetTime wrappedbasics.IWrappedTime,
) ([]domain.FirestoreRegistrationRequest, utilerror.IError) {
	t := targetTime.Time()

	iter := repos.firestore.Collection(repos.collectionName).Where(repos.compareble, repos.operator, t).Documents(repos.firestore.GetContext())
	var responseData []domain.FirestoreRegistrationRequest
	for {
		ok, snapshot, err := iter.Next()
		if !ok {
			break
		}
		if err != nil {
			return nil, err.WrapError("no item")
		}
		snapshotData := snapshot.Data()

		data, err := repos.convert(snapshotData)
		if err != nil {
			logger.Error(err)
			continue
		}

		responseData = append(responseData, data)
	}
	return responseData, nil
}
