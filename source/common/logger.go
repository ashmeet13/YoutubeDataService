package common

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func GetLogger() *logrus.Entry {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return logger
}
