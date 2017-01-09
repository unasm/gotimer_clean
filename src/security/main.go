package main

import (
	"fmt"
	"security/controllers"
	"security/service/timer"
	"security/service/timerUdid"
	"time"

	"github.com/astaxie/beego"
)

func main() {
	time.LoadLocation("America/New_York")

	go timer.ClearPassData()
	go timerUdid.ClearPassData()
	beego.BConfig.RunMode = "prod"
	//beego.BConfig.RecoverPanic = true
	beego.Router("/black/list", &controllers.IpController{}, "*:List")
	beego.Router("/black/add", &controllers.IpController{}, "*:Add")
	beego.Router("/black/del", &controllers.IpController{}, "*:Delete")

	beego.Router("/uri/list", &controllers.UriController{}, "*:List")
	beego.Router("/uri/add", &controllers.UriController{}, "*:Add")
	beego.Router("/uri/del", &controllers.UriController{}, "*:Delete")

	beego.Router("/udid/list", &controllers.UdidController{}, "*:List")
	beego.Router("/udid/add", &controllers.UdidController{}, "*:Add")
	beego.Router("/udid/del", &controllers.UdidController{}, "*:Delete")

	beego.Run()
	fmt.Println("done")
}
