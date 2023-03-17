package interactor

import (
	"fmt"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type FirestoreNewsInteractor struct {
	getFirestoreNewsRepos abstruct.RepositoryRequest[wrappedbasics.IWrappedTime, []domain.FirestoreNews]

	containsScheduleRepos abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, bool]

	insertFullScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse]
	updateFullScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse]

	getPlatformIdRepos abstruct.RepositoryRequest[string, map[string]string]

	getStreamingParticipantsRepos    abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, domain.StreamingParticipants]
	insertStreamingParticipantsRepos abstruct.RepositoryRequest[domain.StreamingParticipants, reference.DBUpdateResponse]
	deleteStreamingParticipantsRepos abstruct.RepositoryRequest[domain.StreamingParticipants, reference.DBUpdateResponse]

	firestoreTargetMin int
	platform           string
}

func NewFirestoreNewsInteractor(
	getFirestoreNewsRepos abstruct.RepositoryRequest[wrappedbasics.IWrappedTime, []domain.FirestoreNews],
	containsScheduleRepos abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, bool],

	insertFullScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse],
	updateFullScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse],

	getPlatformIdRepos abstruct.RepositoryRequest[string, map[string]string],
	getStreamingParticipantsRepos abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, domain.StreamingParticipants],
	insertStreamingParticipantsRepos abstruct.RepositoryRequest[domain.StreamingParticipants, reference.DBUpdateResponse],
	deleteStreamingParticipantsRepos abstruct.RepositoryRequest[domain.StreamingParticipants, reference.DBUpdateResponse],

	firestoreTargetMin int,
	platform string,
) FirestoreNewsInteractor {
	return FirestoreNewsInteractor{
		getFirestoreNewsRepos:            getFirestoreNewsRepos,
		containsScheduleRepos:            containsScheduleRepos,
		insertFullScheduleRepos:          insertFullScheduleRepos,
		updateFullScheduleRepos:          updateFullScheduleRepos,
		getPlatformIdRepos:               getPlatformIdRepos,
		getStreamingParticipantsRepos:    getStreamingParticipantsRepos,
		insertStreamingParticipantsRepos: insertStreamingParticipantsRepos,
		deleteStreamingParticipantsRepos: deleteStreamingParticipantsRepos,
		firestoreTargetMin:               firestoreTargetMin,
		platform:                         platform,
	}
}

func (intr FirestoreNewsInteractor) DataFetchFromFirestore() ([]domain.FirestoreNews, utilerror.IError) {
	now := wrappedbasics.Now()

	targetTime := now.Add(0, 0, 0, 0, -1*intr.firestoreTargetMin, 0)

	logger.Info(fmt.Sprintf("target time UTC: %s - %s", targetTime.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()), now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())))

	list, err := intr.getFirestoreNewsRepos.Execute(targetTime)
	if err != nil {
		return nil, err.WrapError()
	}

	return list, nil
}

func (intr FirestoreNewsInteractor) firestoreNewsToFullSchedule(data domain.FirestoreNews) domain.FullScheduleData {
	fullScheduleData := domain.NewEmptyFullScheduleData(data.VideoID, intr.platform)

	// every firestore news data are status 10.
	fullScheduleData = fullScheduleData.ImportStatusFromFirestore(10)
	return fullScheduleData
}

// dbにすでにデータが存在しているか確認し、存在しているならupdate、存在していないならinsertを行う
func (intr FirestoreNewsInteractor) updateSchedule(data domain.FullScheduleData) utilerror.IError {

	alreadyContains, err := intr.containsScheduleRepos.Execute(reference.NewStreamingIDWithPlatformType(data.StreamingID, data.PlatformType))
	if err != nil {
		return err.WrapError("schedule data count error")
	}

	if alreadyContains {
		logger.Info(fmt.Sprintf("update schedule. id: %s", data.StreamingID))
		if updateResult, err := intr.updateFullScheduleRepos.Execute(data); err != nil {
			return err.WrapError("schedule data update error")
		} else if updateResult.Count == 0 {
			return utilerror.New("schedule data update count is 0", "")
		}
	} else {
		logger.Info(fmt.Sprintf("insert schedule. id: %s", data.StreamingID))
		if insertResult, err := intr.insertFullScheduleRepos.Execute(data); err != nil {
			return err.WrapError("schedule data insert error")
		} else if insertResult.Count == 0 {
			return utilerror.New("schedule data insert count is 0", "")
		}
	}
	return nil
}

