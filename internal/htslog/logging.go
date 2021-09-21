package htslog

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	log *logrus.Logger
)

// Setup sets up the configuration of a global logger variable
// Should only be called once in the main
func Setup(logFile string, logLevel string) {

	// set up a logrus instance for the application
	log = logrus.New()

	// if a log file is specified then lets hook our logger up to that, otherwise leave as default logrus output
	if htsconfig.GetLogFile() != "" {
		file, err := os.OpenFile(htsconfig.GetLogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err == nil {
			log.SetOutput(io.MultiWriter(os.Stderr, file))
		} else {
			log.Warn("Failed to set up writing to configuration log file setting, leaving as default stderr")
		}
	}

	// allow the config to also set the desired log level
	lev, err := logrus.ParseLevel(htsconfig.GetLogLevel())

	if err == nil {
		log.SetLevel(lev)
	} else {
		log.Warn("Did not understand configuration log level setting, leaving as default level")
	}

	log.Debug("Testing debug level")
	log.Info("Testing info level")
	log.Warn("Testing warn level")
	log.Error("Testing error level")

	//log.Formatter = &logrus.JSONFormatter{}
	//log.SetReportCaller(true)
}

func Debug(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}
