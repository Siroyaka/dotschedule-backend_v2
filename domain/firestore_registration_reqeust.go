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

func (frr FirestoreRegistrationRequest) createScheduleID() string {
	switch frr.platform {
	case "YOUTUBE":
		return frr.id
	case "niconico":
		return fmt.Sprintf("niconico_%s", frr.id)
	case "TWITTER_SPACE":
		return fmt.Sprintf("twspace_%s", frr.id)
	default:
		return fmt.Sprintf("OTHERS_%s", frr.id)
	}
}

func (frr FirestoreRegistrationRequest) createPlatformType() string {
	switch frr.platform {
	case "YOUTUBE":
		return "YOUTUBE"
	default:
		return "OTHERS"
	}
}

func (frr FirestoreRegistrationRequest) createStatus() string {
	switch frr.platform {
	default:
		return "100"
	}
}

func (frr FirestoreRegistrationRequest) createYoutubeFullSchedule() {

}

func (frr FirestoreRegistrationRequest) createTWSpaceFullSchedule() {

}

func (frr FirestoreRegistrationRequest) createNiconicoFullSchedule() {

}

func (frr FirestoreRegistrationRequest) createOtherFullSchedule() {

}

func (frr FirestoreRegistrationRequest) CreateFullschedule() {

}

func (frr FirestoreRegistrationRequest) CreateParticipants() {

}
