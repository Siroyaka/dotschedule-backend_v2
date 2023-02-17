package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/httprequest"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlrss"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

const (
	configPath  = "./config.json"
	projectName = "RSS"

	config_public = "PUBLIC"
	config_sql    = "SQLITE"
	config_query  = "QUERY"

	config_sqlPath                 = "PATH"
	config_sqlReplaceTargetString  = "TARGETS_STRING"
	config_sqlReplacedChar         = "REPLACED_CHAR"
	config_sqlReplacedCharSplitter = "REPLACED_CHAR_SPLITTER"

	// root config
	config_scheduleNotCompleteStatus = "NOT_COMPLETE_STATUS"
	config_scheduleCompleteStatus    = "COMPLETE_STATUS"
	config_scheduleUnVisibleStatus   = "UNVISIBLE_STATUS"
	config_platformType              = "PLATFORM_TYPE"
	config_insertStatus              = "INSERT_STATUS"
	config_feedVideoIdElement        = "FEED_VIDEOID_ELEMENT"
	config_feedVideoIdItemName       = "FEED_VIDEOID_ITEMNAME"

	// query
	config_getRSSMaster    = "GET_RSSMASTER"
	config_updateRSSMaster = "UPDATE_RSSMASTER"
	config_getSchedule     = "GET_SCHEDULE"
	config_updateSchedule  = "UPDATE_SCHEDULE"
	config_insertSchedule  = "INSERT_SCHEDULE"
)

func main() {
	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	utility.LoggerStart()
	wrappedbasics.InitializeWrappedTimeProps()

	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)

	// infrastructure
	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	httpRequestHandler := infrastructure.NewHTTPRequest()

	// repository
	getMasterRepos := sqlrss.NewSelectRSSMasterRepository(
		sqlHandler,
		queryConfig.Read(config_getRSSMaster),
	)

	updateMasterRepos := sqlrss.NewUpdateRSSMasterRepository(
		sqlHandler,
		queryConfig.Read(config_updateRSSMaster),
	)

	requestRepos := httprequest.NewRSSRequestRepository(httpRequestHandler)

	getScheduleRepos := sqlrepository.NewSelectIDsFromIDsRepository(
		sqlHandler,
		queryConfig.Read(config_getSchedule),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	insertScheduleRepos := sqlrepository.NewInsertRSSSeedScheduleRepository(
		sqlHandler,
		queryConfig.Read(config_insertSchedule),
	)

	updateScheduleRepos := rssschedule.NewUpdateRepository(
		sqlHandler,
		queryConfig.Read(config_updateSchedule),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	// usecase
	rssInteractor := interactor.NewRSSInteractor(
		getMasterRepos,
		updateMasterRepos,
		requestRepos,
		getScheduleRepos,
		insertScheduleRepos,
		updateScheduleRepos,
		rootConfig.ReadInteger(config_scheduleCompleteStatus),
		rootConfig.ReadInteger(config_scheduleNotCompleteStatus),
		rootConfig.Read(config_platformType),
		rootConfig.Read(config_insertStatus),
		rootConfig.Read(config_feedVideoIdElement),
		rootConfig.Read(config_feedVideoIdItemName),
	)

	// controller
	controller := controller.NewRSSController(rssInteractor)

	controller.Exec()
}
