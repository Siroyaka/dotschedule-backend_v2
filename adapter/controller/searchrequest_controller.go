package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type StreamSearchRequestController struct {
	streamingSearchIntr interactor.StreamingSearchInteractor
	contentType         string
	localTimeDifference int
	defaultMaxResult    int
}

func NewStreamSearchRequestController(
	streamingSearchIntr interactor.StreamingSearchInteractor,
	contentType string,
	localTimeDifference int,
	defaultMaxResult int,
) StreamSearchRequestController {
	return StreamSearchRequestController{
		streamingSearchIntr: streamingSearchIntr,
		contentType:         contentType,
		localTimeDifference: localTimeDifference,
		defaultMaxResult:    defaultMaxResult,
	}
}

func (c StreamSearchRequestController) RequestHandler() http.Handler {
	return http.HandlerFunc(c.searchRequest)
}

type ApiRequestQueries struct {
	member,
	from,
	to,
	tag,
	page,
	title,
	maxresult string
}

type ApiRequestValues struct {
	member    string
	from, to  wrappedbasics.IWrappedTime
	tag       string
	page      int
	title     string
	maxresult int
}

func (c StreamSearchRequestController) readRequestParams(r RequestGet) (string, string, string, string, string, string, string, string) {
	member := r.Get("member")
	from := r.Get("from")
	to := r.Get("to")
	tag := r.Get("tag")
	page := r.Get("page")
	title := r.Get("title")
	maxresult := r.Get("maxresult")
	sortOrderString := r.Get("sort")
	return member, from, to, tag, page, title, maxresult, sortOrderString
}

func (c StreamSearchRequestController) paramsParse(from, to, page, maxresult string) (wrappedbasics.IWrappedTime, wrappedbasics.IWrappedTime, int, int, utilerror.IError) {
	pageCount := 1
	if page != "" {
		var err error
		pageCount, err = strconv.Atoi(page)
		if err != nil {
			return nil, nil, 0, 0, utilerror.New(err.Error(), "")
		}

		if pageCount < 0 {
			return nil, nil, 0, 0, utilerror.New(err.Error(), "")
		}
	}

	var fromDt wrappedbasics.IWrappedTime

	if from != "" {
		var ierr utilerror.IError
		fromDt, ierr = wrappedbasics.NewWrappedTimeFromLocal(from, wrappedbasics.WrappedTimeProps.DateFormat())
		if ierr != nil {
			return nil, nil, 0, 0, ierr.WrapError()
		}
	}

	var toDt wrappedbasics.IWrappedTime
	if to != "" {
		var ierr utilerror.IError
		toDt, ierr = wrappedbasics.NewWrappedTimeFromLocal(to, wrappedbasics.WrappedTimeProps.DateFormat())
		if ierr != nil {
			toDt = nil
		}
	}

	maxResultNum := c.defaultMaxResult
	if maxresult != "" {
		convertValue, err := strconv.Atoi(maxresult)
		if err != nil {
			return nil, nil, 0, 0, utilerror.New(err.Error(), "")
		}

		if convertValue < 0 {
			return nil, nil, 0, 0, utilerror.New(err.Error(), "")
		}

		if convertValue < maxResultNum {
			maxResultNum = convertValue
		}
	}

	return fromDt, toDt, pageCount, maxResultNum, nil
}

func (c StreamSearchRequestController) requestValueLogging(
	member string,
	from wrappedbasics.IWrappedTime,
	to wrappedbasics.IWrappedTime,
	tag string,
	page int,
	titleValue string,
	maxResultNum int,
	sortOrderString string,
) {
	var loggingValues []string

	if member != "" {
		loggingValues = append(loggingValues, fmt.Sprintf("\"member\": \"%s\"", member))
	}

	if from != nil {
		loggingValues = append(loggingValues, fmt.Sprintf("\"from\": \"%s\"", from.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateFormat())))
	}

	if to != nil {
		loggingValues = append(loggingValues, fmt.Sprintf("\"to\": \"%s\"", to.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateFormat())))
	}

	if tag != "" {
		loggingValues = append(loggingValues, fmt.Sprintf("\"tag\": \"%s\"", tag))
	}

	if page != 0 {
		loggingValues = append(loggingValues, fmt.Sprintf("\"page\": \"%d\"", page))
	}

	if titleValue != "" {
		loggingValues = append(loggingValues, fmt.Sprintf("\"title\": \"%s\"", titleValue))
	}

	if maxResultNum != 0 {
		loggingValues = append(loggingValues, fmt.Sprintf("\"maxresult\": \"%d\"", maxResultNum))
	}

	if sortOrderString != "" {
		loggingValues = append(loggingValues, fmt.Sprintf("\"sort\": \"%s\"", sortOrderString))
	}

	if len(loggingValues) == 0 {
		return
	}

	logger.Info(fmt.Sprintf("Streaming Search Request. {%s}", strings.Join(loggingValues, ",")))
}

func (c StreamSearchRequestController) searchRequest(w http.ResponseWriter, r *http.Request) {
	// fromとtoは同じ日付を送られた場合、その日1日分を対象とすることとする
	// fromの翌日日付がtoならfromの日からtoの日の合わせて2日分となるようにすること
	w.Header().Set("Content-Type", c.contentType)

	member, fromValue, toValue, tag, pageValue, titleValue, maxresult, sortOrderString := c.readRequestParams(r.URL.Query())

	from, to, page, maxResultNum, err := c.paramsParse(fromValue, toValue, pageValue, maxresult)

	if err != nil {
		logger.Info(err.WrapError().Error())
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid Request Data.")
		return
	}

	c.requestValueLogging(
		member,
		from,
		to,
		tag,
		page,
		titleValue,
		maxResultNum,
		sortOrderString,
	)

	// toは1日あとの日付にここで加工する
	if to != nil {
		to = to.Add(0, 0, 1, 0, 0, 0)
	}

	value, err := c.streamingSearchIntr.CreateValue(member, from, to, tag, page, titleValue, maxResultNum, sortOrderString)
	if err != nil {
		logger.Info(err.WrapError().Error())
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid Request Data.")
		return
	}

	searchResultCount, err := c.streamingSearchIntr.Count(value)
	if err != nil {
		logger.Fatal(err.WrapError("Streaming Search. Count Query Error"))
		w.WriteHeader(500)
		fmt.Fprint(w, "Service Error")
		return
	}

	if searchResultCount == 0 {
		result := domain.NewAPIResponseData(
			"ok",
			searchResultCount,
			"No Streaming Schedules",
			[]apireference.ScheduleResponse{},
		)
		json, err := result.ToJson()
		if err != nil {
			logger.Fatal(err.WrapError("Streaming Search. Json Convert Error"))
			w.WriteHeader(500)
			fmt.Fprint(w, "Service Error")
			return
		}
		w.WriteHeader(200)
		fmt.Fprint(w, json)
		return
	}

	result, err := c.streamingSearchIntr.Search(value)
	if err != nil {
		logger.Fatal(err.WrapError("Streaming Search. Search Query Error"))
		w.WriteHeader(500)
		fmt.Fprint(w, "Service Error")
		return
	}

	responseData, err := c.streamingSearchIntr.ToJson(result, searchResultCount)
	if err != nil {
		logger.Fatal(err.WrapError("Streaming Search. Json Convert Error"))
		w.WriteHeader(500)
		fmt.Fprint(w, "Service Error")
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, responseData)
}

type RequestGet interface {
	Get(string) string
}
