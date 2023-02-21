package controller

import (
	"fmt"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type ViewScheduleController struct {
	scheduleIntr interactor.DayScheduleInteractor
	contentType  string
}

func NewScheduleController(scheduleIntr interactor.DayScheduleInteractor, contentType string) ViewScheduleController {
	return ViewScheduleController{
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

	fromDate, err := wrappedbasics.NewWrappedTimeFromLocal(date, wrappedbasics.WrappedTimeProps.DateFormat())

	if err != nil {
		utility.LogError(err.WrapError())
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid Date Format. RequestString: %s", date)
		utility.LogInfo(fmt.Sprintf("Invalid Date Format. RequestString: %s", date))
		return
	}

	toDate := fromDate.Add(0, 0, 1, 0, 0, 0)

	utility.LogInfo(fmt.Sprintf("Schedule Request. Date: %s", fromDate.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateFormat())))

	list, err := c.scheduleIntr.GetScheduleData(fromDate, toDate)

	if err != nil {
		utility.LogFatal(err.WrapError())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", date)
		//c.loggerInteractor.Fatal(fmt.Sprintf("Data fetch error: %s", date), err, 1)
		return
	}

	json, err := c.scheduleIntr.ToJson(list)

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
