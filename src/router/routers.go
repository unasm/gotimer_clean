package routers

import (
	"fmt"
	"security/controllers"

	"github.com/astaxie/beego"
)

func init() {
	fmt.Println("rougter.go")
	beego.Router("/", &controllers.IpController{})
	//beego.Router("/", &controllers.MainController{})
}
