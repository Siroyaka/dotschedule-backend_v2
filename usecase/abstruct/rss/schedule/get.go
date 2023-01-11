package rssschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository interface {
	Get([]string, string, ScheduleIsCompleteConverter) (domain.ScheduleIDWithIsComplete, utility.IError)
}

type ScheduleIsCompleteConverter func(string, int) (string, bool)
