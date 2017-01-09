package timer

import (
	"fmt"
	"reflect"
	"security/lib/trace"
	"security/lib/util"
	"security/models/ipModel"
	"security/service/lua"
	"time"
)

type Timer_ip struct {
	lastWakeTime int64
	wakeRow      ipModel.Black
}

//真正执行删除的动作
//func (t *Timer_ip) Process(row ipModel.Black) bool {
func (t *Timer_ip) Process(row EventRow) bool {
	ipArr := make([]string, 1)
	ipArr[0] = row.Data
	dataRes := lua.DeleteIps(ipArr, 0)
	if dataAns, ok := dataRes["data"]; ok {
		t := reflect.ValueOf(dataAns)
		num := t.Len()
		if num != 1 {
			trace.Error("timerProcessGetData", map[string]interface{}{
				"ip":  row.Data,
				"res": dataAns,
			})
			return false
		}
		if _, ok := t.Index(0).Interface().(string); ok {
			return ipModel.DeleteByIps(ipArr)
		}
	}
	return false
}

/*
	设置下次唤醒
	每次都从数据库中获取到新的数据，防止走入更大的偏差，丢失数据
*/
func (this *Timer_ip) SetNew(wakeRow EventRow) int64 {
	//this.wakeRow = this.GetMinTime()
	var rs int64
	if wakeRow.Id == -1 {
		//如果没有时间可设置,十分钟之后检查一次
		return -1
		//return int64(600)
	} else {
		nowTime := time.Now().Unix()
		length := util.GetUnix(wakeRow.Expire_time) - nowTime
		if length > 0 {
			rs = length
			//return length
		} else {
			rs = 1
		}
		//lastWakeTime = rs + nowTime
	}
	fmt.Println("next_wake_length", rs)
	return rs
}

func (t *Timer_ip) GetByToken(newIp string) (bool, EventRow) {
	var tmpData ipModel.Black
	for cnt := 5; cnt >= 0; cnt-- {
		tmpData = ipModel.GetByIp(newIp)
		if tmpData.Id != -1 {
			return true, EventRow{
				Expire_time: tmpData.Expire_time,
				Id:          tmpData.Id,
				Data:        tmpData.Ip,
			}
		}
		time.Sleep(time.Second)
	}
	return false, EventRow{Id: -1}
}

func (t *Timer_ip) GetMinTime() EventRow {
	model := ipModel.GetMinTime()
	return EventRow{
		Expire_time: model.Expire_time,
		Id:          model.Id,
		Data:        model.Ip,
	}
}
