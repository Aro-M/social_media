package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(level string) *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	var l logrus.Level
	switch level {
	case "debug":
		l = logrus.DebugLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	default:
		l = logrus.InfoLevel
	}
	log.SetLevel(l)

	return log
}
