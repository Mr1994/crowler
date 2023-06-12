package middleWare

import (
	"api_client/help"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var (
	logFileName = "D:\\go\\goTest\\system"
)

type LogFormat struct {
	StatusCode  int           `gorm:"statusCode" json:"statusCode";`
	LatencyTime time.Duration `gorm:"latencyTime"`
	ClientIP    string        `gorm:"clientIP"`
	ReqMethod   string        `gorm:"reqMethod"`
	ReqUrl      string        `gorm:"reqUrl"`
}

func LogerMiddleware() gin.HandlerFunc {
	now := time.Now()
	timeDate := now.Format("2006-01-02")
	var logFileNameTime string = logFileName + timeDate + ".log"
	// 日志文件
	fileName := path.Join(logFileNameTime)
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	// 实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.InfoLevel)
	//设置输出
	logger.Out = src

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		logFileNameTime,

		// 生成软链，指向最新日志文件
		//rotatelogs.WithLinkName(logFileNameTime),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()

		LogFormats := LogFormat{
			StatusCode:statusCode,
			LatencyTime:latencyTime,
			ClientIP:clientIP,
			ReqMethod:reqMethod,
			ReqUrl:reqUrl,
		}
		fmt.Println(LogFormats)
		new(help.Common).Log(LogFormats, "test")

	}
}

//func Logger() gin.HandlerFunc {
//	file, err := os.OpenFile(config.ROOT_PATH+"/logs/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
//	if err != nil {
//		panic(err)
//	}
//	return gin.LoggerWithConfig(gin.LoggerConfig{
//		Output: file,
//	})
//}
