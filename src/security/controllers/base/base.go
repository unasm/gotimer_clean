package base

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
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

/*
	输出 格式化的内容
	@params		int64			errno	错误码
	@params		string			errMsg	错误信息
	@params		interface{}		outData		要输出的错误的内容
*/
func (this *BaseController) Output(errno int64, errMsg string, outData interface{}) {
	if errno == 200 {
		errno = 0
	}
	mystruct := Res{
		Error: Ret{
			ReturnCode:        errno,
			ReturnMessage:     "__",
			ReturnUserMessage: errMsg,
		},
		Data: outData,
	}
	this.Data["json"] = mystruct
	this.ServeJSON()
}
