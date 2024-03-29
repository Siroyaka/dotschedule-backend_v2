package main

import (
	"fmt"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/fileio"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlapi"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
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

	config_sqlReplaceTargetString  = "TARGETS_STRING"
	config_sqlReplacedChar         = "REPLACED_CHAR"
	config_sqlReplacedCharSplitter = "REPLACED_CHAR_SPLITTER"

	config_query                   = "QUERY"
	config_getschedule             = "GET_DAYSCHEDULE"
	config_getmonth                = "GET_DAYS_PARTICIPANTS_LIST"
	config_getstreamingsearch      = "GET_STREAMING_SEARCH"
	config_countstreamingsearchlen = "COUNT_STREAMING_SEARCH_LEN"

	config_subquery                    = "SUB"
	config_substreamingsearchmember    = "STREAMING_SEARCH_MEMBER"
	config_substreamingsearchtag       = "STREAMING_SEARCH_TAG"
	config_substreamingsearchfrom      = "STREAMING_SEARCH_FROM"
	config_substreamingsearchto        = "STREAMING_SEARCH_TO"
	config_substreamingsearchtitle     = "STREAMING_SEARCH_TITLE"
	config_substreamingsearchisviewing = "STREAMING_SEARCH_ISVIEWING"
	config_substreamingsearchsortnewer = "STREAMING_SEARCH_SORT_NEWER"
	config_substreamingsearchsortolder = "STREAMING_SEARCH_SORT_OLDER"
	config_substreamingsearchpagelimit = "STREAMING_SEARCH_PAGE_LIMIT"

	config_searchconstructions = "SEARCH_CONSTRUCTIONS"
	config_searchdefaultfrom   = "DEFAULT_FROM"
	config_searchlenLimit      = "LIMIT"
	config_searchdefaultsort   = "DEFAULT_SORT"

	config_localTimeDifference = "LOCAL_TIMEDIFFERENCE"
	config_viewingStatus       = "VIEWING_STATUS"
	config_contentType         = "CONTENT_TYPE"
	config_port                = "PORT"
	config_scheduleRoute       = "ROUTE_SCHEDULE"
	config_monthRoute          = "ROUTE_MONTH"
	config_streamSearchRoute   = "ROUTE_STREAMSEARCH"
)

func heartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintln(w, "ok")
}

func loadConfig() config.IConfig {
	reader := fileio.NewReaderRepository[map[string]interface{}](infrastructure.NewFileReader[map[string]interface{}]())
	data, err := reader.ReadJson(configPath)
	if err != nil {
		panic(err)
	}
	return config.New(data)
}

func main() {
	config.Setup(projectName, loadConfig())

	logger.Start()
	wrappedbasics.InitializeWrappedTimeProps()

	//publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)
	subqueryConfig := queryConfig.ReadChild(config_subquery)
	searchConfig := rootConfig.ReadChild(config_searchconstructions)

	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	selectDaysParticipantsRepository := sqlapi.NewSelectDaysParticipantsRepository(
		sqlHandler,
		queryConfig.Read(config_getmonth),
		rootConfig.ReadInteger(config_viewingStatus),
		rootConfig.ReadInteger(config_localTimeDifference),
	)

	searchRepos := sqlapi.NewSelectStreamingSearchRepository(
		sqlHandler,
		queryConfig.Read(config_getstreamingsearch),
		subqueryConfig.Read(config_substreamingsearchmember),
		subqueryConfig.Read(config_substreamingsearchtag),
		subqueryConfig.Read(config_substreamingsearchfrom),
		subqueryConfig.Read(config_substreamingsearchto),
		subqueryConfig.Read(config_substreamingsearchtitle),
		subqueryConfig.Read(config_substreamingsearchisviewing),
		searchConfig.Read(config_searchdefaultfrom),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
		rootConfig.ReadInteger(config_viewingStatus),
		subqueryConfig.Read(config_substreamingsearchsortnewer),
		subqueryConfig.Read(config_substreamingsearchsortolder),
		subqueryConfig.Read(config_substreamingsearchpagelimit),
	)

	countRepos := sqlapi.NewCountStreamingSearchRepository(
		sqlHandler,
		queryConfig.Read(config_countstreamingsearchlen),
		subqueryConfig.Read(config_substreamingsearchmember),
		subqueryConfig.Read(config_substreamingsearchtag),
		subqueryConfig.Read(config_substreamingsearchfrom),
		subqueryConfig.Read(config_substreamingsearchto),
		subqueryConfig.Read(config_substreamingsearchtitle),
		subqueryConfig.Read(config_substreamingsearchisviewing),
		searchConfig.Read(config_searchdefaultfrom),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
		rootConfig.ReadInteger(config_viewingStatus),
	)

	monthDataInteractor := interactor.NewDaysParticipantsInteractor(
		selectDaysParticipantsRepository,
	)

	streamingSearchInteractor := interactor.NewStreamingSearchInteractor(
		searchRepos,
		countRepos,
	)

	monthController := controller.NewMonthRequestController(
		monthDataInteractor,
		rootConfig.Read(config_contentType),
		rootConfig.ReadInteger(config_localTimeDifference),
	)

	streamSearchController := controller.NewStreamSearchRequestController(
		streamingSearchInteractor,
		rootConfig.Read(config_contentType),
		rootConfig.ReadInteger(config_localTimeDifference),
		searchConfig.ReadInteger(config_searchlenLimit),
	)

	router := infrastructure.NewRouter(rootConfig.Read(config_port))
	router.SetHandle(rootConfig.Read(config_monthRoute), monthController.RequestHandler())
	router.SetHandle(rootConfig.Read(config_streamSearchRoute), streamSearchController.RequestHandler())
	router.SetHandleFunc(heartBeatRoute, heartBeat)
	router.Run()
}