// firestorenewsのデータをparticipantsのデータに変換
func (intr FirestoreNewsInteractor) firestoreNewsToStreamingParticipants(data domain.FirestoreNews) (domain.StreamingParticipants, utilerror.IError) {

	platformToStreamerIdMap, err := intr.getPlatformIdRepos.Execute(intr.platform)
	if err != nil {
		return domain.EmptyStreamingParticipants(), err.WrapError()
	}

	platformParticipants := domain.NewPlatformParticipants(data.VideoID, intr.platform, data.Participants...)
	for k, v := range platformToStreamerIdMap {
		platformParticipants.AddConvertData(v, k)
	}
	res := platformParticipants.Convert()
	return res, nil
}

// Firestoreから取得したparticipantsのデータとDBのparticipantsのデータを揃える(Firestoreのデータを正とする)
func (intr FirestoreNewsInteractor) updateParticipants(data domain.StreamingParticipants) utilerror.IError {

	dbData, err := intr.getStreamingParticipantsRepos.Execute(reference.NewStreamingIDWithPlatformType(data.StreamingID(), data.Platform()))
	if err != nil {
		return err.WrapError("participants data get error")
	}

	// firestoreのデータに入っていないものを削除データとしてピックアップする
	deleteList := domain.NewStreamingParticipants(data.StreamingID(), data.Platform())
	for _, v := range dbData.GetList() {
		if data.Has(v) {
			continue
		}
		deleteList = deleteList.Add(v)
	}

	var responseError utilerror.IError
	responseError = nil
	if !deleteList.IsEmpty() {
		logger.Info(fmt.Sprintf("delete participants data. id: %s [%s]", deleteList.StreamingID(), strings.Join(deleteList.GetList(), ", ")))

		if deleteResult, err := intr.deleteStreamingParticipantsRepos.Execute(deleteList); err != nil {
			logger.Error(err.WrapError("delete participants error"))
			responseError = err.WrapError()
		} else if deleteResult.Count < int64(len(deleteList.GetList())) {
			logger.Error(utilerror.New(fmt.Sprintf("delete data count wrong. list count: %d. delete count: %d", len(deleteList.GetList()), deleteResult.Count), ""))
		}
	}

	// firestoreにあるがDBにないデータを追加用のデータとしてピックアップする
	insertList := domain.NewStreamingParticipants(data.StreamingID(), data.Platform())
	for _, v := range data.GetList() {
		if dbData.Has(v) {
			continue
		}
		insertList = insertList.Add(v)
	}

	if !insertList.IsEmpty() {
		logger.Info(fmt.Sprintf("insert participants data. id: %s [%s]", insertList.StreamingID(), strings.Join(insertList.GetList(), ", ")))

		if insertResult, err := intr.insertStreamingParticipantsRepos.Execute(insertList); err != nil {
			logger.Error(err.WrapError("insert participants error"))
			responseError = err.WrapError()
		} else if insertResult.Count < int64(len(insertList.GetList())) {
			logger.Error(utilerror.New(fmt.Sprintf("insert data count wrong. list count: %d. insert count: %d", len(insertList.GetList()), insertResult.Count), ""))
		}
	}
	return responseError
}

func (intr FirestoreNewsInteractor) UpdateDB(firestoreData []domain.FirestoreNews) {
	for _, data := range firestoreData {
		fullScheduleData := intr.firestoreNewsToFullSchedule(data)
		if err := intr.updateSchedule(fullScheduleData); err != nil {
			logger.Fatal(err.WrapError(fmt.Sprintf("schedule update error id: %s", data.VideoID)))
		}

		streamingParticipants, err := intr.firestoreNewsToStreamingParticipants(data)
		if err != nil {
			logger.Fatal(err.WrapError(fmt.Sprintf("platform id convert to streamer id error. id: %s", data.VideoID)))
			return
		}

		if err := intr.updateParticipants(streamingParticipants); err != nil {
			logger.Fatal(err.WrapError(fmt.Sprintf("participants update error id: %s", data.VideoID)))
		}
	}
}
