package interactor

import (
	"database/sql"
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type NormalizationYoutubeDataInteractor struct {
	getStreamerMasterRepos         sqlwrapper.SelectRepository[domain.StreamerMasterWithPlatformData]
	getScheduleRepos               sqlwrapper.SelectRepository[domain.FullScheduleData]
	updateScheduleRepos            sqlwrapper.UpdateRepository
	updateScheduleToStatus100Repos sqlwrapper.UpdateRepository
	countParticipantsRepos         sqlwrapper.SelectRepository[int]
	insertParticipantsRepos        sqlwrapper.UpdateRepository
	youtubeVideoListAPIRepos       youtubedataapi.VideoListRepository
	common                         utility.Common
	durationParser                 utility.YoutubeDurationParser
	platformType,
	streamingUrlPrefix string
	partList []string
}

func NewNormalizationYoutubeDataInteractor(
	getStreamerMasterRepos sqlwrapper.SelectRepository[domain.StreamerMasterWithPlatformData],
	getScheduleRepos sqlwrapper.SelectRepository[domain.FullScheduleData],
	updateScheduleRepos sqlwrapper.UpdateRepository,
	updateScheduleToStatus100Repos sqlwrapper.UpdateRepository,
	countParticipantsRepos sqlwrapper.SelectRepository[int],
	insertParticipantsRepos sqlwrapper.UpdateRepository,
	youtubeVideoListAPIRepos youtubedataapi.VideoListRepository,
	common utility.Common,
	durationParser utility.YoutubeDurationParser,
	platformType,
	streamingUrlPrefix string,
	partList []string,
) NormalizationYoutubeDataInteractor {
	return NormalizationYoutubeDataInteractor{
		getStreamerMasterRepos:         getStreamerMasterRepos,
		getScheduleRepos:               getScheduleRepos,
		updateScheduleRepos:            updateScheduleRepos,
		updateScheduleToStatus100Repos: updateScheduleToStatus100Repos,
		countParticipantsRepos:         countParticipantsRepos,
		insertParticipantsRepos:        insertParticipantsRepos,
		youtubeVideoListAPIRepos:       youtubeVideoListAPIRepos,
		common:                         common,
		durationParser:                 durationParser,
		platformType:                   platformType,
		streamingUrlPrefix:             streamingUrlPrefix,
		partList:                       partList,
	}
}

func (intr NormalizationYoutubeDataInteractor) streamerMasterScan(s sqlwrapper.IScan) (domain.StreamerMasterWithPlatformData, utility.IError) {
	var streamer_id, platform_id, streamer_name string

	if err := s.Scan(&streamer_id, &platform_id, &streamer_name); err != nil {
		return domain.NewStreamerMasterWithPlatformData(""), utility.NewError(err.Error(), "")
	}

	res := domain.NewStreamerMasterWithPlatformData(streamer_id)
	res.StreamerName = streamer_name
	res.PlatformData[intr.platformType] = domain.StreamerPlatformMaster{
		StreamerID:   streamer_id,
		PlatformID:   platform_id,
		PlatformType: intr.platformType,
	}

	return res, nil
}

func (intr NormalizationYoutubeDataInteractor) createStreamerDataMap(list []domain.StreamerMasterWithPlatformData) map[string]domain.StreamerMaster {
	res := make(map[string]domain.StreamerMaster)
	for _, data := range list {
		v, ok := data.PlatformData[intr.platformType]
		if !ok {
			utility.LogError(utility.NewError(fmt.Sprintf("%s has not %s data", data.StreamerID, intr.platformType), ""))
			continue
		}
		res[v.PlatformID] = data.StreamerMaster
	}
	return res
}

func (intr NormalizationYoutubeDataInteractor) scheduleScan(s sqlwrapper.IScan) (domain.FullScheduleData, utility.IError) {
	var streaming_id, status string
	var publish_datetime sql.NullString

	if err := s.Scan(&streaming_id, &status, &publish_datetime); err != nil {
		return domain.NewEmptyFullScheduleData("", intr.platformType), utility.NewError(err.Error(), "")
	}

	res := domain.NewEmptyFullScheduleData(streaming_id, intr.platformType)
	res.Status = status
	if !publish_datetime.Valid {
		res.PublishDatetime = publish_datetime.String
	}

	return res, nil
}

func (intr *NormalizationYoutubeDataInteractor) status100Data(schedule domain.FullScheduleData) ([]interface{}, utility.IError) {
	now, err := intr.common.Now()
	if err != nil {
		return utility.ToInterfaceSlice(), err.WrapError("create 'now' error")
	}

	afterStatus := 100
	isViewing := 0
	isComplete := 0

	if schedule.PublishDatetime != "" {
		beforePublishDatetime, err := intr.common.CreateNewWrappedTimeFromUTC(schedule.PublishDatetime)

		if err != nil {
			return utility.ToInterfaceSlice(), err.WrapError("publishDateTime parse error")
		}

		if beforePublishDatetime.Before(now.Add(0, 0, -14, 0, 0, 0)) {
			// target schedule data is complete
			isComplete = 1
		}
	}

	result := utility.ToInterfaceSlice(
		now.ToUTCFormatString(),
		afterStatus,
		isViewing,
		isComplete,
		schedule.StreamingID,
		intr.platformType,
	)

	return result, nil
}

func (intr *NormalizationYoutubeDataInteractor) updateStatusTo100(data domain.FullScheduleData) utility.IError {
	queryValues, err := intr.status100Data(data)
	if err != nil {
		return err.WrapError("status update to 100 failed")
	}
	count, _, err := intr.updateScheduleToStatus100Repos.UpdatePrepare(queryValues)
	if err != nil {
		return err.WrapError("status update to 100 failed")
	}
	if count == 0 {
		return utility.NewError("status update to 100 failed", "")
	}
	return nil
}

func (intr NormalizationYoutubeDataInteractor) makeStatus(videoData domain.YoutubeVideoData) (string, utility.IError) {
	if videoData.LiveStreamingDetails.IsEmpty() {
		// this status is 0: not streaming. this is video.
		return "0", nil
	}

	if videoData.Snippet.IsEmpty() {
		return "100", utility.NewError(fmt.Sprintf("snippet is not found: %s", videoData.Id), "")
	}

	switch videoData.Snippet.LiveBroadcastContent {
	case "none":
		// this status is 1: streaming is already finished.
		return "1", nil
	case "live":
		// this status is 2: streaming now.
		return "2", nil
	case "upcoming":
		// this status is 3: streaming is upcoming.
		return "3", nil
	default:
		return "100", utility.NewError(fmt.Sprintf("id: %s, ", videoData.Id), "")
	}
}

func (intr NormalizationYoutubeDataInteractor) makeFullSchedule(videoData domain.YoutubeVideoData, beforeScheduleData domain.FullScheduleData, streamerMasterMap map[string]domain.StreamerMaster, updateAt utility.WrappedTime) (domain.FullScheduleData, utility.IError) {
	result := domain.NewEmptyFullScheduleData(beforeScheduleData.StreamingID, intr.platformType)
	status, err := intr.makeStatus(videoData)
	if err != nil {
		return result, err.WrapError()
	}
	result.Status = status

	var publishDateTime utility.WrappedTime
	switch status {
	case "0":
		publishDateTime, err = intr.common.CreateNewWrappedTimeFromUTC(videoData.Snippet.PublishAt)
	case "1", "2":
		publishDateTime, err = intr.common.CreateNewWrappedTimeFromUTC(videoData.LiveStreamingDetails.ActualStartTime)
	case "3":
		publishDateTime, err = intr.common.CreateNewWrappedTimeFromUTC(videoData.LiveStreamingDetails.ScheduledStartTime)
	}

	if err != nil {
		return result, err.WrapError()
	}
	result.PublishDatetime = publishDateTime.ToUTCFormatString()

	duration := 0
	if status == "0" || status == "1" {
		parserResult := intr.durationParser.Set(videoData.ContentDetails.Duration)
		if err := parserResult.Err(); err != nil {
			return result, err.WrapError()
		}
		duration = parserResult.GetTotalSeconds()
	}
	result.Duration = duration

	if v, ok := streamerMasterMap[videoData.Snippet.ChannelID]; ok {
		result.StreamerID = v.StreamerID
		result.StreamerName = v.StreamerName
	} else {
		result.StreamerName = videoData.Snippet.ChannelName
	}

	if status == "0" || status == "1" {
		result.IsViewing = 1
		result.IsCompleteData = 1
	} else {
		result.IsViewing = 1
		result.IsCompleteData = 0
	}

	result.Title = videoData.Snippet.Title
	result.ThumbnailLink = videoData.Snippet.Thumbnail
	result.Description = videoData.Snippet.Description
	result.UpdateAt = updateAt.ToUTCFormatString()
	result.Url = fmt.Sprintf("%s%s", intr.streamingUrlPrefix, videoData.Id)

	return result, nil
}

func (intr NormalizationYoutubeDataInteractor) participantsCounter(s sqlwrapper.IScan) (int, utility.IError) {
	var count int
	if err := s.Scan(&count); err != nil {
		return 0, utility.NewError(err.Error(), "")
	}
	return count, nil
}

func (intr *NormalizationYoutubeDataInteractor) Normalization() utility.IError {
	streamerMasterSeedData, err := intr.getStreamerMasterRepos.Select(intr.streamerMasterScan)
	if err != nil {
		return err.WrapError()
	}

	// Need to make streamer data from youtube channel id. Create map data for that
	streamerMasterMap := intr.createStreamerDataMap(streamerMasterSeedData)

	targetSchedules, err := intr.getScheduleRepos.Select(intr.scheduleScan)
	if err != nil {
		return err.WrapError()
	}

	for _, data := range targetSchedules {
		now, err := intr.common.Now()
		if err != nil {
			utility.LogFatal(err.WrapError())
			return err
		}

		utility.LogInfo(fmt.Sprintf("target id = %s", data.StreamingID))
		targetStreamingID := data.StreamingID

		apiData, err := intr.youtubeVideoListAPIRepos.IdSearch(intr.partList, []string{targetStreamingID})
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		youtubeVideoData := domain.NewEmptyYoutubeVideoData()
		if len(apiData) != 0 {
			youtubeVideoData = apiData[0]
		}

		if youtubeVideoData.IsEmpty() {
			utility.LogInfo("data is notfound from youtube data api. status to 100.")
			if err := intr.updateStatusTo100(data); err != nil {
				utility.LogError(err.WrapError())
			}
			continue
		}

		// youtube video data to fullschedule data
		afterScheduleData, err := intr.makeFullSchedule(youtubeVideoData, data, streamerMasterMap, now)
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		if (data.Status == "2" || data.Status == "3") && data.Status == afterScheduleData.Status {
			utility.LogInfo("data status not change. not update.")
			// status 2 or 3 ... status not change that not update
			continue
		}

		utility.LogInfo(fmt.Sprintf("schedule update. id = %s, streamer_name = %s, title = %s", afterScheduleData.StreamingID, afterScheduleData.StreamerName, afterScheduleData.Title))
		count, _, err := intr.updateScheduleRepos.UpdatePrepare(utility.ToInterfaceSlice(
			afterScheduleData.Url,
			afterScheduleData.StreamerName,
			afterScheduleData.StreamerID,
			afterScheduleData.Title,
			afterScheduleData.Description,
			afterScheduleData.Status,
			afterScheduleData.PublishDatetime,
			afterScheduleData.Duration,
			afterScheduleData.ThumbnailLink,
			afterScheduleData.UpdateAt,
			afterScheduleData.IsViewing,
			afterScheduleData.IsCompleteData,
			data.StreamingID,
			intr.platformType,
		))
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}
		if count == 0 {
			utility.LogError(utility.NewError(fmt.Sprintf("update count = 0, id = %s", data.StreamingID), ""))
			continue
		}

		if afterScheduleData.StreamerID == "" {
			continue
		}

		countList, err := intr.countParticipantsRepos.SelectPrepare(intr.participantsCounter, utility.ToInterfaceSlice(data.StreamingID, intr.platformType, afterScheduleData.StreamerID))
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		if len(countList) == 0 || countList[0] != 0 {
			continue
		}

		utility.LogInfo(fmt.Sprintf("participants data insert. streamerID = %s", afterScheduleData.StreamerID))
		count, _, err = intr.insertParticipantsRepos.UpdatePrepare(utility.ToInterfaceSlice(data.StreamingID, intr.platformType, afterScheduleData.StreamerID, now.ToUTCFormatString()))
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		if count == 0 {
			utility.LogError(utility.NewError(fmt.Sprintf("participants update count = 0, id = %s", data.StreamingID), ""))
		}

	}

	return nil
}
