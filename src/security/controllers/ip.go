package controllers

import (
	"fmt"
	"reflect"
	"security/controllers/base"
	"security/lib/check"
	"security/models/ipModel"
	"security/service/lua"
	"security/service/timer"
	"security/service/user"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type IpController struct {
	//beego.Controller
	base.BaseController
}

//输出数据的Error格式

type ListOutData struct {
	//Data  string
	Data []*ipModel.Black `json:"data"`
	Cnt  int64            `json:"count"`
}

type Out struct {
	Msg  string
	Name string
}

/*
	类似于beforeAction
*/
func (this *IpController) Prepare() {
	//trace::info("ip_parpare", "")
}

/*
	类似于beforeAction
*/
/*
func (this *IpController) Finish() {
	fmt.Println("after_finish")
	//trace::info("ip_parpare", "")
}
*/
/*
	获取ip 的列表
	@params		int		pageSize	每页的数量
	@params		int		pageNo		页码
*/

func (this *IpController) List() {
	//logText := map[string]string{
	if user.IsOnline(this.Ctx) == false {
		this.Output(1840374119, "请您登陆后访问", "")
		return
	}
	tmpSize := this.Input().Get("pageSize")
	tmpNo := this.Input().Get("pageNo")
	ip := this.Input().Get("ip")
	if len(ip) > 0 {
		if check.IsIp(ip) == false {
			this.Output(2140363147, "ip格式错误", "")
			return
		}
	}
	if check.IsNum(tmpSize) == false {
		tmpSize = "2"
	}
	if check.IsNum(tmpNo) == false {
		tmpNo = "1"
	}
	size, _ := strconv.Atoi(tmpSize)
	pageNo, _ := strconv.Atoi(tmpNo)

	offset := (pageNo - 1) * size
	var data []*ipModel.Black

	cond := orm.NewCondition()
	cond = cond.And("status", ipModel.Status_Online)
	if len(ip) > 0 {
		cond = cond.And("ip", ip)
	}
	orm.NewOrm().QueryTable(ipModel.TableName).SetCond(cond).Limit(size, offset).OrderBy("-update_time").All(&data)
	cnt, _ := orm.NewOrm().QueryTable(ipModel.TableName).SetCond(cond).Count()
	outData := ListOutData{
		Data: data,
		Cnt:  cnt,
	}
	this.Output(200, "ok", outData)
}

/*
	添加ip的接口

	@params		string	ip		添加ip
	@params		string	desc
*/

func (this *IpController) Add() {
	if user.IsOnline(this.Ctx) == false {
		this.Output(1840374119, "请您登陆后访问", "")
		return
	}
	ip := this.Input().Get("ip")
	desc := this.Input().Get("note")
	ipArr := strings.Split(ip, ",")
	//分钟的时间，表示时间的长度
	timeFloat, _ := this.GetFloat("time", 720.0)

	timeLength := int32(timeFloat)

	if timeLength < 0 {
		this.Output(1054031066, "请求的数据内容时间错误", "")
		return
	}
	if len(ipArr) > 50 {
		this.Output(1840375100, "添加ip数量超过50", "")
		return
	}

	o := orm.NewOrm()
	tableName := ipModel.GetName()
	for _, val := range ipArr {
		if check.IsIp(val) == false {
			this.Output(184039872, val+" 格式错误", "")
			return
		}
		if o.QueryTable(tableName).Filter("ip", val).Filter("status", ipModel.Status_Online).Exist() == true {
			this.Output(184039812, val+" 已经存在", "")
			return
		}
	}
	userName := user.GetUserName(this.Ctx)
	fmt.Printf("the timeLength is  %d", timeLength)
	if ipModel.Add(ipArr, desc, userName, timeLength) == false {
		this.Output(184039814, "添加失败", "")
		return
	}
	rs := lua.AddIp(ipArr, lua.Type_ip)
	fmt.Println(rs)
	rsError, ok := rs["error"].(map[string]interface{})
	if ok == false {
		this.Output(500, "添加lua异常", "")
	}
	returnCode := rsError["returnCode"].(float64)

	if returnCode == 0 {
		//删除成功
		data := rs["data"].([]interface{})
		cnt := 0
		delIpArr := make([]string, len(data))
		for _, v := range data {
			vstring := v.(string)
			delIpArr[cnt] = vstring
			cnt++
			fmt.Println("send_to_chanel : " + vstring)
			timer.EventCenter.InputChannel <- timer.EventChannel{
				Typ:  timer.Type_Ip,
				Data: vstring,
			}
			//timer.NewIp <- vstring
		}
		if ipModel.UpdateByIps(&orm.Params{"status": ipModel.Status_Online}, delIpArr) == false {
			this.Output(500, "数据库更新失败，添加ng成功", len(delIpArr))
			return
		} else {
			this.Output(200, "ok", len(delIpArr))
		}

	} else {
		this.Output(1540089188, "删除lua失败", rsError)
		return
	}
	this.Output(200, "ok", "error")
	return
}

/*
	添加ip的接口

	@params		string	id		要删除的id
*/

func (this *IpController) Delete() {
	ips := this.Input().Get("ip")

	ipArr := strings.Split(ips, ",")
	if len(ipArr) > 50 {
		this.Output(25403131135, "删除的ip超过50", "")
		return
	}
	for _, val := range ipArr {
		if check.IsIp(val) == false {
			this.Output(184039872, val+" 格式错误", "")
			return
		}
	}
	//o := orm.NewOrm()
	//tableName := ipModel.GetName()
	var dbData []*ipModel.Black
	//var dataLua []*string
	orm.NewOrm().QueryTable(ipModel.GetName()).Filter("ip__in", ipArr).All(&dbData)
	var found = 0
	for _, ip := range ipArr {
		found = 0
		for _, tmp := range dbData {
			if tmp.Ip == ip {
				found = 1
			}
		}
		if found == 0 {
			this.Output(1540314149, ip+" 并不存在", "")
			return
		}
	}

	//dataRes := lua.AddIp(ipArr)
	dataRes := lua.DeleteIps(ipArr, 0)
	fmt.Println(dataRes)
	var ipRes []string
	if dataAns, ok := dataRes["data"]; ok {
		t := reflect.ValueOf(dataAns)
		num := t.Len()
		for i := 0; i < num; i++ {
			if str, ok := t.Index(i).Interface().(string); ok {
				ipRes = append(ipRes, str)
			} else {
				this.Output(17190217, "返回值异常", "")
				return
			}
		}
		if ipModel.DeleteByIps(ipRes) == true {
			this.Output(200, "ok", dataAns)
		} else {
			this.Output(19250219197, "删除成功，数据库更新失败", ipRes)
		}
	}
	this.Output(200, "ok", "")
	return
}
