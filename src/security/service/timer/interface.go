package timer

type EventInterface interface {
	//当一个处理完毕之后，获取下一个
	SetNew(EventRow) int64
	//获取当前最小的一个
	GetMinTime() EventRow
	//SetLuaType() int

	//具体发生时间处理程序
	Process(EventRow) bool
	GetByToken(string) (bool, EventRow)
}

type EventRow struct {
	//当一个处理完毕之后，获取下一个
	Expire_time string
	Id          int64
	Data        string
}
