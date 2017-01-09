package timer

import (
	"fmt"
	"security/lib/trace"
	"security/lib/util"
	"time"
)

var NO_EVENT_INTERVAL = 10

const Type_Ip = "ip"

//即时，到时候唤醒
//var timer time.Timer

/*
	需要一个定时器
	需要一个变量，记录下一次醒来的时间
*/

func init() {
	fmt.Println("timer_________started")
	//EventCenter.Register("ip", &timer.Timer_ip{})
	//NewIp = make(chan string, 1)
}

type EventChannel struct {
	Data string
	Typ  string
}

type _eventCenter struct {
	//需要一个信道，与外界通信，获取外界新的时间，
	InputChannel chan EventChannel

	//下一次唤醒的时间戳
	LastWakeTime int64
	// 下一次下一次唤醒的时候的typ
	NowType string

	//醒来之后，处理的数据
	WakeRow EventRow

	//具体处理的对象和时间调度器的映射
	EventMap map[string]*EventInterface
}

var (
	EventCenter = &_eventCenter{
		InputChannel: make(chan EventChannel, 1),
		LastWakeTime: int64(0),
		//WakeRow:      make(EventRow, 1),
		EventMap: make(map[string]*EventInterface),
	}
)

//Event.RegisterModel

func (this *_eventCenter) Register(typ string, event EventInterface) (bool, error) {
	if _, ok := this.EventMap[typ]; ok {
		return false, fmt.Errorf("compare datetime miss format")
	}
	this.EventMap[typ] = &event
	return true, nil
}

func (this *_eventCenter) GetMinRow() (string, EventRow) {
	var (
		timeUnix int64
		minRow   EventRow
		eventRow EventRow
		minType  string
	)
	minRow.Expire_time = "2038-01-01 19:34:39"
	minType = ""
	for typ, event := range this.EventMap {
		eventRow = (*event).GetMinTime()
		timeUnix = util.GetUnix(eventRow.Expire_time)
		if timeUnix != 0 && timeUnix < util.GetUnix(minRow.Expire_time) {
			minRow = eventRow
			minType = typ
		}
		//fmt.Println(timeUnix)
		//fmt.Println("timer_start : " + typ)
	}

	return minType, minRow
}

func (this *_eventCenter) SetMinTime(minType string, minRow EventRow) int64 {
	if minType != "" {
		nowEvent := this.EventMap[minType]
		nextTime := (*nowEvent).SetNew(minRow)
		if nextTime == -1 {
			nextTime = int64(NO_EVENT_INTERVAL)
			//return int64(NO_EVENT_INTERVAL)
			this.WakeRow = EventRow{Id: -1}
		} else {
			this.WakeRow = minRow
		}
		this.LastWakeTime = int64(time.Now().Unix() + nextTime)
		return nextTime
	} else {
		this.LastWakeTime = time.Now().Unix() + int64(NO_EVENT_INTERVAL)
		return int64(NO_EVENT_INTERVAL)
		//return int64(time.Now().Unix() + 60)
	}
}

/*
	用于清理过期的数据

	在开始的时候，先从数据库获取到所有的数据，然后根据最小的一个时间去 设置定时器,按时清除
*/
func (this *_eventCenter) Start() {
	//下一次唤醒的时间
	minType, minRow := this.GetMinRow()
	startTime := this.SetMinTime(minType, minRow)
	fmt.Printf("next_wake 1 : %d\n", startTime)
	//startTime := int64(2)
	//var timer *time.Timer
	timer := time.NewTimer(time.Duration(startTime) * time.Second)

	trace.Info("clear_pass_start", map[string]string{
		"startTime": util.IntToStr(int32(startTime)),
	})
	for {
		select {
		case input := <-this.InputChannel:
			trace.Info("timer_get_data", map[string]string{
				"newIp": input.Data,
				"type":  input.Typ,
			})
			//如果有新的数据插入的话,检查是否更加靠前

			nowEvent := this.EventMap[input.Typ]
			status, rowData := (*nowEvent).GetByToken(input.Data)
			if status != true || rowData.Id == -1 {
				//查找10S 没有找到数据,默认数据不存在
				break
			}
			stamp := util.GetUnix(rowData.Expire_time)

			if stamp > this.LastWakeTime && this.LastWakeTime > time.Now().Unix() {
				//晚于 最近一次唤醒时间，不管
				break
			}

			startTime = this.SetMinTime(input.Typ, rowData)

			//fmt.Printf("reSetTime %d \n", startTime)
			timer.Reset(time.Duration(startTime) * time.Second)
			break
		case cleanTime := <-timer.C:
			fmt.Printf("cleanTime %d \n", cleanTime.Unix())
			if this.WakeRow.Id > 0 {
				//如果有内容， 真正执行代码
				nowEvent := this.EventMap[minType]
				fmt.Println(this.WakeRow)
				(*nowEvent).Process(this.WakeRow)
				this.WakeRow = EventRow{Id: -1}
			}
			//获取最小的数据
			minType, minRow = this.GetMinRow()
			startTime = this.SetMinTime(minType, minRow)
			timer.Reset(time.Duration(startTime) * time.Second)

			break
		}
		trace.Info("clean_pass_ending", "")
	}
	defer func() {
		trace.Info("ClearPassData_end", "")
		timer.Stop()
	}()
}
