package lua

import (
	"encoding/json"
	"fmt"
	"security/conf"
	"security/lib/check"
	"security/lib/common"
	"security/lib/trace"
	"security/lib/util"
)

const Status_Add = "0"
const Status_Del = "1"

const Type_ip int32 = 0
const Type_udid int32 = 1

type LuaUriJsonAddData struct {
	Uri          string `json:"uri"`
	Type         string `json:"type"`
	Total        string `json:"total"`
	Time         string `json:"time"`
	StrategyType string `json:"strategyType"`
}
type service_lua_log_format struct {
	Url    string
	Res    string
	Params interface{}
}

/*
	添加ip列表

	@params		ip	[]string			要添加的数组列表
	@params		typ	int					类型， 0是ip， 1 是udid
	@return		map[string]interface{}	返回的数组列表
*/
func AddIp(ips []string, typ int32) map[string]interface{} {
	config := conf.New()
	url := config.Dict["luaHost"] + config.Dict["luaSetBlackIp"]

	ipJson, _ := json.Marshal(ips)
	data := map[string]string{
		"ipLists": string(ipJson),
		"status":  Status_Add,
		"type":    util.IntToStr(typ),
	}

	res := doRequest(url, data)
	return res
}

/*
	获取ip列表是否在用
*/
func GetList(ip []string) {

}

/*
 添加uri
*/
func DelUri(data *[]LuaUriJsonAddData) map[string]interface{} {
	ipJson, _ := json.Marshal(data)

	config := conf.New()
	url := config.Dict["luaHost"] + config.Dict["luaSetUri"]
	uriData := map[string]string{
		"status": Status_Del,
		"config": string(ipJson),
	}
	//fmt.Println(uriData)
	return doRequest(url, uriData)
}

/*
 添加uri
*/
func AddUri(data *[]LuaUriJsonAddData) map[string]interface{} {
	ipJson, _ := json.Marshal(data)

	config := conf.New()
	url := config.Dict["luaHost"] + config.Dict["luaSetUri"]
	uriData := map[string]string{
		"status": Status_Add,
		"config": string(ipJson),
	}
	//fmt.Println(uriData)
	return doRequest(url, uriData)
}

/*
	从lua删除ip
*/
func DeleteIps(ips []string, typ int32) map[string]interface{} {
	config := conf.New()
	url := config.Dict["luaHost"] + config.Dict["luaSetBlackIp"]

	ipJson, _ := json.Marshal(ips)
	data := map[string]string{
		"ipLists": string(ipJson),
		"status":  Status_Del,
		"type":    util.IntToStr(typ),
	}

	res := doRequest(url, data)
	return res
}

/*
	发送http请求
*/
func doRequest(url string, data map[string]string) map[string]interface{} {
	config := conf.New()
	fmt.Println("config____")
	fmt.Println(config.Dict)
	data["appkey"] = config.Dict["luaAppKey"]
	addStr := common.HttpPost(url, data)

	trace.Info("service_lua_doRequest", service_lua_log_format{
		Url:    url,
		Params: data,
		Res:    addStr,
	})
	//fmt.Println(addStr)
	var jsonData map[string]interface{}
	//_ := json.Unmarshal([]byte(addStr), &jsonData)
	//json.Unmarshal([]byte(addStr), &jsonData)
	err := json.Unmarshal([]byte(addStr), &jsonData)
	check.Err(err)
	return jsonData
}
