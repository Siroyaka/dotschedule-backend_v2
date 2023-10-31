package httprequest

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type RSSRequestRepository struct {
	request abstruct.HTTPRequest
}

func NewRSSRequestRepository(request abstruct.HTTPRequest) RSSRequestRepository {
	return RSSRequestRepository{
		request: request,
	}
}

func (repos RSSRequestRepository) converter(feedText string) (utility.IFeed, utilerror.IError) {
	rssParser := utility.NewRSSParser(string(wrappedbasics.WrappedTimeProps.DateTimeFormat()))
	feed, err := rssParser.Parse(feedText)
	if err != nil {
		return nil, err.WrapError()
	}
	return feed, nil
}

func (repos RSSRequestRepository) Execute(url string) (utility.IFeed, utilerror.IError) {
	response, err := repos.request.Get(url)
	if err != nil {
		return nil, err.WrapError("rss data get falied.")
	}
	logger.Debug(fmt.Sprintf("http request: {url: %s, status:%s, statusCode:%b}", url, response.Status(), response.StatusCode()))

	body := response.Body()
	if body == "" {
		return nil, utilerror.New("RSS responsebody has not content.", utilerror.ERR_RSS_PARSE)
	}

	feedData, err := repos.converter(response.Body())
	if err != nil {
		return nil, err.WrapError(fmt.Sprintf("response: { %s }", response.Body()))
	}

	return feedData, nil
}
