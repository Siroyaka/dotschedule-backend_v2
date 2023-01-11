package domain

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type ScheduleData struct {
	StreamerName, StreamerID, VideoLink string
	VideoStatus                         VideoStatus
	VideoTitle, Thumbnail               string
	StartDate                           utility.WrappedTime
	Duration                            int
	ParticipantsList                    []ParticipantsData
}

func NewScheduleData(
	streamerId, streamerName, videoLink string,
	videoStatus int,
	videoTitle, thumbnail string,
	startDate utility.WrappedTime,
	duration int,
) (ScheduleData, utility.IError) {
	vs, err := NewVideoStatus(videoStatus)
	if err != nil {
		return ScheduleData{}, err
	}
	return ScheduleData{
		StreamerID:       streamerId,
		StreamerName:     streamerName,
		VideoLink:        videoLink,
		VideoStatus:      vs,
		VideoTitle:       videoTitle,
		Thumbnail:        thumbnail,
		StartDate:        startDate,
		Duration:         duration,
		ParticipantsList: []ParticipantsData{},
	}, nil
}

func (sd ScheduleData) AddParticipants(id, name, icon string) ScheduleData {
	sd.ParticipantsList = append(sd.ParticipantsList, ParticipantsData{Id: id, Name: name, Icon: icon})
	return sd
}

type ParticipantsData struct {
	Id, Name, Icon string
}

type VideoStatus int

func NewVideoStatus(s int) (VideoStatus, utility.IError) {
	switch s {
	case 0, 1, 2, 3, 10, 20, 100:
		return VideoStatus(s), nil
	default:
		return VideoStatus(0), utility.NewError("", utility.ERR_INVALIDVALUE, "VideoStatus", s)
	}
}

func (vs VideoStatus) ToInt() int {
	return int(vs)
}