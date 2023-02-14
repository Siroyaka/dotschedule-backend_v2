package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/discordpost"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	infraYoutube "github.com/Siroyaka/dotschedule-backend_v2/infrastructure/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/dbmodels"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
)

const (
	configPath  = "./config.json"
	projectName = "YOUTUBE_SCHEDULE_DATA_UPDATE"

	config_public  = "PUBLIC"
	config_sql     = "SQLITE"
	config_parser  = "PARSER"
	config_query   = "QUERY"
	config_discord = "DISCORD"

	config_sqlPath = "PATH"

	// root
	config_platform           = "PLATFORM"
	config_status100Deferment = "STATUS100_DEFERMENT"
	config_developerKey       = "DEVELOPER_KEY"
	config_streamingUrlPrefic = "STREAMING_URL_PREFIX"
	config_partList           = "PART_LIST"

	// parser
	config_parserHour   = "HOUR"
	config_parserMinite = "MINITE"
	config_parserSecond = "SECOND"

	// query
	config_getSchedule         = "GET_SCHEDULE"
	config_getStreamerMaster   = "GET_STREAMER_MASTER"
	config_updateSchedule      = "UPDATE_SCHEDULE"
	config_updateScheduleTo100 = "UPDATE_SCHEDULE_TO_STATUS100"
	config_getParticipants     = "GET_PARTICIPANTS"
	config_insertParticipants  = "INSERT_PARTICIPANTS"

	// discord
	config_discordUrl                = "NORTIFICATION_URL"
	config_discordNortificationRange = "TARGET_TIME_RANGE"
)

func main() {
	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	utility.LoggerStart()

	publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)
	parserConfig := rootConfig.ReadChild(config_parser)
	discordConfig := rootConfig.ReadChild(config_discord)

	// import
	common := utility.NewCommon(publicConfig)

	// infrastructure
	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	youtubeDataAPI := infraYoutube.NewYoutubeDataAPI(rootConfig.Read(config_developerKey))

	httpRequest := infrastructure.NewHTTPRequest()

	// repository
	youtubeVideoListRepos := youtubedataapi.NewVideoListRepository(youtubeDataAPI)

	getScheduleRepos := sqlwrapper.NewSelectRepository[domain.FullScheduleData](sqlHandler, queryConfig.Read(config_getSchedule))

	getStreamerMasterRepos := sqlwrapper.NewSelectRepository[domain.StreamerMasterWithPlatformData](sqlHandler, queryConfig.Read(config_getStreamerMaster))

	updateScheduleRepos := sqlwrapper.NewUpdateRepository(sqlHandler, queryConfig.Read(config_updateSchedule))

	updateScheduleTo100Repos := sqlwrapper.NewUpdateRepository(sqlHandler, queryConfig.Read(config_updateScheduleTo100))

	getParticipantsRepos := sqlwrapper.NewSelectRepository[dbmodels.KeyValue[string, string]](sqlHandler, queryConfig.Read(config_getParticipants))

	insertParticipantsRepos := sqlwrapper.NewUpdateRepository(sqlHandler, queryConfig.Read(config_insertParticipants))

	discordPostRepos := discordpost.NewDiscordPostRepository(httpRequest, discordConfig.Read(config_discordUrl))

	// interactor
	intr := interactor.NewNormalizationYoutubeDataInteractor(
		getStreamerMasterRepos,
		getScheduleRepos,
		updateScheduleRepos,
		updateScheduleTo100Repos,
		getParticipantsRepos,
		insertParticipantsRepos,
		youtubeVideoListRepos,
		discordPostRepos,
		common,
		utility.NewYoutubeDurationParser(
			parserConfig.Read(config_parserHour),
			parserConfig.Read(config_parserMinite),
			parserConfig.Read(config_parserSecond),
		),
		rootConfig.Read(config_platform),
		rootConfig.Read(config_streamingUrlPrefic),
		rootConfig.ReadStringList(config_partList),
		discordConfig.ReadInteger(config_discordNortificationRange),
	)

	// controller
	controller := controller.NewNormalizationController(intr)

	controller.Execute()
}
