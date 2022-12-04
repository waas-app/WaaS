package config

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var logLevel zapcore.Level

func init() {
	viper.AutomaticEnv()
	loggerConfig := zap.NewProductionConfig()
	// make Test env to only log error level
	if IsDevelopment() {
		loggerConfig.Level.SetLevel(zapcore.DebugLevel)
	} else {
		logLevel = zapcore.InfoLevel
		loggerConfig.Level.SetLevel(logLevel)
	}
	// loggerConfig.Level.SetLevel(zapcore.DebugLevel)

	var err error

	// tracer, _ := apm.NewTracerOptions(apm.TracerOptions{
	// 	ServiceName:        config.Spec.ServiceName,
	// 	ServiceEnvironment: config.Spec.Environment,
	// })

	logger, err = loggerConfig.Build(zap.Fields(
		zap.Int("pid", os.Getpid()),
		zap.String("env", viper.GetString("environment")),
	),
	)

	if err != nil {
		panic(err.Error())
	}
	defer logger.Sync()
}

func Logger(ctx context.Context) *otelzap.LoggerWithCtx {
	newLogger := logger
	if ctx == nil {
		ctx = context.Background()
	}

	l := otelzap.New(newLogger, otelzap.WithTraceIDField(true), otelzap.WithMinLevel(logLevel), otelzap.WithStackTrace(true)).Ctx(ctx)
	return &l
}

func IsDevelopment() bool {
	return strings.EqualFold(viper.GetString("environment"), Development)
}
