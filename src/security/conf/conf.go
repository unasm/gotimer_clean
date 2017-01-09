package conf

import (
	"security/conf/dev"
	"security/conf/prod"
)

/*
	dev,beta, prod 三种
*/
func Runtime() string {
	return "dev"
}

//输出数据的Error格式
type Ret struct {
	ReturnCode        int64  `json:"returnCode"`
	ReturnMessage     string `json:"returnMessage"`
	ReturnUserMessage string `json:"returnUserMessage"`
}

//输出的数据
type Res struct {
	Error Ret         `json:"error"`
	Data  interface{} `json:"data"`
	//Data  OutData `json:"data"`
}

var _instance *Config

type Config struct {
	Dict map[string]string
}

func New() *Config {
	if _instance == nil {
		_instance = new(Config)
		_instance.Init()
	}
	return _instance
}

func getCommon() map[string]string {
	config := map[string]string{
		//lua 的相关参数
		"luaAppKey":     "a8dfd32bb2413958ab3874d491d777ce",
		"luaHost":       "http://127.0.0.1",
		"luaSetBlackIp": "/interface/setBlackIPLists",
		"luaSetUri":     "/interface/setConfigByUri",

		// 日志相关的
		"logFile": "test.log",
	}
	return config
}
func (this *Config) Disp() {

}
func (this *Config) Init() {
	var envConfig map[string]string
	env := Runtime()
	this.Dict = getCommon()
	if env == "dev" {
		envConfig = dev.GetConf()
	} else if env == "prod" {
		envConfig = prod.GetConf()
	} else {
		panic("wrong conf runtime")
	}
	for key, val := range envConfig {
		this.Dict[key] = val
		//	fmt.Println(key, "_____", val)
	}
	/*
		for key, val := range this.Dict {
			fmt.Println(key, "_____", val)
		}
	*/
	//this.Dict = comConfig
	/*
		else if env == "beta" {
			return
			envConfig := beta.getConf()
		} else if env == "prod" {
			return beta.getConf()
		}*/
}
