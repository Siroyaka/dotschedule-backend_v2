package interactor

import (
	"fmt"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/dbmodels"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/participants"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type NormalizationYoutubeDataInteractor struct {
	getStreamerMasterRepos abstruct.RepositoryRequest[string, []domain.StreamerMasterWithPlatformData]
	getScheduleRepos       abstruct.RepositoryRequest[reference.VoidStruct, []domain.FullScheduleData]

	updateScheduleRepos            abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse]
	updateScheduleToStatus100Repos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse]

	youtubeVideoListAPIRepos abstruct.RepositoryRequest[string, domain.YoutubeVideoData]

	getParticipantsRepos    abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, []dbmodels.KeyValue[string, string]]
	insertParticipantsRepos abstruct.RepositoryRequest[participants.SingleInsertData, reference.DBUpdateResponse]

	discordPostRepos abstruct.RepositoryRequest[domain.DiscordWebhookParams, string]

	durationParser utility.YoutubeDurationParser
	platformType,
	streamingUrlPrefix string
	discordNortificationRange int
}

func NewNormalizationYoutubeDataInteractor(
	getStreamerMasterRepos abstruct.RepositoryRequest[string, []domain.StreamerMasterWithPlatformData],
	getScheduleRepos abstruct.RepositoryRequest[reference.VoidStruct, []domain.FullScheduleData],

	updateScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse],
	updateScheduleToStatus100Repos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse],
	getParticipantsRepos abstruct.RepositoryRequest[reference.StreamingIDWithPlatformType, []dbmodels.KeyValue[string, string]],
	insertParticipantsRepos abstruct.RepositoryRequest[participants.SingleInsertData, reference.DBUpdateResponse],
	youtubeVideoListAPIRepos abstruct.RepositoryRequest[string, domain.YoutubeVideoData],
	discordPostRepos abstruct.RepositoryRequest[domain.DiscordWebhookParams, string],
	durationParser utility.YoutubeDurationParser,
	platformType,
	streamingUrlPrefix string,
	discordNortificationRange int,
) NormalizationYoutubeDataInteractor {
	return NormalizationYoutubeDataInteractor{
		getStreamerMasterRepos:         getStreamerMasterRepos,
		getScheduleRepos:               getScheduleRepos,
		updateScheduleRepos:            updateScheduleRepos,
		updateScheduleToStatus100Repos: updateScheduleToStatus100Repos,
		getParticipantsRepos:           getParticipantsRepos,
		insertParticipantsRepos:        insertParticipantsRepos,
		youtubeVideoListAPIRepos:       youtubeVideoListAPIRepos,
		discordPostRepos:               discordPostRepos,
		durationParser:                 durationParser,
		platformType:                   platformType,
		streamingUrlPrefix:             streamingUrlPrefix,
		discordNortificationRange:      discordNortificationRange,
	}
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

func (intr NormalizationYoutubeDataInteractor) makeFullSchedule(videoData domain.YoutubeVideoData, beforeScheduleData domain.FullScheduleData, streamerMasterMap map[string]domain.StreamerMaster, updateAt wrappedbasics.IWrappedTime) (domain.FullScheduleData, utility.IError) {
	result := domain.NewEmptyFullScheduleData(beforeScheduleData.StreamingID, intr.platformType)
	status, err := intr.makeStatus(videoData)
	if err != nil {
		return result, err.WrapError()
	}
	result.Status = status

	var publishDateTime wrappedbasics.IWrappedTime
	datetimeFormat := wrappedbasics.WrappedTimeProps.DateTimeFormat()
	switch status {
	case "0":
		publishDateTime, err = wrappedbasics.NewWrappedTimeFromUTC(videoData.Snippet.PublishAt, datetimeFormat)
	case "1", "2":
		publishDateTime, err = wrappedbasics.NewWrappedTimeFromUTC(videoData.LiveStreamingDetails.ActualStartTime, datetimeFormat)
	case "3":
		publishDateTime, err = wrappedbasics.NewWrappedTimeFromUTC(videoData.LiveStreamingDetails.ScheduledStartTime, datetimeFormat)
	}

	if err != nil {
		return result, err.WrapError()
	}
	result.PublishDatetime = publishDateTime.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

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
	result.UpdateAt = updateAt.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	result.Url = fmt.Sprintf("%s%s", intr.streamingUrlPrefix, videoData.Id)

	return result, nil
}

func (intr NormalizationYoutubeDataInteractor) isParticipantsUpdate(data domain.FullScheduleData, participantsMap map[string]string) bool {
	// すでにparticipantsテーブルに対象の配信の配信者idが登録されている
	if _, ok := participantsMap[data.StreamerID]; ok {
		return false
	}

	// streamerIDがない = DBのマスターテーブルに対象配信者のyoutubeIDが登録されていない = participantsに登録する必要がない(できない)
	if data.StreamerID == "" {
		return false
	}

	return true
}

// # discord通知する条件
//
//   - 初めてYoutubeからデータ取得した対象であること
//
//   - 元のステータスが10 or 20のものを初めて取得する対象と定義する
//
//   - 動画の場合は以下に条件はなし。配信の場合は以下の条件を満たすこと
//
//   - 配信の開始時間から指定時間以上経過していないこと(時間は設定ファイルで定義)
func (intr NormalizationYoutubeDataInteractor) isDiscordNortification(beforeScheduleData domain.FullScheduleData, afterScheduleData domain.FullScheduleData) bool {
	if beforeScheduleData.Status != "10" && beforeScheduleData.Status != "20" {
		return false
	}

	// 動画
	if afterScheduleData.Status == "0" {
		return true
	}

	now := wrappedbasics.Now(wrappedbasics.WrappedTimeProps.LocalLocation())

	startDate, err := wrappedbasics.NewWrappedTimeFromUTC(afterScheduleData.PublishDatetime, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		utility.LogError(err.WrapError())
		return false
	}

	deadLine := startDate.Add(0, 0, 0, 0, intr.discordNortificationRange, 0)

	if now.After(deadLine) {
		return false
	}

	utility.LogDebug(fmt.Sprintf("Discord nortification. { \"title\": \"%s\"}", afterScheduleData.Title))

	return true
}

// discordへ投げるデータを作成する
func (intr NormalizationYoutubeDataInteractor) CreateDiscordPostData(afterScheduleData domain.FullScheduleData, participantsMaps map[string]string) domain.DiscordWebhookParams {
	var members []string
	if afterScheduleData.StreamerID != "" {
		members = append(members, afterScheduleData.StreamerName)
	}

	for k, v := range participantsMaps {
		if afterScheduleData.StreamerID == k {
			continue
		}
		members = append(members, v)
	}

	embedDescription := fmt.Sprintf("参加者: %s", strings.Join(members, "、"))

	embed := domain.DiscordWebhookEmbed{
		Title:       afterScheduleData.Title,
		TimeStamp:   afterScheduleData.PublishDatetime,
		Url:         afterScheduleData.Url,
		Description: embedDescription,
	}

	content := fmt.Sprintf("%sの配信が追加されたよ", afterScheduleData.StreamerName)

	return domain.DiscordWebhookParams{
		Content: content,
		Embeds:  []domain.DiscordWebhookEmbed{embed},
	}
}

func (intr *NormalizationYoutubeDataInteractor) Normalization() utility.IError {
	streamerMasterSeedData, err := intr.getStreamerMasterRepos.Execute(intr.platformType)
	if err != nil {
		return err.WrapError()
	}

	// Need to make streamer data from youtube channel id. Create map data for that
	streamerMasterMap := intr.createStreamerDataMap(streamerMasterSeedData)

	targetSchedules, err := intr.getScheduleRepos.Execute(reference.Void())
	if err != nil {
		return err.WrapError()
	}

	for _, data := range targetSchedules {
		now := wrappedbasics.Now(wrappedbasics.WrappedTimeProps.LocalLocation())

		utility.LogInfo(fmt.Sprintf("target id = %s", data.StreamingID))

		youtubeVideoData, err := intr.youtubeVideoListAPIRepos.Execute(data.StreamingID)
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		if youtubeVideoData.IsEmpty() {
			utility.LogInfo(fmt.Sprintf("notfound from youtube data api. change status to 100. { \"StreamingID\": \"%s\" }", data.StreamingID))

			if updateResult, err := intr.updateScheduleToStatus100Repos.Execute(data); err != nil {
				utility.LogError(err.WrapError(fmt.Sprintf("{\"StreamingID\": \"%s\"}", data.StreamingID)))
				continue
			} else if updateResult.Count == 0 {
				utility.LogError(utility.NewError(fmt.Sprintf("schedule update to 100 count 0. { \"StreamingID\": \"%s\"}", data.StreamingID), ""))
				continue
			}
			utility.LogInfo(fmt.Sprintf("Finished change status to 100. { \"StreamingID\": \"%s\" }", data.StreamingID))
			continue
		}

		// youtube video data to fullschedule data
		afterScheduleData, err := intr.makeFullSchedule(youtubeVideoData, data, streamerMasterMap, now)
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		if (data.Status == "2" || data.Status == "3") && data.Status == afterScheduleData.Status {
			// status 2 or 3 ... status not change that not update
			utility.LogInfo(fmt.Sprintf("data status not change. not update. { id: %s, streamer_name: %s, title: %s }", afterScheduleData.StreamingID, afterScheduleData.StreamerName, afterScheduleData.Title))
			continue
		}

		utility.LogInfo(fmt.Sprintf("schedule update. { id: %s, streamer_name: %s, title: %s }", afterScheduleData.StreamingID, afterScheduleData.StreamerName, afterScheduleData.Title))

		updateResult, err := intr.updateScheduleRepos.Execute(afterScheduleData)
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}
		if updateResult.Count == 0 {
			utility.LogError(utility.NewError(fmt.Sprintf("update count = 0, id = %s", data.StreamingID), ""))
			continue
		}

		streamingIdWithPlatformType := reference.NewStreamingIDWithPlatformType(data.StreamingID, intr.platformType)

		participantsIdNames, err := intr.getParticipantsRepos.Execute(streamingIdWithPlatformType)
		if err != nil {
			utility.LogError(err.WrapError())
			continue
		}

		participantsIdNameMap := dbmodels.KeyValueToMap(participantsIdNames)

		if intr.isParticipantsUpdate(afterScheduleData, participantsIdNameMap) {
			utility.LogInfo(fmt.Sprintf("participants data insert. streamerID = %s", afterScheduleData.StreamerID))

			participantsSingleInsertData := participants.NewSingleInsertData(data.StreamingID, intr.platformType, afterScheduleData.StreamerID, now)

			if response, err := intr.insertParticipantsRepos.Execute(participantsSingleInsertData); err != nil {
				utility.LogError(err.WrapError())
			} else if response.Count == 0 {
				utility.LogError(utility.NewError(fmt.Sprintf("participants update count = 0, id = %s", data.StreamingID), ""))
			}

		}

		// discordへの通知 discordへの通知対象であるかを判定し、通知対象ならば通知を行う
		if intr.isDiscordNortification(data, afterScheduleData) {
			// discordへの通知用データを作成する
			discordPostData := intr.CreateDiscordPostData(afterScheduleData, participantsIdNameMap)

			if message, err := intr.discordPostRepos.Execute(discordPostData); err != nil {
				utility.LogError(err.WrapError(message))
			}
		}

	}

	return nil
}
