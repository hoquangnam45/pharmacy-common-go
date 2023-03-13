package util

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = logger
	SugaredLogger = Logger.Sugar()
}
