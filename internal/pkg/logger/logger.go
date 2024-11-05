package logger

import (
	"github.com/sirupsen/logrus"
)

func NewWithDefaultLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.InfoLevel)
	return l
}
