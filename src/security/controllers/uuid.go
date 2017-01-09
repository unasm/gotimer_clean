package controllers

import (
	"reflect"
	"security/controllers/base"
	"security/lib/check"
	"security/models/udid"
	"security/service/lua"
	"security/service/timerUdid"
	"security/service/user"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type UdidController struct {
	//beego.Controller
	base.BaseController
}

//输出数据的Error格式

type UdidOutData struct {
	//Data  string
	Data []*udid.Udid `json:"data"`
	Cnt  int64        `json:"count"`
}

/*
	添加ip的接口

	  `way` '阻断方式,1.图形验证码、2.短信验证码、3.直接返回' 目前只支持两种，3，1 对应的lua 是0(阻断)，1(图形验证码)

	@params		string	ip		添加ip
	@params		string	desc
*/

func (this *UdidController) Add() {
	/*
		if user.IsOnline(this.Ctx) == false {
			this.Output(1840374119, "请您登陆后访问", "")
			return
		}
	*/
	udidStr := this.Input().Get("udid")
	desc := this.Input().Get("note")
	//udidArr := strings.Split(udids, ",")
	//分钟的时间，表示时间的长度
	timeFloat, err := this.GetFloat("time", 720.0)

	timeLength := int32(timeFloat)

	if timeLength < 0 {
		this.Output(1054031066, "请求的数据内容时间错误", "")
		return
	}

	o := orm.NewOrm()
	tableName := udid.GetName()
	if o.QueryTable(tableName).Filter("udid", udidStr).Filter("status", udid.Status_Online).Exist() == true {
		this.Output(184039812, udidStr+" 已经存在", "")
		return
	}
	userName := "unasm"
	//userName := user.GetUserName(this.Ctx)
	var insertId int32
	if insertId, err = udid.Adds(udidStr, desc, userName, timeLength); err != nil {
		this.Output(184039814, "添加失败", "")
		return
	}
	rs := lua.AddIp([]string{udidStr}, lua.Type_udid)
	rsError, ok := rs["error"].(map[string]interface{})
	if ok == false {
		this.Output(500, "添加lua异常", "")
	}
	returnCode := rsError["returnCode"].(float64)

	if returnCode == 0 {
		//删除成功
		data := rs["data"].([]interface{})
		cnt := 0
		udidSucc := make([]string, len(data))
		if len(data) > 0 {
			for _, v := range data {
				vstring := v.(string)
				udidSucc[cnt] = vstring
				timerUdid.NewIp <- vstring
				if udid.UpdateById(&orm.Params{"status": udid.Status_Online}, insertId) == false {
					this.Output(500, "数据库更新失败，添加ng成功, 请重试", "")
				} else {
					this.Output(200, "ok", udidSucc)
				}
				break
			}
		} else {
			this.Output(16500100, "更新失败", "")
		}
		return
	} else {
		this.Output(1540089188, "删除lua失败", rsError)
		return
	}
	this.Output(200, "ok", "error")
	return
}

func (this *UdidController) List() {
	if user.IsOnline(this.Ctx) == false {
		this.Output(1840374119, "请您登陆后访问", "")
		return
	}
	tmpSize := this.Input().Get("pageSize")
	tmpNo := this.Input().Get("pageNo")
	udidInput := this.Input().Get("udid")

	if check.IsNum(tmpSize) == false {
		tmpSize = "2"
	}
	if check.IsNum(tmpNo) == false {
		tmpNo = "1"
	}
	size, _ := strconv.Atoi(tmpSize)
	pageNo, _ := strconv.Atoi(tmpNo)

	offset := (pageNo - 1) * size
	var data []*udid.Udid

	cond := orm.NewCondition()
	cond = cond.And("status", udid.Status_Online)
	if len(udidInput) > 0 {
		cond = cond.And("udid", udidInput)
	}
	orm.NewOrm().QueryTable(udid.GetName()).SetCond(cond).Limit(size, offset).OrderBy("-update_time").All(&data)
	cnt, _ := orm.NewOrm().QueryTable(udid.GetName()).SetCond(cond).Count()
	outData := UdidOutData{
		Data: data,
		Cnt:  cnt,
	}
	this.Output(200, "ok", outData)
}

/*
	添加ip的接口

	@params		string	id		要删除的id
*/

func (this *UdidController) Delete() {
	udids := this.Input().Get("udid")

	if len(udids) == 0 {
		this.Output(15403161208, "删除数据为空", "")
		return
	}
	udidArr := strings.Split(udids, ",")
	if len(udidArr) > 50 {
		this.Output(25403131135, "删除的ip超过50", "")
		return
	}

	var dbData []*udid.Udid
	//var dataLua []*string
	orm.NewOrm().QueryTable(udid.GetName()).Filter("udid__in", udidArr).All(&dbData)
	var found = 0
	for _, val := range udidArr {
		found = 0
		for _, tmp := range dbData {
			if tmp.Udid == val {
				found = 1
			}
		}
		if found == 0 {
			this.Output(1540314149, val+" 并不存在", "")
			return
		}
	}

	//dataRes := lua.AddIp(ipArr)
	dataRes := lua.DeleteIps(udidArr, lua.Type_udid)
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
		if udid.DeleteByUdids(ipRes) == true {
			this.Output(200, "ok", dataAns)
		} else {
			this.Output(19250219197, "删除成功，数据库更新失败", ipRes)
		}
	}
	this.Output(200, "ok", "")
	return
}
