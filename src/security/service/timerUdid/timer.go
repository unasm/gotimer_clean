package timerUdid

import (
	"fmt"
	"reflect"
	"security/lib/trace"
	"security/lib/util"
	"security/models/udid"
	"security/service/lua"
	"time"
)

//需要一个信道，与外界通信，获取外界新的时间，
var NewIp chan string

//下一次唤醒的时间戳
var lastWakeTime int64

//醒来之后，处理的数据
var wakeRow udid.Udid

//即时，到时候唤醒
//var timer time.Timer

/*
	需要一个定时器
	需要一个变量，记录下一次醒来的时间
*/

func init() {
	fmt.Println("timer_________started")
	NewIp = make(chan string, 1)
}

//真正执行删除的动作
func Process(row udid.Udid) bool {
	udidArr := make([]string, 1)
	udidArr[0] = row.Udid
	dataRes := lua.DeleteIps(udidArr, lua.Type_udid)
	if dataAns, ok := dataRes["data"]; ok {
		fmt.Println(dataAns)
		t := reflect.ValueOf(dataAns)
		num := t.Len()
		if num != 1 {
			trace.Error("timerProcessGetDataUdid", map[string]interface{}{
				"udid": row.Udid,
				"res":  dataAns,
			})
			return false
		}
		if _, ok := t.Index(0).Interface().(string); ok {
			return udid.DeleteByUdids(udidArr)
		}
	}
	return false
}

/*
	设置下次唤醒
	每次都从数据库中获取到新的数据，防止走入更大的偏差，丢失数据
*/
func setNew() int64 {
	wakeRow = udid.GetMinTime()
	var rs int64
	if wakeRow.Id == -1 {
		//如果没有时间可设置,十分钟之后检查一次
		rs = 600
		return 0
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
		lastWakeTime = rs + nowTime
	}
	fmt.Println("next_wake_length_udid : ", rs, "\t wakeRow : ", wakeRow.Udid)
	return rs
}

/*
	用于清理过期的数据

	在开始的时候，先从数据库获取到所有的数据，然后根据最小的一个时间去 设置定时器,按时清除
*/
func ClearPassData() {
	//下一次唤醒的时间
	timer := time.NewTimer(time.Duration(setNew()) * time.Second)
	var newIp string
	var tmpRow udid.Udid
	for {
		select {
		case newIp = <-NewIp:
			trace.Info("timer_get_data", map[string]string{
				"newIp": newIp,
			})
			//如果有新的数据插入的话,检查是否更加靠前
			for cnt := 5; cnt >= 0; cnt-- {
				tmpRow = udid.GetByUdid(newIp)
				fmt.Println("start_ending : " + newIp)
				fmt.Println(tmpRow)
				if tmpRow.Id != -1 {
					break
				}
				time.Sleep(time.Second)
			}
			if tmpRow.Id == -1 {
				trace.Info("timer_no_data", map[string]string{
					"newIp": newIp,
				})
				//查找10S 没有找到数据,默认数据不存在
				break
			}
			stamp := util.GetUnix(tmpRow.Expire_time)

			fmt.Printf("yes_data %d \n", stamp)
			fmt.Printf("lastWakeTime %d \n", lastWakeTime)
			if stamp > lastWakeTime && lastWakeTime > time.Now().Unix() {
				//晚于 最近一次唤醒时间，不管
				break
			}
			/*
				lastWakeTime = stamp
				wakeRow = tmpRow
				//有数据，并且 //timer.Reset(2 * time.Second)
				fmt.Printf("length is  %d \n", (stamp - time.Now().Unix()))
				timer.Reset(time.Duration((stamp - time.Now().Unix())) * time.Second)
			*/
			timer.Reset(time.Duration(setNew()) * time.Second)
			break
		case cleanTime := <-timer.C:
			//fmt.Println(cleanTime.Unix())
			fmt.Println("cleaNTime", cleanTime)
			trace.Info("wake_to_process", map[string]interface{}{
				"time.Unix": cleanTime.Unix(),
				"ip":        wakeRow.Udid,
			})
			//真正执行代码
			Process(wakeRow)
			timer.Reset(time.Duration(setNew()) * time.Second)
			break
		}
		trace.Info("clean_pass_ending", "")
	}
	defer func() {
		trace.Info("ClearPassData_end", "")
		timer.Stop()
	}()
}
