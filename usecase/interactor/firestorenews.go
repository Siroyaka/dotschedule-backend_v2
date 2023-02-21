package interactor

import (
	"fmt"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamermaster"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type FirestoreNewsInteractor struct {
	getFirestoreNewsRepos abstruct.RepositoryRequest[wrappedbasics.IWrappedTime, []domain.FirestoreNews]

	containsScheduleRepos abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, bool]

	insertFullScheduleRepos fullschedule.InsertRepository
	updateFullScheduleRepos fullschedule.UpdateAnyColumnRepository

	getPlatformIdRepos streamermaster.GetPlatformIdRepository

	getStreamingParticipantsRepos    streamingparticipants.GetRepository
	insertStreamingParticipantsRepos streamingparticipants.InsertRepository
	deleteStreamingParticipantsRepos streamingparticipants.DeleteRepository

	common             utility.Common
	firestoreTargetMin int
	platform           string
}

func NewFirestoreNewsInteractor(
	getFirestoreNewsRepos abstruct.RepositoryRequest[wrappedbasics.IWrappedTime, []domain.FirestoreNews],
	containsScheduleRepos abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, bool],

	insertFullScheduleRepos fullschedule.InsertRepository,
	updateFullScheduleRepos fullschedule.UpdateAnyColumnRepository,
	getPlatformIdRepos streamermaster.GetPlatformIdRepository,
	getStreamingParticipantsRepos streamingparticipants.GetRepository,
	insertStreamingParticipantsRepos streamingparticipants.InsertRepository,
	deleteStreamingParticipantsRepos streamingparticipants.DeleteRepository,
	common utility.Common,
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
		common:                           common,
		firestoreTargetMin:               firestoreTargetMin,
		platform:                         platform,
	}
}

func (intr FirestoreNewsInteractor) DataFetchFromFirestore() ([]domain.FirestoreNews, utility.IError) {
	now := wrappedbasics.Now()

	targetTime := now.Add(0, 0, 0, 0, -1*intr.firestoreTargetMin, 0)

	utility.LogInfo(fmt.Sprintf("target time UTC: %s - %s", targetTime.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()), now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())))

	list, err := intr.getFirestoreNewsRepos.Execute(targetTime)
	if err != nil {
		return nil, err.WrapError()
	}

	return list, nil
}

func (intr FirestoreNewsInteractor) firestoreNewsToFullSchedule(data domain.FirestoreNews) domain.FullScheduleData {
	fullScheduleData := domain.NewEmptyFullScheduleData(data.VideoID, intr.platform)
	fullScheduleData = fullScheduleData.ImportStatusFromFirestore(data.VideoStatus)
	return fullScheduleData
}

// dbにすでにデータが存在しているか確認し、存在しているならupdate、存在していないならinsertを行う
func (intr FirestoreNewsInteractor) updateSchedule(data domain.FullScheduleData) utility.IError {
	now, err := intr.common.Now()
	if err != nil {
		return err.WrapError()
	}

	alreadyContains, err := intr.containsScheduleRepos.Execute(reference.NewStreamingIDWithPlatformType(data.StreamingID, data.PlatformType))
	if err != nil {
		return err.WrapError("schedule data count error")
	}

	if alreadyContains {
		utility.LogInfo(fmt.Sprintf("update schedule. id: %s", data.StreamingID))
		updateCount, err := intr.updateFullScheduleRepos.Update(now, data.IsCompleteData, data.StreamingID, data.PlatformType)
		if err != nil {
			return err.WrapError("schedule data update error")
		}
		if updateCount == 0 {
			return utility.NewError("schedule data update count is 0", "")
		}
	} else {
		utility.LogInfo(fmt.Sprintf("insert schedule. id: %s", data.StreamingID))
		insertCount, err := intr.insertFullScheduleRepos.Insert(data, now)
		if err != nil {
			return err.WrapError("schedule data insert error")
		}
		if insertCount == 0 {
			return utility.NewError("schedule data insert count is 0", "")
		}
	}
	return nil
}

func (intr FirestoreNewsInteractor) firestoreNewsToStreamingParticipants(data domain.FirestoreNews) (domain.StreamingParticipants, utility.IError) {
	platformToStreamerIdMap, err := intr.getPlatformIdRepos.GetPlatformIdToStreamerId(intr.platform)
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

func (intr FirestoreNewsInteractor) updateParticipants(data domain.StreamingParticipants) utility.IError {
	now, err := intr.common.Now()
	if err != nil {
		return err.WrapError()
	}

	dbData, err := intr.getStreamingParticipantsRepos.Get(data.StreamingID(), data.Platform())
	if err != nil {
		return err.WrapError("participants data get error")
	}

	deleteList := domain.NewStreamingParticipants(data.StreamingID(), data.Platform())
	for _, v := range dbData.GetList() {
		if data.Has(v) {
			continue
		}
		deleteList = deleteList.Add(v)
	}

	var responseError utility.IError
	responseError = nil
	if !deleteList.IsEmpty() {
		utility.LogInfo(fmt.Sprintf("delete participants data. id: %s [%s]", deleteList.StreamingID(), strings.Join(deleteList.GetList(), ", ")))
		deleteCount, err := intr.deleteStreamingParticipantsRepos.Delete(deleteList)
		if err != nil {
			utility.LogError(err.WrapError("delete participants error"))
			responseError = err.WrapError()
		} else if deleteCount < int64(len(deleteList.GetList())) {
			utility.LogError(utility.NewError(fmt.Sprintf("delete data count wrong. list count: %d. delete count: %d", len(deleteList.GetList()), deleteCount), ""))
		}
	}

	insertList := domain.NewStreamingParticipants(data.StreamingID(), data.Platform())
	for _, v := range data.GetList() {
		if dbData.Has(v) {
			continue
		}
		insertList = insertList.Add(v)
	}

	if !insertList.IsEmpty() {
		utility.LogInfo(fmt.Sprintf("insert participants data. id: %s [%s]", insertList.StreamingID(), strings.Join(insertList.GetList(), ", ")))
		insertCount, err := intr.insertStreamingParticipantsRepos.InsertStreamingParticipants(insertList, now)
		if err != nil {
			utility.LogError(err.WrapError("insert participants error"))
			responseError = err.WrapError()
		} else if insertCount < int64(len(insertList.GetList())) {
			utility.LogError(utility.NewError(fmt.Sprintf("insert data count wrong. list count: %d. insert count: %d", len(insertList.GetList()), insertCount), ""))
		}
	}
	return responseError
}

func (intr FirestoreNewsInteractor) UpdateDB(firestoreData []domain.FirestoreNews) {
	for _, data := range firestoreData {
		fullScheduleData := intr.firestoreNewsToFullSchedule(data)
		if err := intr.updateSchedule(fullScheduleData); err != nil {
			utility.LogFatal(err.WrapError(fmt.Sprintf("schedule update error id: %s", data.VideoID)))
		}

		streamingParticipants, err := intr.firestoreNewsToStreamingParticipants(data)
		if err != nil {
			utility.LogFatal(err.WrapError(fmt.Sprintf("platform id convert to streamer id error. id: %s", data.VideoID)))
			return
		}
		if err := intr.updateParticipants(streamingParticipants); err != nil {
			utility.LogFatal(err.WrapError(fmt.Sprintf("participants update error id: %s", data.VideoID)))
		}
	}
}
