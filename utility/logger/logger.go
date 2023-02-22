package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Siroyaka/dotschedule-backend_v2/mode"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

const (
	logConfigParent          = "LOG"
	logConfigDir             = "DIR"
	logConfigName            = "NAME"
	logConfigTimeFormat      = "TIME_FORMAT"
	logConfigOSOutput        = "OS_OUTPUT"
	logConfigIsFileSwitch    = "FILE_SWITCH_OF_TIME"
	defaultLogDir            = "./"
	defaultLogName           = "default"
	defaultLogFileTimeFormat = "20060102"
	defaultLogOSOutput       = true
	defaultLogSwitch         = false
)

var (
	suffix      = ""
	logIsSwitch = defaultLogSwitch
)

func loadConfig() config.IConfig {
	if projectConfig := config.ReadProjectConfig(); projectConfig != nil && projectConfig.Has(logConfigParent) {
		return projectConfig.ReadChild(logConfigParent)
	}
	return nil
}

func loadLogFileTimeFormat(loggerConfig config.IConfig) string {
	if loggerConfig == nil {
		return defaultLogFileTimeFormat
	}

	if !loggerConfig.Has(logConfigTimeFormat) {
		return defaultLogFileTimeFormat
	}

	return loggerConfig.Read(logConfigTimeFormat)
}

func isLogOsOutput(loggerConfig config.IConfig) bool {
	if loggerConfig != nil && loggerConfig.Has(logConfigOSOutput) {
		return loggerConfig.ReadBoolean(logConfigOSOutput)
	}
	return defaultLogOSOutput
}

func loadLogFileBasics(loggerConfig config.IConfig) (string, string) {
	if loggerConfig == nil {
		return defaultLogDir, defaultLogName
	}

	logDir := ""
	logName := ""

	if loggerConfig.Has(logConfigDir) {
		logDir = loggerConfig.Read(logConfigDir)
	} else {
		logDir = defaultLogDir
	}

	if loggerConfig.Has(logConfigName) {
		logName = loggerConfig.Read(logConfigName)
	} else {
		logName = defaultLogName
	}

	return logDir, logName
}

func isLogConfigSwitch(loggerConfig config.IConfig) bool {
	if loggerConfig == nil {
		return false
	}

	if !loggerConfig.Has(logConfigIsFileSwitch) {
		return false
	}

	return loggerConfig.ReadBoolean(logConfigIsFileSwitch)
}

func makeLogFilePath(logDir, logName, logSuffix string) string {
	if logDir == "" {
		return ""
	}

	if f, err := os.Stat(logDir); os.IsNotExist(err) || !f.IsDir() {
		Error(utilerror.New(err.Error(), ""))
		return ""
	}

	logFileName := fmt.Sprintf("%s_%s.log", logName, logSuffix)
	return filepath.Join(logDir, logFileName)
}

func makeLogFileSuffix(loggerConfig config.IConfig) string {
	if loggerConfig == nil {
		return ""
	}
	timeFormat := loadLogFileTimeFormat(loggerConfig)
	return time.Now().Format(timeFormat)
}

func loggerSwitch() {
	if !logIsSwitch {
		return
	}

	loggerConfig := loadConfig()
	if loggerConfig == nil {
		return
	}

	newSuffix := makeLogFileSuffix(loggerConfig)
	if suffix == newSuffix {
		return
	}

	suffix = newSuffix

	logDir, logName := loadLogFileBasics(loggerConfig)

	logFilePath := makeLogFilePath(logDir, logName, suffix)

	if logFilePath == "" {
		return
	}

	logIsSwitch = isLogConfigSwitch(loggerConfig)

	isOsOutput := isLogOsOutput(loggerConfig)

	err := loggerSetup(logFilePath, isOsOutput)
	if err != nil {
		Fatal(err.WrapError())
		return
	}
}

func loggerSetup(logFilePath string, isOsOutput bool) utilerror.IError {
	logfile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return utilerror.New(err.Error(), "")
	}
	if isOsOutput {
		log.SetOutput(io.MultiWriter(os.Stdout, logfile))
	} else {
		log.SetOutput(logfile)
	}
	return nil
}

// 実行したディレクトリ直下にyyyyMMdd.logのファイルで出力し、コンソールにも出力する
func defaultLoggerSetup() {
	dateString := time.Now().Format(defaultLogFileTimeFormat)

	logFileName := fmt.Sprintf("%s.log", dateString)

	logFilePath := filepath.Join(defaultLogDir, logFileName)

	loggerSetup(logFilePath, defaultLogOSOutput)
}

func Start() {
	loggerConfig := loadConfig()

	if loggerConfig == nil {
		defaultLoggerSetup()
		return
	}

	logDir, logName := loadLogFileBasics(loggerConfig)

	suffix = makeLogFileSuffix(loggerConfig)

	logFilePath := makeLogFilePath(logDir, logName, suffix)

	if logFilePath == "" {
		return
	}

	logIsSwitch = isLogConfigSwitch(loggerConfig)

	isOsOutput := isLogOsOutput(loggerConfig)

	err := loggerSetup(logFilePath, isOsOutput)
	if err != nil {
		panic(err.Error())
	}
}

func Debug(msg string) {
	if !mode.DEBUG {
		return
	}

	loggerSwitch()
	log.Printf("DEBUG\t%s\n", msg)
}

func Info(msg string) {
	loggerSwitch()
	log.Printf("INFO\t%s\n", msg)
}

func Error(err error) {
	loggerSwitch()
	log.Printf("ERROR\t%s\n", err)
}

func Fatal(err error) {
	loggerSwitch()
	log.Printf("FATAL\t%s\n", err)
}
