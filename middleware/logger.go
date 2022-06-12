package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	filePath := "log/log"
	linkName := "latest_log.log" //软链接
	src,err := os.OpenFile(filePath,os.O_RDWR|os.O_CREATE, 0755)  //将日志写入文件
	if err != nil {
		fmt.Println("err:",err)
	}
	logger := logrus.New()
	logger.Out = src  // 实例化之后写入文件
	//切割日志
	logger.SetLevel(logrus.DebugLevel)  //设置日志级别
	logWriter,_ := retalog.New(
		filePath+"%Y%m%d.log",
		retalog.WithMaxAge(7*24*time.Hour), //保存时间
		retalog.WithRotationTime(24*time.Hour),//多长时间分隔一次
		retalog.WithLinkName(linkName), //软链接(也就是可以在项目根目录看到最新的日志文件)
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	hook := lfshook.NewHook(writeMap,&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(hook)
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms",int(math.Ceil(float64(stopTime.Nanoseconds())/100000.0)))
		hostName,err := os.Hostname()
		if err != nil {
			hostName = "unknow"
		}
		statusCode := c.Writer.Status()  // 请求码
		clientIp := c.ClientIP()  //客户端IP
		userAgent := c.Request.UserAgent()  //客户端信息
		dataSize := c.Writer.Size()  //请求数据大小
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":hostName,
			"status":statusCode,
			"SpendTime":spendTime,
			"IP":clientIp,
			"Method":method,
			"Path":path,
			"DataSize":dataSize,
			"Agent":userAgent,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())  //记录系统内部错误
		}
		if statusCode >= 500 {
			entry.Error()
		}else if statusCode >= 400 {
			entry.Warn()
		}else {
			entry.Info()
		}
	}
}
