package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var log *logrus.Logger

func InitLogrus() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.TraceLevel)
	log.SetReportCaller(true)

	path := "./logs/log.log"

	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		//rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationCount(10),
		rotatelogs.WithRotationTime(time.Hour*24))

	log.SetOutput(io.MultiWriter(writer, os.Stdout))
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func InfoWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	log.WithFields(fields).Infof(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func ErrorWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	log.WithFields(fields).Errorf(format, args...)
}
