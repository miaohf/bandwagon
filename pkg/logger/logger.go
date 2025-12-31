package logger

import (
	"io"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

// Init 初始化日志系统
func Init() {
	// 创建不同级别的日志记录器
	infoLogger = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
}

// Logger 日志接口
type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// logger 实现
type logger struct{}

// GetLogger 获取日志实例
func GetLogger() Logger {
	return &logger{}
}

func (l *logger) Info(v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Println(v...)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}

func (l *logger) Error(v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Println(v...)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
}

func (l *logger) Debug(v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Println(v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}

func (l *logger) Fatal(v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Fatal(v...)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Fatalf(format, v...)
	}
}

// SetOutput 设置日志输出
func SetOutput(w io.Writer) {
	if infoLogger != nil {
		infoLogger.SetOutput(w)
	}
	if errorLogger != nil {
		errorLogger.SetOutput(w)
	}
	if debugLogger != nil {
		debugLogger.SetOutput(w)
	}
}
