package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlapi"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

const (
	configPath  = "./config.json"
	projectName = "API"

	heartBeatRoute = "/"

	config_public = "PUBLIC"
	config_sql    = "SQLITE"

	config_sqlPath          = "PATH"
	config_sqlDataSplitter  = "DATA_SPLITTER"
	config_sqlArraySplitter = "ARRAY_SPLITTER"

	config_query       = "QUERY"
	config_getschedule = "GET_DAYSCHEDULE"
	config_getmonth    = "GET_DAYS_PARTICIPANTS_LIST"

	config_localTimeDifference = "LOCAL_TIMEDIFFERENCE"
	config_viewingStatus       = "VIEWING_STATUS"
	config_contentType         = "CONTENT_TYPE"
	config_port                = "PORT"
	config_scheduleRoute       = "ROUTE_SCHEDULE"
	config_monthRoute          = "ROUTE_MONTH"
)

func heartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintln(w, "ok")
}

func main() {
	//confLoader := fileio.NewReaderRepository[config.ConfigData](infrastructure.NewFileReader[config.ConfigData]())

	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	//configValue, err := confLoader.ReadJson(configPath)

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	utility.LoggerStart()
	wrappedbasics.InitializeWrappedTimeProps()

	//publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)

	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	scheduleRepository := sqlapi.NewSelectSchedulesRepository(
		sqlHandler,
		queryConfig.Read(config_getschedule),
		rootConfig.ReadInteger(config_viewingStatus),
	)

	selectDaysParticipantsRepository := sqlapi.NewSelectDaysParticipantsRepository(
		sqlHandler,
		queryConfig.Read(config_getmonth),
		rootConfig.ReadInteger(config_viewingStatus),
		rootConfig.ReadInteger(config_localTimeDifference),
	)

	scheduleInteractor := interactor.NewDayScheduleInteractor(
		scheduleRepository,
	)

	monthDataInteractor := interactor.NewDaysParticipantsInteractor(
		selectDaysParticipantsRepository,
	)

	scController := controller.NewScheduleController(
		scheduleInteractor,
		rootConfig.Read(config_contentType),
	)

	monthController := controller.NewMonthRequestController(
		monthDataInteractor,
		rootConfig.Read(config_contentType),
		rootConfig.ReadInteger(config_localTimeDifference),
	)

	router := infrastructure.NewRouter(rootConfig.Read(config_port))
	router.SetHandle(rootConfig.Read(config_scheduleRoute), scController.ScheduleRequestHandler())
	router.SetHandle(rootConfig.Read(config_monthRoute), monthController.MonthRequestHandler())
	router.SetHandleFunc(heartBeatRoute, heartBeat)
	router.Run()
}
