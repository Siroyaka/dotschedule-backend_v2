package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/viewschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
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
	config_getschedule = "GET_SCHEDULE"
	config_getmonth    = "GET_MONTH_DATA"

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
	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	if err := utility.LoggerStart(); err != nil {
		log.Fatal(err)
	}

	publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)

	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	common := utility.NewCommon(publicConfig)

	scheduleRepository := viewschedule.NewGetRepository(
		sqlHandler,
		queryConfig.Read(config_getschedule),
		queryConfig.Read(config_getmonth),
		rootConfig.Read(config_localTimeDifference),
	)

	scheduleInteractor := usecase.NewScheduleInteractor(
		scheduleRepository,
		sqlConfig.Read(config_sqlDataSplitter),
		sqlConfig.Read(config_sqlArraySplitter),
		rootConfig.ReadInteger(config_viewingStatus),
		common,
	)
	monthDataInteractor := usecase.NewMonthInteractor(
		scheduleRepository,
		sqlConfig.Read(config_sqlDataSplitter),
		sqlConfig.Read(config_sqlArraySplitter),
		rootConfig.ReadInteger(config_viewingStatus),
		common,
	)

	scController := controller.NewScheduleController(
		common,
		scheduleInteractor,
		rootConfig.Read(config_contentType),
	)
	monController := controller.NewMonthRequestController(
		common,
		monthDataInteractor,
		rootConfig.Read(config_contentType),
	)

	router := infrastructure.NewRouter(rootConfig.Read(config_port))
	router.SetHandle(rootConfig.Read(config_scheduleRoute), scController.ScheduleRequestHandler())
	router.SetHandle(rootConfig.Read(config_monthRoute), monController.MonthRequestHandler())
	router.SetHandleFunc(heartBeatRoute, heartBeat)
	router.Run()
}
