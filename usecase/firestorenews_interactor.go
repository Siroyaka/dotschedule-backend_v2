package usecase

import (
	"fmt"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/firestorenews"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamermaster"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type FirestoreNewsInteractor struct {
	getFirestoreNewsRepos            firestorenews.GetRepository
	countScheduleRepos               fullschedule.CountRepository
	insertFullScheduleRepos          fullschedule.InsertRepository
	updateFullScheduleRepos          fullschedule.UpdateAnyColumnRepository
	getPlatformIdRepos               streamermaster.GetPlatformIdRepository
	getStreamingParticipantsRepos    streamingparticipants.GetRepository
	insertStreamingParticipantsRepos streamingparticipants.InsertRepository
	deleteStreamingParticipantsRepos streamingparticipants.DeleteRepository
	common                           utility.Common
	firestoreTargetMin               int
	platform                         string
	firestoreData                    []domain.FirestoreNews
}

func NewFirestoreNewsInteractor(
	getFirestoreNewsRepos firestorenews.GetRepository,
	countScheduleRepos fullschedule.CountRepository,
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
		countScheduleRepos:               countScheduleRepos,
		insertFullScheduleRepos:          insertFullScheduleRepos,
		updateFullScheduleRepos:          updateFullScheduleRepos,
		getPlatformIdRepos:               getPlatformIdRepos,
		getStreamingParticipantsRepos:    getStreamingParticipantsRepos,
		insertStreamingParticipantsRepos: insertStreamingParticipantsRepos,
		deleteStreamingParticipantsRepos: deleteStreamingParticipantsRepos,
		common:                           common,
		firestoreTargetMin:               firestoreTargetMin,
		platform:                         platform,
		firestoreData:                    []domain.FirestoreNews{},
	}
}

func (intr FirestoreNewsInteractor) converter(videoID string, videoStatus int, updateAt string, participants map[string]bool) (domain.FirestoreNews, utility.IError) {
	var participantsList []string
	for k := range participants {
		participantsList = append(participantsList, k)
	}

	// update at convert to wrapped time. need common.
	updateAtTime, err := intr.common.CreateNewWrappedTimeFromLocal(updateAt)
	if err != nil {
		return domain.FirestoreNews{}, err.WrapError()
	}

	return domain.NewFirestoreNews(videoID, videoStatus, updateAtTime, participantsList), nil
}

func (intr *FirestoreNewsInteractor) DataFetchFromFirestore() utility.IError {
	now, err := intr.common.Now()
	if err != nil {
		return err.WrapError()
	}
	targetTime := now.Add(0, 0, 0, 0, -1*intr.firestoreTargetMin, 0)

	utility.LogDebug(fmt.Sprintf("firestore target time: %s - %s", targetTime.ToLocalFormatString(), now.ToLocalFormatString()))

	list, err := intr.getFirestoreNewsRepos.Get(intr.converter, targetTime)
	if err != nil {
		return err.WrapError()
	}
	intr.firestoreData = list
	return nil
}

func (intr FirestoreNewsInteractor) firestoreNewsToFullSchedule(data domain.FirestoreNews) domain.FullScheduleData {
	fullScheduleData := domain.NewEmptyFullScheduleData(data.VideoID, intr.platform)
	fullScheduleData = fullScheduleData.ImportStatusFromFirestore(data.VideoStatus)
	return fullScheduleData
}

func (intr FirestoreNewsInteractor) updateSchedule(data domain.FullScheduleData) utility.IError {
	now, err := intr.common.Now()
	if err != nil {
		return err.WrapError()
	}

	cnt, err := intr.countScheduleRepos.Count(data.StreamingID, data.PlatformType)
	if err != nil {
		return err.WrapError("schedule data count error")
	}

	if cnt == 0 {
		utility.LogInfo(fmt.Sprintf("insert schedule. id: %s", data.StreamingID))
		insertCount, err := intr.insertFullScheduleRepos.Insert(data, now)
		if err != nil {
			return err.WrapError("schedule data insert error")
		}
		if insertCount == 0 {
			return utility.NewError("schedule data insert count is 0", "")
		}
	} else {
		utility.LogInfo(fmt.Sprintf("update schedule. id: %s", data.StreamingID))
		updateCount, err := intr.updateFullScheduleRepos.Update(now, data.IsCompleteData, data.StreamingID, data.PlatformType)
		if err != nil {
			return err.WrapError("schedule data update error")
		}
		if updateCount == 0 {
			return utility.NewError("schedule data update count is 0", "")
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

func (intr FirestoreNewsInteractor) UpdateDB() {
	if len(intr.firestoreData) == 0 {
		utility.LogDebug("firestoreNews no data")
		return
	}
	for _, data := range intr.firestoreData {
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