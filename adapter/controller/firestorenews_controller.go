package controller

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type FirestoreNewsController struct {
	firestoreNewsIntr usecase.FirestoreNewsInteractor
}

func NewFirestoreNewsController(firestoreNewsIntr usecase.FirestoreNewsInteractor) FirestoreNewsController {
	return FirestoreNewsController{
		firestoreNewsIntr: firestoreNewsIntr,
	}
}

func (controller FirestoreNewsController) Exec() {
	if err := controller.firestoreNewsIntr.DataFetchFromFirestore(); err != nil {
		utility.LogFatal(err.WrapError())
		return
	}

	controller.firestoreNewsIntr.UpdateDB()
}
