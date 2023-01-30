package usecase

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	rssmaster "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/master"
	rssrequest "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/request"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type RSSInteractor struct {
	index      int
	startTime  utility.WrappedTime
	masterList []domain.RSSMaster

	common utility.Common

	getMasterRepository    rssmaster.GetRepository
	updateMasterRepository rssmaster.UpdateRepository

	requestRepository rssrequest.RequestRepository

	getScheduleRepository    rssschedule.GetRepository
	insertScheduleRepository rssschedule.InsertRepository
	updateScheduleRepository rssschedule.UpdateRepository

	completeStatus, notCompleteStatus                       int
	platform, insertStatus, videoIdElement, videoIdItemName string
}

func NewRSSInteractor(
	common utility.Common,
	getMasterRepository rssmaster.GetRepository,
	updateMasterRepository rssmaster.UpdateRepository,
	requestRepository rssrequest.RequestRepository,
	getScheduleRepository rssschedule.GetRepository,
	insertScheduleRepository rssschedule.InsertRepository,
	updateScheduleRepository rssschedule.UpdateRepository,
	completeStatus, notCompleteStatus int,
	platform, insertStatus, videoIdElement, videoIdItemName string,
) RSSInteractor {
	startTime, err := common.Now()
	if err != nil {
		utility.LogFatal(err.WrapError())
		panic(err)
	}
	var intr = RSSInteractor{
		index:                    -1,
		startTime:                startTime,
		masterList:               []domain.RSSMaster{},
		common:                   common,
		getMasterRepository:      getMasterRepository,
		updateMasterRepository:   updateMasterRepository,
		requestRepository:        requestRepository,
		getScheduleRepository:    getScheduleRepository,
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

func (intr RSSInteractor) masterDataAdapter(id, name, url, date string) (master domain.RSSMaster, err utility.IError) {
	t, err := intr.common.CreateNewWrappedTimeFromUTC(date)
	return domain.NewRSSMaster(id, name, url, t), err
}

func (intr *RSSInteractor) GetMaster() utility.IError {
	list, err := intr.getMasterRepository.Get(intr.masterDataAdapter)
	intr.masterList = list
	return err
}

func (intr *RSSInteractor) Next() bool {
	if len(intr.masterList) > intr.index {
		intr.index++
	}
	return len(intr.masterList) > intr.index
}

func (intr RSSInteractor) rssDataConverter(data, platformType, feedStatus, videoIdElement, videoIdItemName string, common utility.Common, master domain.RSSMaster) ([]domain.SeedSchedule, utility.IError) {
	rssParser := utility.NewRSSParser(common.TimeFormat())
	feed, err := rssParser.Parse(data)
	if err != nil {
		return nil, err.WrapError()
	}

	var seedScheduleList []domain.SeedSchedule
	for i := 0; i < feed.GetItemLength(); i++ {
		feedItem := feed.GetItem(i)
		if feedItem == nil {
			utility.LogError(utility.NewError(fmt.Sprintf("items out of index. index: %v", i), utility.ERR_RSS_PARSE))
			continue
		}

		feedPublishedAt, err := common.CreateNewWrappedTimeFromUTC(feedItem.GetPublishedAt())
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

	return seedScheduleList, nil
}

func (intr RSSInteractor) GetRSSData() (list []domain.SeedSchedule, err utility.IError) {
	if len(intr.masterList) <= intr.index {
		err = utility.NewError("", utility.ERR_OUTOFINDEX, "ri.masterList", strconv.Itoa(len(intr.masterList)), strconv.Itoa(intr.index))
		return
	}
	targetMaster := intr.masterList[intr.index]

	converter := func(feedText string) (list []domain.SeedSchedule, err utility.IError) {
		return intr.rssDataConverter(feedText, intr.platform, intr.insertStatus, intr.videoIdElement, intr.videoIdItemName, intr.common, targetMaster)
	}
	utility.LogDebug(fmt.Sprintf("feed request: %s %s", targetMaster.ID, targetMaster.Name))

	list, err = intr.requestRepository.Request(targetMaster.Url, converter)
	if err != nil {
		err = err.WrapError()
	}
	return
}

func (intr RSSInteractor) scheduleIsCompleteConverter(id string, iscomplete int) (string, bool) {
	return id, iscomplete == intr.completeStatus
}

func (intr RSSInteractor) PushToDB(list []domain.SeedSchedule) (insertCount, updateCount int, isError bool, newestPublishedAt utility.WrappedTime, ierr utility.IError) {
	var idList []string
	isError = false

	master := intr.masterList[intr.index]
	newestPublishedAt = master.LastUpdate

	for _, seedSchedule := range list {
		idList = append(idList, seedSchedule.GetID())
	}

	withIsComplete, err := intr.getScheduleRepository.Get(idList, intr.platform, intr.scheduleIsCompleteConverter)
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

		if withIsComplete.Has(seedSchedule.GetID()) {

			// dbにあるデータがcomplete状態であれば更新する、そうでなければ更新もしない
			if withIsComplete.IsComplete(seedSchedule.GetID()) {
				updateList = append(updateList, seedSchedule.GetID())
			}
			continue
		}

		if err = intr.insertScheduleRepository.Insert(seedSchedule); err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("schedule insert error. name: %s, streaming_id: [ %s ]", master.Name, seedSchedule.GetID())))
			isError = true
			continue
		}

		utility.LogInfo(fmt.Sprintf("schedule insert finished. name: %s, streaming_id: [ %s ]", master.Name, seedSchedule.GetID()))

		insertCount++
	}

	if err = intr.updateScheduleRepository.Update(updateList, intr.platform, intr.notCompleteStatus); err != nil {
		ierr = err.WrapError(fmt.Sprintf("schedule update ERROR. name: %s, streaming_id: [ %s ]", master.Name, strings.Join(updateList, ", ")))
		isError = true
		return
	}

	if updateList != nil {
		utility.LogInfo(fmt.Sprintf("schedule update finished. name: %s, streaming_id: [ %s ]", master.Name, strings.Join(updateList, ", ")))
	}

	updateCount = len(updateList)
	ierr = nil
	return
}

func (intr RSSInteractor) EndRow(updateAt utility.WrappedTime) utility.IError {
	master := intr.masterList[intr.index]

	// マスターの更新日付より後の日付が渡されたら更新する
	if updateAt.After(master.LastUpdate) {
		return intr.updateMasterRepository.UpdateTime(master.ID, updateAt)
	}

	return nil
}
