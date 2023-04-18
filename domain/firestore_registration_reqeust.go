package domain

import (
	"fmt"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type FirestoreRegistrationRequest struct {
	id           string
	participants []string
	url          string
	platform     string
	updateAt     wrappedbasics.IWrappedTime
	startDate    wrappedbasics.IWrappedTime
	streamerID   string
	streamerName string
	title        string
}

func NewFirestoreRegistrationRequestYoutube(
	id string,
	participants []string,
	url string,
	platform string,
	updateAt wrappedbasics.IWrappedTime,
) FirestoreRegistrationRequest {
	return FirestoreRegistrationRequest{
		id:           id,
		participants: participants,
		url:          url,
		platform:     platform,
		updateAt:     updateAt,
	}
}

func NewFirestoreRegistrationRequestOthers(
	id string,
	participants []string,
	url string,
	platform string,
	updateAt wrappedbasics.IWrappedTime,
	startDate wrappedbasics.IWrappedTime,
	streamerID string,
	streamerName string,
	title string,
) FirestoreRegistrationRequest {
	return FirestoreRegistrationRequest{
		id:           id,
		participants: participants,
		url:          url,
		platform:     platform,
		updateAt:     updateAt,
		startDate:    startDate,
		streamerID:   streamerID,
		streamerName: streamerName,
		title:        title,
	}
}

func (frr FirestoreRegistrationRequest) Tostring() string {
	var paramValues []string
	paramValues = append(paramValues, fmt.Sprintf(" \"id\": \"%s\"", frr.id))
	paramValues = append(paramValues, fmt.Sprintf(" \"platform\": \"%s\"", frr.platform))
	paramValues = append(paramValues, fmt.Sprintf(" \"url\": \"%s\"", frr.url))
	paramValues = append(paramValues, fmt.Sprintf(" \"participants\": [%s]", strings.Join(frr.participants, ", ")))
	paramValues = append(paramValues, fmt.Sprintf(" \"updateAt\": \"%s\"", frr.updateAt.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())))

	if frr.platform == "YOUTUBE" {
		return fmt.Sprintf("{ %s }", strings.Join(paramValues, ", "))
	}

	paramValues = append(paramValues, fmt.Sprintf(" \"streamerName\": \"%s\"", frr.streamerName))
	paramValues = append(paramValues, fmt.Sprintf(" \"streamerID\": \"%s\"", frr.streamerID))
	paramValues = append(paramValues, fmt.Sprintf(" \"title\": \"%s\"", frr.title))
	paramValues = append(paramValues, fmt.Sprintf(" \"startDate\": \"%s\"", frr.startDate.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())))

	return fmt.Sprintf("{ %s }", strings.Join(paramValues, ", "))
}

func (frr FirestoreRegistrationRequest) createStreamingID() string {
	switch frr.platform {
	case "YOUTUBE", "niconico", "TWITTER_SPACE", "TWITTER":
		return frr.id
	default:
		return fmt.Sprintf("%s_%s", frr.platform, frr.id)
	}
}

func (frr FirestoreRegistrationRequest) createPlatformType() string {
	switch frr.platform {
	case "YOUTUBE":
		return "YOUTUBE"
	case "niconico":
		return "NICONICO"
	case "TWITTER_SPACE", "TWITTER":
		return "TWITTER"
	default:
		return "OTHERS"
	}
}

func (frr FirestoreRegistrationRequest) createStatus() string {
	switch frr.platform {
	case "YOUTUBE":
		return "20"
	case "niconico":
		return "90"
	case "TWITTER":
		return "90"
	default:
		return "90"
	}
}

func (frr FirestoreRegistrationRequest) CreateFullschedule() FullScheduleData {
	streamingID := frr.createStreamingID()
	platformType := frr.createPlatformType()
	status := frr.createStatus()

	isViewing := 0
	isCompleteData := 0

	if status == "90" {
		isViewing = 1
		isCompleteData = 1
	}

	publishDateTime := ""
	if status == "90" {
		publishDateTime = frr.startDate.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	}

	return FullScheduleData{
		StreamingID:     streamingID,
		PlatformType:    platformType,
		Url:             frr.url,
		StreamerName:    frr.streamerName,
		StreamerID:      frr.streamerID,
		Title:           frr.title,
		Status:          status,
		PublishDatetime: publishDateTime,
		IsViewing:       isViewing,
		IsCompleteData:  isCompleteData,
	}
}

func (frr FirestoreRegistrationRequest) CreateParticipants() StreamingParticipants {
	streamingID := frr.createStreamingID()
	platformType := frr.createPlatformType()
	return NewStreamingParticipants(streamingID, platformType, frr.participants...)
}
