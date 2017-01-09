package trace

import (
	"encoding/json"
	"security/conf"
	"time"

	"github.com/astaxie/beego/logs"
)

var logger = logs.NewLogger()
var trackId string

type logFormat struct {
	//访问的日志id
	TrackId string
	//context 日志上下文
	Context interface{}
	// 记录的时间
	Time string
	//日志的名字
	MsgName string

	//访问的server
	ServerHost string
	//用户的uid
	Uid string
}

func init() {
	confObj := conf.New()
	//logger =
	if conf.Runtime() != "prod" {
		//节省性能，线上不显示这个
		logger.EnableFuncCallDepth(true)
	}
	if _, ok := confObj.Dict["logFile"]; ok {
		//logger.SetLogger("console")
		logger.SetLogger("file", `{"filename":"runtime/security.log", "daily":true, "maxdays":10}`)
	}
	trackId = "adfaa"
	logger.SetLogFuncCallDepth(3)
	logger.Async(100)
}

func getLogData() *logFormat {
	timeStr := time.Now().String()
	//.Format("2006-01-02 15:04:05")
	//time.Now().Format("2006-01-02 15:04:05")
	data := logFormat{
		TrackId:    trackId,
		Time:       timeStr,
		ServerHost: "host",
		Uid:        "1",
	}
	return &data
}
func AddLog(msgName string, context interface{}) {

}

/*
	info 级别的日志
*/
//func Info(msgName string, context interface{}) {
func Info(msgName string, context interface{}) {
	logger.Info(string(formatData(msgName, context)))
}

func formatData(msgName string, context interface{}) []byte {
	data := getLogData()
	data.Context = context
	data.MsgName = msgName
	ipJson, _ := json.Marshal(*data)
	return ipJson
}

/*
	warning 级别的日志
*/
func Warn(msgName string, context interface{}) {
	logger.Warning(string(formatData(msgName, context)))
}

/*
	记录 error 级别的日志
*/
func Error(msgName string, context interface{}) {
	logger.Error(string(formatData(msgName, context)))
}
