package interactor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type RSSInteractor struct {
	index      int
	masterList []domain.RSSMaster

	getMasterRepository    abstruct.RepositoryRequest[reference.VoidStruct, []domain.RSSMaster]
	updateMasterRepository abstruct.RepositoryRequest[reference.IDWithTime, reference.DBUpdateResponse]

	requestRepository abstruct.RepositoryRequest[string, utility.IFeed]

	getScheduleIDRepository  abstruct.RepositoryRequest[reference.StreamingIDListWithPlatformID, utility.HashSet[string]]
	insertScheduleRepository abstruct.RepositoryRequest[domain.SeedSchedule, reference.DBUpdateResponse]
	updateScheduleRepository abstruct.RepositoryRequest[reference.StreamingIDListWithPlatformID, reference.DBUpdateResponse]

	completeStatus, notCompleteStatus                       int
	platform, insertStatus, videoIdElement, videoIdItemName string
}

func NewRSSInteractor(
	getMasterRepository abstruct.RepositoryRequest[reference.VoidStruct, []domain.RSSMaster],
	updateMasterRepository abstruct.RepositoryRequest[reference.IDWithTime, reference.DBUpdateResponse],

	requestRepository abstruct.RepositoryRequest[string, utility.IFeed],

	getScheduleIDRepository abstruct.RepositoryRequest[reference.StreamingIDListWithPlatformID, utility.HashSet[string]],
	insertScheduleRepository abstruct.RepositoryRequest[domain.SeedSchedule, reference.DBUpdateResponse],
	updateScheduleRepository abstruct.RepositoryRequest[reference.StreamingIDListWithPlatformID, reference.DBUpdateResponse],

	completeStatus, notCompleteStatus int,
	platform, insertStatus, videoIdElement, videoIdItemName string,
) RSSInteractor {
	var intr = RSSInteractor{
		index:                    -1,
		masterList:               []domain.RSSMaster{},
		getMasterRepository:      getMasterRepository,
		updateMasterRepository:   updateMasterRepository,
		requestRepository:        requestRepository,
		getScheduleIDRepository:  getScheduleIDRepository,
		insertScheduleRepository: insertScheduleRepository,
		updateScheduleRepository: updateScheduleRepository,
		completeStatus:           completeStatus,
		notCompleteStatus:        notCompleteStatus,
		platform:                 platform,
		insertStatus:             insertStatus,
		videoIdElement:           videoIdElement,
		videoIdItemName:          videoIdItemName,
	}
	return intr
}

func (intr *RSSInteractor) GetMaster() utility.IError {
	list, err := intr.getMasterRepository.Execute(reference.Void())
	intr.masterList = list
	return err
}

func (intr *RSSInteractor) Next() bool {
	if len(intr.masterList) > intr.index {
		intr.index++
	}
	return len(intr.masterList) > intr.index
}

func (intr RSSInteractor) feedDataToSeedSchedule(
	feedData utility.IFeed,
	platformType, feedStatus, videoIdElement, videoIdItemName string,
	master domain.RSSMaster,
) []domain.SeedSchedule {
	var seedScheduleList []domain.SeedSchedule
	for i := 0; i < feedData.GetItemLength(); i++ {
		feedItem := feedData.GetItem(i)
		if feedItem == nil {
			utility.LogError(utility.NewError(fmt.Sprintf("items out of index. index: %v", i), utility.ERR_RSS_PARSE))
			continue
		}

		feedPublishedAt, err := wrappedbasics.NewWrappedTimeFromUTC(feedItem.GetPublishedAt(), wrappedbasics.WrappedTimeProps.DateTimeFormat())
		if err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("feed updatetime parse error. Title: %s, UpdateTime: %s", feedItem.GetTitle(), feedItem.GetUpdateAt())))
			continue
		}

		if !feedPublishedAt.After(master.LastUpdate) {
			// feedPublishedAt <= master.LastUpdate
			continue
		}

		_, videoId, _ := feedItem.GetExtension(videoIdElement, videoIdItemName)

		seedSchedule := domain.NewSeedSchedule(videoId, platformType, feedStatus, feedPublishedAt)
		seedScheduleList = append(seedScheduleList, seedSchedule)
	}

	return seedScheduleList
}

