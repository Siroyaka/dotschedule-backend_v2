package rssrequest

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	rssrequest "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/request"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type RequestRepository struct {
	request abstruct.HTTPRequest
}

func NewRequestRepository(request abstruct.HTTPRequest) rssrequest.RequestRepository {
	return RequestRepository{
		request: request,
	}
}

func (r RequestRepository) Request(url string, converter rssrequest.RequestDataConverter) ([]domain.SeedSchedule, utility.IError) {
	a, err := r.request.Get(url)
	if err != nil {
		return nil, err.WrapError()
	}
	utility.LogDebug(fmt.Sprintf("http request: {url: %s, status:%s, statusCode:%b}", url, a.Status(), a.StatusCode()))

	return converter(a.Body())
}
