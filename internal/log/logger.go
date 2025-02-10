package log

import (
	"sync"

	"github.com/K4rian/dslogger"
	"github.com/K4rian/kfrs/internal/config"
)

var (
	Logger *dslogger.Logger
	once   sync.Once
)

func Init() {
	once.Do(initLogger)
}

func initLogger() {
	conf := config.Get()

	logLevel := conf.LogLevel
	loggerConfig := &dslogger.Config{
		LogFile:       conf.LogFile,
		LogFileFormat: dslogger.LogFormat(conf.LogFileFormat),
		MaxSize:       conf.LogMaxSize,
		MaxBackups:    conf.LogMaxBackups,
		MaxAge:        conf.LogMaxAge,
		Level:         logLevel,
	}

	if conf.LogToFile {
		Logger = dslogger.NewLogger(logLevel, loggerConfig)
	} else {
		Logger = dslogger.NewConsoleLogger(logLevel, loggerConfig)
	}
}
