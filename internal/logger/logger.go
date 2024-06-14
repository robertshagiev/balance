package logger

import (
	"log"
	"os"
)

type Logger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
}

func New() *Logger {
	return &Logger{
		debugLogger:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger:   log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Debug(msg string) {
	l.debugLogger.Println(msg)
}

func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *Logger) Warning(msg string) {
	l.warningLogger.Println(msg)
}

func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

func (l *Logger) Fatal(msg string) {
	l.fatalLogger.Fatal(msg)
}