func (intr RSSInteractor) GetRSSData() ([]domain.SeedSchedule, utility.IError) {
	if len(intr.masterList) <= intr.index {
		return nil, utility.NewError("", utility.ERR_OUTOFINDEX, "ri.masterList", strconv.Itoa(len(intr.masterList)), strconv.Itoa(intr.index))
	}
	targetMaster := intr.masterList[intr.index]

	utility.LogDebug(fmt.Sprintf("feed request: %s %s", targetMaster.ID, targetMaster.Name))

	feedData, err := intr.requestRepository.Execute(targetMaster.Url)
	if err != nil {
		return nil, err.WrapError()
	}

	list := intr.feedDataToSeedSchedule(
		feedData,
		intr.platform,
		intr.insertStatus,
		intr.videoIdElement,
		intr.videoIdItemName,
		targetMaster,
	)

	return list, nil
}

func (intr RSSInteractor) PushToDB(list []domain.SeedSchedule) (insertCount, updateCount int, isError bool, newestPublishedAt wrappedbasics.WrappedTime, ierr utility.IError) {
	var idList []string
	isError = false

	master := intr.masterList[intr.index]
	newestPublishedAt = master.LastUpdate

	for _, seedSchedule := range list {
		idList = append(idList, seedSchedule.GetID())
	}

	alreadyInsertedIdList, err := intr.getScheduleIDRepository.Execute(reference.NewStreamingIDListWithPlatformID(idList, intr.platform))
	if err != nil {
		ierr = err.WrapError()
		isError = true
		return
	}

	var updateList []string

	for _, seedSchedule := range list {
		if seedSchedule.GetPublishedAt().After(newestPublishedAt) {
			newestPublishedAt = seedSchedule.GetPublishedAt()
		}

		if alreadyInsertedIdList.Has(seedSchedule.GetID()) {
			updateList = append(updateList, seedSchedule.GetID())
			continue
		}

		if result, err := intr.insertScheduleRepository.Execute(seedSchedule); err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("schedule insert error. name: %s, streaming_id: [ %s ]", master.Name, seedSchedule.GetID())))
			isError = true
			continue
		} else if result.Count == 0 {
			utility.LogError(err.WrapError(fmt.Sprintf("schedule insert error. count 0. name: %s, streaming_id: [ %s ]", master.Name, seedSchedule.GetID())))
			isError = true
			continue
		}

		utility.LogInfo(fmt.Sprintf("schedule insert finished. name: %s, streaming_id: [ %s ]", master.Name, seedSchedule.GetID()))

		insertCount++
	}

	if updateList != nil {
		if result, err := intr.updateScheduleRepository.Execute(reference.NewStreamingIDListWithPlatformID(updateList, intr.platform)); err != nil {
			utility.LogDebug(err.Error())
			ierr = err.WrapError(fmt.Sprintf("schedule update ERROR. name: %s, streaming_id: [ %s ]", master.Name, strings.Join(updateList, ", ")))
			isError = true
			return
		} else if result.Count == 0 {
			utility.LogDebug("test4")
			ierr = err.WrapError(fmt.Sprintf("schedule update ERROR. name: %s, streaming_id: [ %s ]", master.Name, strings.Join(updateList, ", ")))
			isError = true
			return
		}

		utility.LogInfo(fmt.Sprintf("schedule update finished. name: %s, streaming_id: [ %s ]", master.Name, strings.Join(updateList, ", ")))
	}

	updateCount = len(updateList)
	ierr = nil
	return
}

func (intr RSSInteractor) EndRow(updateAt wrappedbasics.WrappedTime) utility.IError {
	master := intr.masterList[intr.index]

	// マスターの更新日付より後の日付が渡されたら更新する
	if updateAt.After(master.LastUpdate) {
		updateData := reference.NewIDWithTime(master.ID, updateAt)
		if result, err := intr.updateMasterRepository.Execute(updateData); err != nil {
			return err.WrapError()
		} else if result.Count == 0 {
			return utility.NewError("RSS Master Update Error. Update Count 0.", "")
		}
		return nil
	}

	return nil
}
