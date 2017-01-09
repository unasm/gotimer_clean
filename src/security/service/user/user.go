package user

import (
	"encoding/base64"
	"fmt"
	"security/conf"
	"security/controllers/base"
	"security/lib/check"
	"security/lib/common"
	"sort"
	"strings"

	"github.com/astaxie/beego/context"
)

//login 的返回值
type LoginRetData struct {
	Status int64
	Level  int64
	//定义成map , 更加灵活一些;
	Cookies map[string]string
}
type LoginRes struct {
	Error base.Ret
	Data  LoginRetData
}

/*
	获取用户的登陆信息
*/
func GetUserName(ctx *context.Context) string {
	headerAuth := strings.Trim(ctx.Request.Header.Get("Authorization"), " \n\r\t")
	//headerAuth = "Basic cGFueXU6cGFueXU="
	if len(headerAuth) <= 4 {
		panic("登陆状态异常")
	}
	arr := strings.Split(headerAuth, "Basic ")
	if len(arr) != 2 {
		panic("登陆状态异常")
	}
	data, err := base64.StdEncoding.DecodeString(strings.Trim(arr[1], " \n\r\t"))
	check.Err(err)

	arr = strings.Split(string(data), ":")
	if len(arr) != 2 {
		panic("登陆状态异常")
	}
	return arr[0]
}

func Login(userName string, password string) {
	config := conf.New()
	systemId := config.Dict["userSystemId"]
	data := map[string]string{
		"name":     userName,
		"password": password,
		"systemId": systemId,
	}
	fmt.Println(data)
}

/*
	判断用户是否登陆
*/
func IsOnline(ctx *context.Context) bool {
	//conf := conf.New()
	user_name := ctx.GetCookie("user_name")
	token := ctx.GetCookie("token")
	fmt.Println(user_name)
	fmt.Println(token)
	return true
}

/*
	要生成签名的数组
*/
func createSign(params map[string]string, appKey string) string {
	keys := make([]string, len(params))
	cnt := 0
	for k, _ := range params {
		keys[cnt] = k
		cnt++
	}
	sort.Sort(sort.StringSlice(keys))
	md5Str := ""
	for _, idx := range keys {
		md5Str = md5Str + params[idx] + "|"
	}
	md5Str = md5Str + appKey
	fmt.Println(md5Str)
	return common.GetMd5(md5Str)
}

/*
	发送http请求
*/
func doRequest(url string, data map[string]string) string {
	config := conf.New()
	data["appKey"] = config.Dict["userAppKeyPublic"]
	data["sign"] = createSign(data, config.Dict["userAppKeyPrivate"])
	fmt.Println(data)
	fmt.Println(url)
	resStr := common.HttpPost(url, data)
	return resStr
}
