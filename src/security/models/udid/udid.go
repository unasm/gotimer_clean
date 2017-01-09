package udid

import (
	"fmt"
	"security/lib/check"
	"security/models/connection"
	"time"

	"github.com/astaxie/beego/orm"
)

const TableName = "udid"

//刚刚添加
const Status_Init = 1

//线上在用
const Status_Online = 3

//添加失败
const Status_Failed = 7

//删除
const Status_Del = 14

type Udid struct {
	Id          int64 `-1`
	Udid        string
	Status      int64
	Desc        string
	Creator     string
	Create_time string
	Update_time string
	Expire_time string
}

func init() {
	connection.Init()
	orm.RegisterModel(new(Udid))
	//orm.RunSyncdb("default", false, true)
}

//获取数据库的表名
func GetName() string {
	return TableName
}

/*
	插入数据库数据
*/
func Adds(udidStr string, desc string, creator string, timeLength int32) (int32, error) {
	//data Black
	//tableName := GetName()
	timer := time.Now()
	timeStr := timer.Format("2006-01-02 15:04:05")

	o := orm.NewOrm()
	err := o.Begin()

	//换算成秒
	timeLength = timeLength * 60
	expireTime := int32(timer.Unix()) + timeLength
	expireTimeFormat := time.Unix(int64(expireTime), 0).Format("2006-01-02 15:04:05")

	dataRow := Udid{
		Creator:     creator,
		Status:      Status_Init,
		Udid:        udidStr,
		Desc:        desc,
		Create_time: timeStr,
		Update_time: timeStr,
		Expire_time: expireTimeFormat,
	}
	Id, err := o.Insert(&dataRow)
	if err != nil {
		o.Rollback()
		check.Err(err)
	}
	err = o.Commit()
	check.Err(err)
	return int32(Id), nil
}

/*
	标记删除
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
func DeleteByUdids(udids []string) bool {
	params := orm.Params{
		"status": Status_Del,
	}
	return UpdateByIps(&params, udids)
}

/*
	查找最小的一一行数据
*/
func GetMinTime() Udid {
	o := orm.NewOrm()
	//通过默认值区分是否获得了数组
	row := Udid{
		Id: -1,
	}
	//select udid.* from udid inner join  (select id from udid where expire_time  > 0 && status = 3 order by expire_time  limit 1 ) t2 on udid.id  = t2.id\G;
	sql := "SELECT t1.* FROM  " + TableName + " t1 inner join (select id from " + TableName + " where `expire_time` > 0 && `status` = " + fmt.Sprintf("%d", Status_Online) + " order by expire_time limit 1) t2 on t1.id = t2.id"
	//fmt.Println(sql)
	o.Raw(sql).QueryRow(&row)
	return row
}

/*
   根据udid 更新
*/
func UpdateByIps(data *orm.Params, udids []string) bool {
	o := orm.NewOrm()
	err := o.Begin()
	for _, val := range udids {
		_, err := o.QueryTable(TableName).Filter("udid", val).Update(*data)
		if err != nil {
			o.Rollback()
			check.Err(err)
			return false
		}
	}
	err = o.Commit()
	check.Err(err)
	return true
}

//只有一列,online状态的只有一列
func GetByUdid(udid string) Udid {
	o := orm.NewOrm()
	tableName := GetName()
	var data = Udid{
		Id: -1,
	}
	o.QueryTable(tableName).Filter("udid", udid).Filter("status", Status_Online).One(&data)
	return data
}

func UpdateById(data *orm.Params, id int32) bool {
	o := orm.NewOrm()
	o.Begin()
	_, err := o.QueryTable(GetName()).Filter("id", id).Update(*data)
	if err != nil {
		o.Rollback()
		check.Err(err)
		return false
	}
	err = o.Commit()
	check.Err(err)
	return true
}
