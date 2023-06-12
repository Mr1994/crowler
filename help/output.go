package help

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"net/http"
	"os"
	"path"
	"time"
)

// 错误处理的结构体
type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	Success     = NewError(http.StatusOK, 0, "success")
	ServerError = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	NotFound    = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

type Body struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errMsg"`
}

func OutError(request *gin.Context, errMsg string) {
	var data Body

	if length := len(request.Errors); length > 0 {
		e := request.Errors[length-1]
		err := e.Err
		if err != nil {
			var Err *Error
			if e, ok := err.(*Error); ok {
				Err = e
			} else if e, ok := err.(error); ok {
				Err = OtherError(e.Error())
			} else {
				Err = ServerError
			}

			// 输出国定格式
			data.Code = http.StatusOK
			data.ErrMsg = Err.Msg
			data.Data = make(map[string]string)
			//非阻塞报警
			request.JSON(http.StatusOK, data)

		}
	} else {
		data.Code = http.StatusSwitchingProtocols
		data.ErrMsg = errMsg
		data.Data = make(map[string]string)
		// 记录一个错误的日志
		request.JSON(http.StatusOK, data)
		return

	}
}

func Output(msg interface{}, requeset *gin.Context) {
	var data Body
	data.Code = http.StatusOK
	data.ErrMsg = "请求成功"
	data.Data = msg
	requeset.JSON(http.StatusOK, data)
	return
}

func (c *Common) Log(Msg interface{}, logFileName string) {
	// 获取ini配置文件
	Conf, _ := ini.Load(Appini) //此路径为ini文件的路径
	logName := Conf.Section("LOG").Key("LogPath").String()
	msgString := c.InterFaceToString(Msg)
	now := time.Now()
	timeDate := now.Format("2006-01-02")
	var logFileNameTime string = logName + logFileName + timeDate + ".log"
	fmt.Println(logFileNameTime)
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
		FieldMap: logrus.FieldMap{
			"info": msgString,
		},
	}))

	// 日志格式
	logger.WithFields(logrus.Fields{
		"info": msgString,
	}).Info()
}
