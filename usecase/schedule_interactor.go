package usecase

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/viewschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type ScheduleInteractor struct {
	getRepos                 viewschedule.GetRepository
	participantsDataSplitter string
	participantsArrSplitter  string
	displayScheduleStatus    int
	common                   utility.Common
}

func NewScheduleInteractor(
	getRepos viewschedule.GetRepository,
	participantsDataSplitter, participantsArrSplitter string,
	displayScheduleStatus int,
	common utility.Common,
) ScheduleInteractor {
	return ScheduleInteractor{
		getRepos:                 getRepos,
		participantsDataSplitter: participantsDataSplitter,
		participantsArrSplitter:  participantsArrSplitter,
		common:                   common,
		displayScheduleStatus:    displayScheduleStatus,
	}
}

func (intr ScheduleInteractor) dataAdapter(
	id, platform, url, streamer_name, streamer_id, title, description string,
	status int,
	publish_datetime string,
	duration int,
	thumbnail, icon, participants_data string,
) (domain.ScheduleData, utility.IError) {
	startDate, err := intr.common.CreateNewWrappedTimeFromUTC(publish_datetime)
	if err != nil {
		return domain.ScheduleData{}, err.WrapError()
	}

	newScheduleData, err := domain.NewScheduleData(
		streamer_id,
		streamer_name,
		url,
		status,
		title,
		thumbnail,
		startDate,
		duration,
	)
	if err != nil {
		return domain.ScheduleData{}, err.WrapError()
	}

	if participants_data != "" {
		for _, v := range strings.Split(participants_data, intr.participantsArrSplitter) {
			vSplit := strings.Split(v, intr.participantsDataSplitter)
			mId := vSplit[0]
			mName := vSplit[1]
			mIcon := vSplit[2]
			newScheduleData = newScheduleData.AddParticipants(mId, mName, mIcon)
		}
	}
	return newScheduleData, nil
}

func (intr ScheduleInteractor) GetScheduleData(baseDate utility.WrappedTime) ([]domain.ScheduleData, utility.IError) {
	fromDate := baseDate
	toDate := fromDate.Add(0, 0, 1, 0, 0, 0)
	return intr.getRepos.GetScheduleData(fromDate, toDate, intr.displayScheduleStatus, intr.dataAdapter)
}

func (intr ScheduleInteractor) ToJson(list []domain.ScheduleData) (string, utility.IError) {
	var res []map[string]interface{}
	len := 0
	for _, v := range list {
		m := make(map[string]interface{})
		m["StreamerID"] = v.StreamerID
		m["StreamerName"] = v.StreamerName
		m["VideoLink"] = v.VideoLink
		m["VideoStatus"] = v.VideoStatus.ToInt()
		m["VideoTitle"] = v.VideoTitle
		m["Thumbnail"] = v.Thumbnail
		m["StartDate"] = v.StartDate.ToLocalFormatString()
		m["Duration"] = v.Duration
		m["Participants"] = v.ParticipantsList
		res = append(res, m)
		len++
	}
	d := domain.NewAPIResponseData("ok", len, "", res)
	return d.ToJson()

}
