package rssrequest

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type RequestRepository interface {
	Request(url string, converter RequestDataConverter) ([]domain.SeedSchedule, utility.IError)
}

type RequestDataConverter func(data string) ([]domain.SeedSchedule, utility.IError)
