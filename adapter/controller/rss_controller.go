package controller

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type RSSController struct {
	intr usecase.RSSInteractor
}

func NewRSSController(intr usecase.RSSInteractor) RSSController {
	return RSSController{
		intr: intr,
	}
}

func (rc RSSController) Exec() {
	err := rc.intr.GetMaster()
	if err != nil {
		utility.LogFatal(err.WrapError())
		return
	}
	totalInsert := 0
	totalUpdate := 0
	for rc.intr.Next() {
		list, err := rc.intr.GetRSSData()
		if err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}
		insert, update, isError, newestUpdate, err := rc.intr.PushToDB(list)
		if err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}

		if isError {
			continue
		}

		if (insert + update) == 0 {
			continue
		}

		if err = rc.intr.EndRow(newestUpdate); err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}
		totalInsert += insert
		totalUpdate += update
	}
	utility.LogInfo(fmt.Sprintf("Get RSSFeed data end. insert_count: %d, update_count: %d", totalInsert, totalUpdate))
}
