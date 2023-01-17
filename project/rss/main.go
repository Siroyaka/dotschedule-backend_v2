package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	rssmaster "github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/rss/master"
	rssrequest "github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/rss/request"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
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

	if err := utility.LoggerSetup(); err != nil {
		log.Fatal(err)
	}

	publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)

	// import
	common := utility.NewCommon(publicConfig)

	// infrastructure
	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	httpRequestHandler := infrastructure.NewHTTPRequest()

	// repository
	getMasterRepos := rssmaster.NewGetRepository(
		sqlHandler,
		queryConfig.Read(config_getRSSMaster),
	)

	updateMasterRepos := rssmaster.NewUpdateRepository(
		sqlHandler,
		queryConfig.Read(config_updateRSSMaster),
	)

	requestRepos := rssrequest.NewRequestRepository(httpRequestHandler)

	getScheduleRepos := rssschedule.NewGetRepository(
		sqlHandler,
		common,
		queryConfig.Read(config_getSchedule),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	insertScheduleRepos := rssschedule.NewInsertRepository(
		sqlHandler,
		common,
		queryConfig.Read(config_insertSchedule),
	)

	updateScheduleRepos := rssschedule.NewUpdateRepository(
		sqlHandler,
		common,
		queryConfig.Read(config_updateSchedule),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	// usecase
	rssInteractor := usecase.NewRSSInteractor(
		common,
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
