package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

var customLog = logrus.New()
var Fields = make(map[string]interface{})

func init() {
	customLog.Level = logrus.TraceLevel
	customLog.Out = os.Stdout
	customLog.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func LogInfo(message string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		customLog.Infof(message)
	} else {
		customLog.WithFields(fields[0]).Info(message)
	}
}

func LogWarning(message string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		customLog.Warn(message)
	} else {
		customLog.WithFields(fields[0]).Warn(message)
	}
}

func LogError(message string, err error) {
	customLog.WithError(err).Error(message)
}
