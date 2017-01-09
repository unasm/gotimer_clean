package controllers

import (
	"fmt"
	"security/controllers/base"
	"security/lib/check"
	"security/models/rule"
	"security/service/lua"
	"security/service/user"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

var wayMap = map[int64]int64{
	3: 0,
	1: 1,
}

type UriController struct {
	//beego.Controller
	base.BaseController
}

type ReturnError struct {
	ReturnCode        int64
	ReturnMessage     string
	ReturnUserMessage string
}

//输出数据的Error格式

type ListUriOutData struct {
	//Data  string
	Data []*rule.Rule `json:"data"`
	Cnt  int64        `json:"count"`
}

/*
	删除uri

	@params		string	id		要删除的id
	@todo  修改成批量删除
*/

func (this *UriController) Delete() {
	//uris := make([]string, 50)
	//err := this.Ctx.Input.Bind(&uris, "uri")
	inputUri := strings.Trim(this.Input().Get("uri"), " ")
	uris := strings.Split(inputUri, ",")
	lenUris := len(uris)
	if lenUris > 50 {
		this.Output(1540350224, "删除数据过多，超过50条", "")
		return
	}
	if lenUris == 0 {
		this.Output(1540357225, "提交数据为空", "")
		return
	}
	fmt.Println(lenUris)
	var dbRow []*rule.Rule
	var tmpData = make([]lua.LuaUriJsonAddData, 1)
	fmt.Println(inputUri)
	fmt.Println(uris)
	for cnt, inputUri := range uris {
		inputUri = strings.Trim(inputUri, " ")
		fmt.Println(inputUri)
		if check.IsUri(inputUri) == false {
			this.Output(1640264228, inputUri+" 格式错误", "")
			return
		}
		orm.NewOrm().QueryTable(rule.GetName()).Filter("uri", inputUri).
			Filter("status", rule.Status_Online).All(&dbRow)
		if len(dbRow) != 1 {
			this.Output(184039812, inputUri+" 不存在,或者未生效", "")
			return
		}
		fmt.Println(cnt)
		//dbIds[cnt] = dbRow[0].Id
		//addUriData[cnt] = lua.LuaUriJsonAddData{

		strategy, ok := wayMap[int64(dbRow[0].Way)]
		if !ok {
			this.Output(15403211241, "阻断方式错误", "")
			return
		}
		tmpData[0] = lua.LuaUriJsonAddData{
			Uri:          inputUri,
			Type:         strconv.Itoa(int(dbRow[0].Dim_id)),
			Total:        strconv.Itoa(int(dbRow[0].Times)),
			Time:         strconv.Itoa(int(dbRow[0].Expire)),
			StrategyType: strconv.Itoa(int(strategy)),
		}

		rs := lua.DelUri(&tmpData)
		rsError, ok := rs["error"].(map[string]interface{})
		if ok == false {
			this.Output(500, "添加lua异常", "")
		}
		returnCode := rsError["returnCode"].(float64)

		if returnCode == 0 {
			//删除成功
			if rule.Delete(dbRow[0].Id) == false {
				this.Output(1650083190, "删除lua成功，修改db失败", "")
				return
			}
		} else {
			this.Output(1540089188, "删除lua失败", rsError)
			return
		}
	}
	this.Output(200, "ok", "")
	return
}

/*
	获取uri列表的接口
*/
func (this *UriController) List() { //size := this.Input().Get("pageSize").Int()
	size, err := this.GetInt("pageSize")
	check.Err(err)
	pageNo, err := this.GetInt("pageNo")
	check.Err(err)
	//pageNo := this.Input().Get("pageNo").Int()
	uri := this.GetString("uri")
	fmt.Println(uri)
	fmt.Println(pageNo)
	fmt.Println(size)
	offset := (pageNo - 1) * size
	if offset < 0 {
		offset = 0
	}

	cond := orm.NewCondition()

	cond = cond.And("status", rule.Status_Online)
	if len(uri) > 0 {
		cond = cond.And("uri", uri)
	}
	var data []*rule.Rule
	orm.NewOrm().QueryTable(rule.GetName()).SetCond(cond).Limit(size, offset).OrderBy("-update_time").All(&data)
	cnt, _ := orm.NewOrm().QueryTable(rule.GetName()).SetCond(cond).Count()
	outData := ListUriOutData{
		Data: data,
		Cnt:  cnt,
	}
	this.Output(200, "ok", outData)
	return
}

/*
	添加ip的接口

	  `way` '阻断方式,1.图形验证码、2.短信验证码、3.直接返回' 目前只支持两种，3，1 对应的lua 是0(阻断)，1(图形验证码)

	@params		string	ip		添加ip
	@params		string	desc
*/

func (this *UriController) Add() {

	//fmt.Println(user.GetUserName(this.Ctx))
	if user.IsOnline(this.Ctx) == false {
		this.Output(1840374119, "请您登陆后访问", "")
		return
	}

	//要处理的uri
	uri := strings.Trim(this.Input().Get("uri"), " ")

	//原因
	reason := strings.Trim(this.Input().Get("reason"), " ")
	//次数
	tmpInput := this.Input().Get("times")
	times, err := strconv.Atoi(tmpInput)
	check.Err(err)

	// 时间，每天，每小时,按分钟计算
	tmpInput = this.Input().Get("expire")
	expire, err := strconv.Atoi(tmpInput)
	check.Err(err)

	//添加阻断方式, 0, 1
	tmpInput = this.Input().Get("way")
	way, err := strconv.Atoi(tmpInput)
	check.Err(err)

	//统计的维度，1 : memberId+ip，2 : uuid,0 : memberId, 3 : ip
	Dim_id, err := this.GetInt("dim")
	check.Err(err)
	if Dim_id > 3 || Dim_id < 0 {
		this.Output(25540318195, fmt.Sprintf("%d", Dim_id)+" 统计条件错误", "")
		return
	}

	if check.IsUri(uri) == false {
		this.Output(24403119165, uri+" 格式错误", "")
		return
	}
	if orm.NewOrm().QueryTable(rule.GetName()).Filter("uri", uri).Filter("status", rule.Status_Online).Exist() == true {
		this.Output(184039812, uri+" 已经存在", "")
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	data := rule.Rule{
		Uri:         uri,
		Dim_id:      int64(Dim_id),
		User:        user.GetUserName(this.Ctx),
		Times:       int64(times),
		Expire:      int64(expire),
		Way:         int64(way),
		Reason:      reason,
		Status:      rule.Status_Init,
		Create_time: timeStr,
		Update_time: timeStr,
	}
	uriId := rule.Add(&data)
	strategy, ok := wayMap[int64(way)]
	if !ok {
		this.Output(15403211241, "阻断方式错误", "")
		return
	}
	//strategy =
	addUriData := []lua.LuaUriJsonAddData{
		{
			Uri:          uri,
			Type:         strconv.Itoa(Dim_id),
			Total:        strconv.Itoa(times),
			Time:         strconv.Itoa(expire),
			StrategyType: strconv.Itoa(int(strategy)),
		},
	}
	fmt.Println(addUriData)
	rs := lua.AddUri(&addUriData)

	rsError, ok := rs["error"].(map[string]interface{})
	if ok == false {
		this.Output(500, "添加lua异常", "")
	}
	returnCode := rsError["returnCode"].(float64)
	if int(returnCode) == 0 {
		if rule.Update(&orm.Params{"status": rule.Status_Online}, uriId) == false {
			this.Output(500, "数据库更新失败，添加ng成功", "")
		} else {
			this.Output(200, "ok", "")
		}
		return
	}

	this.Output(14500185187, "调用ng前置异常", rsError["returnCode"])
	return
}
