package controller

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
)

type NormalizationController struct {
	intr interactor.NormalizationYoutubeDataInteractor
}

func NewNormalizationController(
	intr interactor.NormalizationYoutubeDataInteractor,
) NormalizationController {
	return NormalizationController{
		intr: intr,
	}
}

func (controller NormalizationController) Execute() {
	logger.Info("Start Normalization")
	if err := controller.intr.Normalization(); err != nil {
		logger.Error(err.WrapError())
	}
	logger.Info("End Normalization")
}
