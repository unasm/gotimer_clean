package main

import (
	"fmt"
	_ "security/routers"

	"github.com/astaxie/beego"
)

/*type MainController struct {*/
//beego.Controller
//}

//func (this *MainController) Get() {
//this.Data["Website"] = "beego.me"
//this.Data["Email"] = "astaxie@gmail.com"
//this.TplName = "index.tpl"
/*}*/

func main() {
	fmt.Println("vim-go")
	//beego.Router("/", &controllers.MainController{})
	beego.Run()
}
