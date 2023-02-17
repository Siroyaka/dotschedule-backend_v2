package reference

import "github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"

type IDWithTime struct {
	id   string
	time wrappedbasics.IWrappedTime
}

func NewIDWithTime(id string, time wrappedbasics.IWrappedTime) IDWithTime {
	return IDWithTime{
		id:   id,
		time: time,
	}
}

func (iwt IDWithTime) Id() string {
	return iwt.id
}

func (iwt IDWithTime) Time() wrappedbasics.IWrappedTime {
	return iwt.time
}
