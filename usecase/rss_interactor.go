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

	isOldItem := func(item utility.IFeedItem) bool {
		feedTime, err := common.CreateNewWrappedTimeFromUTC(item.GetUpdateAt())
		if err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("feed updatetime parse error: %s", item.GetTitle())))
			return true
		}
		return feedTime.Before(master.LastUpdate)
	}

	var oldContents []string
	var seedScheduleList []domain.SeedSchedule
	for i := 0; i < feed.GetItemLength(); i++ {
		feedItem := feed.GetItem(i)
		if feedItem == nil {
			utility.LogError(utility.NewError(fmt.Sprintf("items out of index: %v", i), utility.ERR_RSS_PARSE))
			continue
		}
		if isOldItem(feedItem) {
			oldContents = append(oldContents, feedItem.GetTitle())
			continue
		}
		_, videoId, _ := feedItem.GetExtension(videoIdElement, videoIdItemName)

		seedSchedule := domain.NewSeedSchedule(videoId, platformType, feedStatus)
		seedScheduleList = append(seedScheduleList, seedSchedule)
	}
	if len(oldContents) > 0 {
		utility.LogDebug(fmt.Sprintf("old contents: %s", strings.Join(oldContents, ",")))
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
	utility.LogInfo(fmt.Sprintf("feed request: %s %s", targetMaster.ID, targetMaster.Name))

	list, err = intr.requestRepository.Request(targetMaster.Url, converter)
	if err != nil {
		err = err.WrapError()
	}
	return
}

func (intr RSSInteractor) scheduleIsCompleteConverter(id string, iscomplete int) (string, bool) {
	return id, iscomplete == intr.completeStatus
}

func (intr RSSInteractor) PushToDB(list []domain.SeedSchedule) (insertCount, updateCount int, ierr utility.IError) {
	var idList []string
	for _, seedSchedule := range list {
		idList = append(idList, seedSchedule.GetID())
	}

	withIsComplete, err := intr.getScheduleRepository.Get(idList, intr.platform, intr.scheduleIsCompleteConverter)
	if err != nil {
		ierr = err.WrapError()
		return
	}
	var updateList []string

	for _, seedSchedule := range list {
		if withIsComplete.Has(seedSchedule.GetID()) {
			if withIsComplete.IsComplete(seedSchedule.GetID()) {
				updateList = append(updateList, seedSchedule.GetID())
			}
			continue
		}

		utility.LogDebug(fmt.Sprintf("insert: %s", seedSchedule.GetID()))

		if err = intr.insertScheduleRepository.Insert(seedSchedule); err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("insert error id: %s", seedSchedule.GetID())))
			continue
		}
		insertCount++
	}

	if err = intr.updateScheduleRepository.Update(updateList, intr.platform, intr.notCompleteStatus); err != nil {
		ierr = err.WrapError(fmt.Sprintf("update error id: %s", strings.Join(updateList, ",")))
		return
	}
	updateCount = len(updateList)
	ierr = nil
	return
}

func (intr RSSInteractor) EndRow() utility.IError {
	return intr.updateMasterRepository.UpdateTime(intr.masterList[intr.index].ID, intr.startTime)
}
