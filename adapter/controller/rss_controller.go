package controller

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
)

type RSSController struct {
	intr interactor.RSSInteractor
}

func NewRSSController(intr interactor.RSSInteractor) RSSController {
	return RSSController{
		intr: intr,
	}
}

func (rc RSSController) Exec() {
	err := rc.intr.GetMaster()
	if err != nil {
		logger.Fatal(err.WrapError())
		return
	}
	totalInsert := 0
	totalUpdate := 0
	for rc.intr.Next() {
		list, err := rc.intr.GetRSSData()
		if err != nil {
			logger.Fatal(err.WrapError())
			continue
		}
		insert, update, isError, newestUpdate, err := rc.intr.PushToDB(list)
		if err != nil {
			logger.Fatal(err.WrapError())
			continue
		}

		if isError {
			continue
		}

		if (insert + update) == 0 {
			continue
		}

		if err = rc.intr.EndRow(newestUpdate); err != nil {
			logger.Fatal(err.WrapError())
			continue
		}
		totalInsert += insert
		totalUpdate += update
	}
	logger.Info(fmt.Sprintf("Get RSSFeed data end. insert_count: %d, update_count: %d", totalInsert, totalUpdate))
}
