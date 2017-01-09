package ipModel

import (
	"fmt"
	"security/lib/check"
	"security/models/connection"
	"time"

	"github.com/astaxie/beego/orm"
)

const TableName = "black"

type Black struct {
	Id          int64 `-1`
	Ip          string
	Status      int64
	Desc        string
	Creator     string
	Create_time string
	Update_time string
	Expire_time string
}

//刚刚添加
const Status_Init = 1

//线上在用
const Status_Online = 3

//添加失败
const Status_Failed = 7

//删除
const Status_Del = 14

func init() {
	connection.Init()
	orm.RegisterModel(new(Black))
	//orm.RunSyncdb("default", false, true)
}

//获取数据库的表名
func GetName() string {
	return TableName
}

/*
	插入数据库数据
*/
func Add(ips []string, desc string, creator string, timeLength int32) bool {
	//data Black

	tableName := GetName()
	timer := time.Now()
	timeStr := timer.Format("2006-01-02 15:04:05")

	var dbData []*Black
	o := orm.NewOrm()
	err := o.Begin()

	//换算成秒
	timeLength = timeLength * 60
	expireTime := int32(timer.Unix()) + timeLength
	expireTimeFormat := time.Unix(int64(expireTime), 0).Format("2006-01-02 15:04:05")

	params := orm.Params{}
	for _, val := range ips {
		dataRow := Black{
			Ip: val,
		}

		o.QueryTable(tableName).Filter("ip", val).All(&dbData)
		if len(dbData) > 0 {
			len := len(dbData)
			for i := 0; i < len; i++ {
				params = orm.Params{
					"creator":     creator,
					"status":      Status_Init,
					"desc":        desc,
					"create_time": timeStr,
					"update_time": timeStr,
					"expire_time": expireTimeFormat,
				}
				Update(&params, dbData[i].Id)
			}

		} else {
			dataRow = Black{
				Creator:     creator,
				Status:      Status_Init,
				Ip:          val,
				Desc:        desc,
				Create_time: timeStr,
				Update_time: timeStr,
				Expire_time: expireTimeFormat,
			}
			_, err := o.Insert(&dataRow)
			if err != nil {
				o.Rollback()
				check.Err(err)
			}
		}
	}
	err = o.Commit()
	check.Err(err)
	return true
}

/*
	标记删除
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
func Delete(id int64) int64 {
	params := orm.Params{
		"status": Status_Del,
	}
	return Update(&params, id)
	//num, err := o.QueryTable(TableName).Filter("id", id).Update()
	//fmt.Printf("Affected Num: %s, %s", num, err)
	//return num
}

/*
	标记删除
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
func DeleteByIps(ips []string) bool {
	params := orm.Params{
		"status": Status_Del,
	}
	return UpdateByIps(&params, ips)
	//num, err := o.QueryTable(TableName).Filter("id", id).Update()
	//fmt.Printf("Affected Num: %s, %s", num, err)
	//return num
}

/*
	根据ip更新
*/
func Update(data *orm.Params, id int64) int64 {
	o := orm.NewOrm()
	num, err := o.QueryTable(TableName).Filter("id", id).Update(*data)
	check.Err(err)
	return num
}

/*
   根据ip 更新
*/
func UpdateByIps(data *orm.Params, ips []string) bool {
	o := orm.NewOrm()
	err := o.Begin()
	for _, val := range ips {
		_, err := o.QueryTable(TableName).Filter("ip", val).Update(*data)
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

/*
	查找最小的一一行数据
*/
func GetMinTime() Black {
	o := orm.NewOrm()
	//通过默认值区分是否获得了数组
	row := Black{
		Id: -1,
	}
	//o.Raw("SELECT * FROM  " + TableName + " t1 , (select id, min(expire_time) from " + TableName + " where `expire_time` > 0 && `status` != " + fmt.Sprintf("%d", Status_Del) + ") t2 WHERE t1.id = t2.id").QueryRow(&row)

	sql := "SELECT t1.* FROM  " + TableName + " t1 inner join (select id from " + TableName + " where `expire_time` > 0 && `status` = " + fmt.Sprintf("%d", Status_Online) + " order by expire_time limit 1) t2 on t1.id = t2.id"
	o.Raw(sql).QueryRow(&row)
	return row
}

//只有一列,online状态的只有一列
func GetByIp(ip string) Black {
	o := orm.NewOrm()
	tableName := GetName()
	var data = Black{
		Id: -1,
	}
	o.QueryTable(tableName).Filter("ip", ip).Filter("status", Status_Online).One(&data)
	return data
}
