package controller

import (
	"fmt"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type ViewScheduleController struct {
	common       utility.Common
	scheduleIntr usecase.ScheduleInteractor
	contentType  string
}

func NewScheduleController(common utility.Common, scheduleIntr usecase.ScheduleInteractor, contentType string) ViewScheduleController {
	return ViewScheduleController{
		common:       common,
		scheduleIntr: scheduleIntr,
		contentType:  contentType,
	}
}

func (c ViewScheduleController) ScheduleRequestHandler() http.Handler {
	return http.HandlerFunc(c.scheduleRequest)
}

func (c ViewScheduleController) scheduleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", c.contentType)
	date := r.URL.Query().Get("date")
	if date == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid Request")
		return
	}

	baseDate, err := c.common.CreateNewWrappedTimeFromLocalDate(date)

	if err != nil {
		utility.LogError(err.WrapError())
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid Date Format: %s", date)
		return
	}

	list, err := c.scheduleIntr.GetScheduleData(baseDate)

	if err != nil {
		utility.LogFatal(err.WrapError())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", date)
		//c.loggerInteractor.Fatal(fmt.Sprintf("Data fetch error: %s", date), err, 1)
		return
	}

	json, err := usecase.ScheduleDataToResponseJson(list)

	if err != nil {
		utility.LogFatal(err.WrapError())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", date)
		//c.loggerInteractor.Fatal(fmt.Sprintf("Json convert error: %s", date), err, 1)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, json)
}
