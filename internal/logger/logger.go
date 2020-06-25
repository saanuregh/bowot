package logger

import (
	"bowot/internal/config"
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/andersfylling/disgord"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var CustomLogger *LoggerZap

func init() {
	conf := zap.NewProductionConfig()
	if config.C.Logging.Level == "debug" {
		conf.Level.SetLevel(zap.DebugLevel)
	} else {
		conf.Level.SetLevel(zap.InfoLevel)
	}
	logger, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	cfg := zapsentry.Configuration{
		Level: zapcore.ErrorLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(config.C.Logging.SentryDsn))
	if err != nil {
		panic(err)
	}
	logger = zapsentry.AttachCoreToLogger(core, logger)
	CustomLogger = &LoggerZap{
		instance: logger.With(),
	}
}

type LoggerZap struct {
	instance *zap.Logger
}

var _ disgord.Logger = (*LoggerZap)(nil)

func (log *LoggerZap) Debug(v ...interface{}) {
	log.instance.Debug(fmt.Sprint(v...))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Info(v ...interface{}) {
	log.instance.Info(fmt.Sprint(v...))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Error(v ...interface{}) {
	log.instance.Error(fmt.Sprint(v...))
	_ = log.instance.Sync()
}

func (log *LoggerZap) Panic(v ...interface{}) {
	log.instance.Panic(fmt.Sprint(v...))
	_ = log.instance.Sync()
}
