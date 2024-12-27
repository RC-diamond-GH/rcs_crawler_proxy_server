package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// 初始化全局日志
func InitLogger(logFile *os.File) {
	logger = logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.InfoLevel)
}

// 获取全局日志实例
func GetLogger() *logrus.Logger {
	return logger
}
