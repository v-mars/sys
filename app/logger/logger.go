package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/v-mars/frame/pkg/logger"
)

var Logger = logrus.New()
var Gin = logrus.New()

func init()  {
	logger.Logger = Logger
}


func Initial(logLevel,logPath string)  {
	formatter := &logger.Formatter{
		//LogFormat:       "",
		LogFormat:       "%time% [%lvl%] %msg%",
		TimestampFormat: "2006-01-02 15:04:05",
	}
	ginFormatter := &logger.Formatter{
		LogFormat:       "%msg%",
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logger.InitLog(logLevel,"std.log",logPath,formatter, Logger)
	logger.InitLog(logLevel,"gin.log",logPath,ginFormatter, Gin)

}