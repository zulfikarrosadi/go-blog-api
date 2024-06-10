package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logrus *logrus.Logger = logrus.New()

func init() {
	Logrus.SetFormatter(&logrus.JSONFormatter{})
	f, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	// defer f.Close()
	fmt.Println(f, err)
	Logrus.SetOutput(f)
}

func ErrorLog(action string, message string, err error) {
	recover()
	Logrus.WithFields(logrus.Fields{
		"timestamp": time.Now(),
		"details":   err.Error(),
		"context": map[string]any{
			"action": action,
		},
	}).Error(message)
}
