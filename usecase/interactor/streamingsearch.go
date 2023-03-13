package interactor

import (
	"regexp"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type StreamingSearchInteractor struct {
	searchRepos abstruct.RepositoryRequest[apireference.StreamingSearchValues, []apireference.ScheduleResponse]
	countRepos  abstruct.RepositoryRequest[apireference.StreamingSearchValues, int]
}

func NewStreamingSearchInteractor(
	searchRepos abstruct.RepositoryRequest[apireference.StreamingSearchValues, []apireference.ScheduleResponse],
	countRepos abstruct.RepositoryRequest[apireference.StreamingSearchValues, int],
) StreamingSearchInteractor {
	return StreamingSearchInteractor{
		searchRepos: searchRepos,
		countRepos:  countRepos,
	}
}

func (intr StreamingSearchInteractor) CreateValue(
	members string,
	from wrappedbasics.IWrappedTime,
	to wrappedbasics.IWrappedTime,
	tags string,
	page int,
	titleValue string,
) (apireference.StreamingSearchValues, utilerror.IError) {
	textRegex := regexp.MustCompile("([0-9]|[a-z])*")

	var memberIDList []string

	if len(members) > 0 {
		memberIDList = strings.Split(members, ",")
		for _, id := range memberIDList {
			if !textRegex.MatchString(id) {
				return apireference.EmptyStreamingSearchValues(), utilerror.New("id regex error", "")
			}
		}
	}

	var tagList []string

	if len(tags) > 0 {
		tagList = strings.Split(tags, ",")
		for _, tag := range tagList {
			if !textRegex.MatchString(tag) {
				return apireference.EmptyStreamingSearchValues(), utilerror.New("tag regex error", "")
			}
		}
	}

	if from != nil && to != nil {
		if from.After(to) {
			return apireference.EmptyStreamingSearchValues(), utilerror.New("date from over to", "")
		}
	}

	return apireference.NewStreamingSearchValues(memberIDList, from, to, tagList, page, titleValue), nil
}

func (intr StreamingSearchInteractor) Count(searchValue apireference.StreamingSearchValues) (int, utilerror.IError) {
	searchLength, err := intr.countRepos.Execute(searchValue)
	if err != nil {
		return 0, err.WrapError()
	}
	if searchLength == 0 {
		return 0, nil
	}
	return searchLength, nil
}

func (intr StreamingSearchInteractor) Search(searchValue apireference.StreamingSearchValues) ([]apireference.ScheduleResponse, utilerror.IError) {
	resultList, err := intr.searchRepos.Execute(searchValue)
	if err != nil {
		return nil, err.WrapError()
	}
	return resultList, nil
}

func (intr StreamingSearchInteractor) ToJson(data []apireference.ScheduleResponse, count int) (string, utilerror.IError) {
	d := domain.NewAPIResponseData("ok", count, "", data)
	result, err := d.ToJson()
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
