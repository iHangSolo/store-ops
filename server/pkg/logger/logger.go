package logger

import (
	"log"
	"os"
)

var (
	infoLogger  = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
)

func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

func Error(v ...interface{}) {
	errorLogger.Println(v...)
}

func Debug(v ...interface{}) {
	debugLogger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	debugLogger.Printf(format, v...)
}