package controller

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
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
	utility.LogInfo("Start Normalization")
	if err := controller.intr.Normalization(); err != nil {
		utility.LogError(err.WrapError())
	}
	utility.LogInfo("End Normalization")
}
