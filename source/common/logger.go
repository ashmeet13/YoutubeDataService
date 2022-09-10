package common

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func GetLogger() *logrus.Logger {
	if Logger == nil {
		Logger = logrus.StandardLogger()
	}
	return Logger
}
