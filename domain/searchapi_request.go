package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

const (
	defaultPage      = 1
	defaultMaxResult = 20
)

type SearchAPIRequestQueries struct {
	member,
	from,
	to,
	tag,
	page,
	title,
	maxresult string
}

type SearchAPIRequestParams struct {
	member          []string
	from            wrappedbasics.IWrappedTime
	to              wrappedbasics.IWrappedTime
	tags            []string
	page, maxResult int
	title           string
}

func NewSearchAPIRequestQueries(
	member,
	from,
	to,
	tag,
	page,
	title,
	maxresult string,
) SearchAPIRequestQueries {
	return SearchAPIRequestQueries{
		member:    member,
		from:      from,
		to:        to,
		tag:       tag,
		page:      page,
		title:     title,
		maxresult: maxresult,
	}
}

func (q SearchAPIRequestQueries) Convert() (SearchAPIRequestParams, utilerror.IError) {
	textRegex := regexp.MustCompile("([0-9]|[a-z])*")

	var memberList []string

	if len(q.member) > 0 {
		memberList = strings.Split(q.member, ",")
		for _, id := range memberList {
			if !textRegex.MatchString(id) {
				return SearchAPIRequestParams{}, utilerror.New("request member id regex error", "")
			}
		}
	}

	var tagList []string

	if len(q.tag) > 0 {
		tagList = strings.Split(q.tag, ",")
		for _, tag := range tagList {
			if !textRegex.MatchString(tag) {
				return SearchAPIRequestParams{}, utilerror.New("request tag regex error", "")
			}
		}
	}

	page := defaultPage
	if q.page != "" {
		var err error
		page, err = strconv.Atoi(q.page)
		if err != nil {
			return SearchAPIRequestParams{}, utilerror.New(err.Error(), "").WrapError("request page value parse error")
		}

		if page < 0 {
			return SearchAPIRequestParams{}, utilerror.New("request page count under 0", "")
		}
	}

	var from wrappedbasics.IWrappedTime

	if q.from != "" {
		var err utilerror.IError
		from, err = wrappedbasics.NewWrappedTimeFromLocal(q.from, wrappedbasics.WrappedTimeProps.DateFormat())
		if err != nil {
			return SearchAPIRequestParams{}, err.WrapError("request from value parse error")
		}
	}

	var to wrappedbasics.IWrappedTime
	if q.to != "" {
		var err utilerror.IError
		to, err = wrappedbasics.NewWrappedTimeFromLocal(q.to, wrappedbasics.WrappedTimeProps.DateFormat())
		if err != nil {
			return SearchAPIRequestParams{}, err.WrapError("request to value parse error")
		}
	}

	if from != nil && to != nil {
		if from.After(to) {
			return SearchAPIRequestParams{}, utilerror.New("request date from over to", "")
		}
	}

	maxResult := defaultMaxResult
	if q.maxresult != "" {
		var err error
		maxResult, err = strconv.Atoi(q.maxresult)
		if err != nil {
			return SearchAPIRequestParams{}, utilerror.New(err.Error(), "").WrapError("request maxresult value parse error")
		}

		if maxResult < 0 {
			return SearchAPIRequestParams{}, utilerror.New("request maxResult number under 0", "")
		}
	}

	return NewSearchAPIRequestParams(
		memberList,
		from,
		to,
		tagList,
		page,
		maxResult,
		q.title,
	), nil
}

func NewSearchAPIRequestParams(
	member []string,
	from wrappedbasics.IWrappedTime,
	to wrappedbasics.IWrappedTime,
	tags []string,
	page, maxResult int,
	title string,
) SearchAPIRequestParams {
	return SearchAPIRequestParams{
		member:    member,
		from:      from,
		to:        to,
		tags:      tags,
		page:      page,
		maxResult: maxResult,
		title:     title,
	}
}

func (p SearchAPIRequestParams) ToString() string {
	var stList []string

	if len(p.member) > 0 {
		var list []string
		for _, v := range p.member {
			list = append(list, fmt.Sprintf("\"%s\"", v))
		}
		stList = append(stList, fmt.Sprintf("\"member\": [%s]", strings.Join(list, ", ")))
	}

	if p.from != nil {
		stList = append(stList, fmt.Sprintf("\"from\": \"%s\"", p.from.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateFormat())))
	}

	if p.to != nil {
		stList = append(stList, fmt.Sprintf("\"to\": \"%s\"", p.to.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateFormat())))
	}

	if len(p.tags) > 0 {
		var list []string
		for _, v := range p.tags {
			list = append(list, fmt.Sprintf("\"%s\"", v))
		}
		stList = append(stList, fmt.Sprintf("\"tags\": [%s]", strings.Join(list, ", ")))
	}

	stList = append(stList, fmt.Sprintf("\"page\": [%d]", p.page))
	stList = append(stList, fmt.Sprintf("\"maxResult\": [%d]", p.maxResult))

	if p.title != "" {
		stList = append(stList, fmt.Sprintf("\"title\": \"%s\"", p.title))
	}

	return fmt.Sprintf("{%s}", strings.Join(stList, ","))
}
