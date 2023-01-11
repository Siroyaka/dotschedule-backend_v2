package domain

type FullScheduleData struct {
	StreamingID     string
	PlatformType    string
	Url             string
	StreamerName    string
	StreamerID      string
	Title           string
	Description     string
	Status          string
	PublishDatetime string
	Duration        int
	ThumbnailLink   string
	InsertAt        string
	UpdateAt        string
	IsViewing       int
	IsCompleteData  int
}

func NewEmptyFullScheduleData(streamingID, platformType string) FullScheduleData {
	return FullScheduleData{
		StreamingID:  streamingID,
		PlatformType: platformType,
	}
}

func (fsd FullScheduleData) ImportStatusFromFirestore(statusInFirestore int) FullScheduleData {
	status, isViewing, isCompleteData := scheduleStatusExchange(statusInFirestore)
	fsd.Status = status
	fsd.IsViewing = isViewing
	fsd.IsCompleteData = isCompleteData
	return fsd
}

func scheduleStatusExchange(statusInFirestore int) (status string, isViewing int, isCompleteData int) {
	switch statusInFirestore {
	case 0:
		status = "10"
		isViewing = 0
		isCompleteData = 0
		return
	case 1:
		status = "3"
		isViewing = 1
		isCompleteData = 0
		return
	case 2:
		status = "2"
		isViewing = 1
		isCompleteData = 0
		return
	case 3:
		status = "1"
		isViewing = 1
		isCompleteData = 1
		return
	case 4:
		status = "0"
		isViewing = 1
		isCompleteData = 1
		return
	case -1:
		status = "100"
		isViewing = 0
		isCompleteData = 0
		return
	case 10:
		status = "20"
		isViewing = 0
		isCompleteData = 0
		return
	case 100:
		status = "100"
		isViewing = 0
		isCompleteData = 1
		return
	default:
		status = "100"
		isViewing = 0
		isCompleteData = 0
		return
	}
}

type FullScheduleWithPlatformParticipantsData struct {
	FullScheduleData
	PlatformIdList []string
}

func NewEmptyFullScheduleWithPlatformParticipantsData(streamingID, platformType string) FullScheduleWithPlatformParticipantsData {
	return FullScheduleWithPlatformParticipantsData{
		FullScheduleData: NewEmptyFullScheduleData(streamingID, platformType),
	}

}

func (ppd FullScheduleWithPlatformParticipantsData) ImportStatusFromFirestore(statusInFirestore int) FullScheduleWithPlatformParticipantsData {
	status, isViewing, isCompleteData := scheduleStatusExchange(statusInFirestore)
	ppd.FullScheduleData.Status = status
	ppd.FullScheduleData.IsViewing = isViewing
	ppd.FullScheduleData.IsCompleteData = isCompleteData
	return ppd
}

func (ppd FullScheduleWithPlatformParticipantsData) AppendParticipants(idList ...string) FullScheduleWithPlatformParticipantsData {
	ppd.PlatformIdList = append(ppd.PlatformIdList, idList...)
	return ppd
}

type GroupStreamerData struct {
	StreamerID,
	StreamerName string
}
