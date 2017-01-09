package rule

import (
	"security/lib/check"
	"security/models/connection"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const TableName = "rule"

//图形验证码阻断方式
const Way_Img = 1

//短信验证码阻断方式
const Way_Msg = 2

//直接拒绝
const Way_Deny = 3

//刚刚添加
const Status_Init = 1

//线上在用
const Status_Online = 3

//添加失败
const Status_Failed = 7

//删除
const Status_Del = 14

type Rule struct {
	Id     int64
	Uri    string
	Dim_id int64
	User   string
	//	Name        string
	Times       int64
	Expire      int64
	Way         int64
	Ext         string
	Reason      string
	Status      int64
	Create_time string
	Update_time string
}

/*
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`uri` varchar(256) NOT NULL DEFAULT '' COMMENT 'URI',
`dim_id` int(11) NOT NULL DEFAULT '0',
`user` varchar(64) NOT NULL DEFAULT '',
`name` varchar(64) NOT NULL DEFAULT '' COMMENT '策略名称',
`total` int(11) NOT NULL DEFAULT '0' COMMENT '次数',
`expire` int(11) NOT NULL DEFAULT '0' COMMENT '超时周期, 分、小时、天',
`way` int(11) NOT NULL DEFAULT '0' COMMENT '阻断方式,1.图形验证码、2.短信验证码、3.直接返回',
`ext` varchar(512) NOT NULL DEFAULT '扩展字段、建议存放json',
`reason` varchar(512) NOT NULL DEFAULT '' COMMENT '创建、修改策略原因',
`status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0未授权1已生效2未生效',
`create_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
`update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
*/

func init() {
	connection.Init()
	orm.RegisterModel(new(Rule))
}

//获取数据库的表名
func GetName() string {
	return TableName
}

/*
	插入数据库数据
*/
func Add(data *Rule) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(data)
	check.Err(err)
	return id
}

/*
 *  更新数据库
    @param	data	Rule	要更新的数据内容
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
func Update(data *orm.Params, id int64) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable(TableName).Filter("id", id).Update(*data)
	check.Err(err)
	if num > 0 {
		return true
	}
	return false
}

/*
	标记删除
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
func Delete(id int64) bool {
	params := orm.Params{
		"status": Status_Del,
	}
	if Update(&params, id) == false {
		return false
	}
	return true
}

/*
	标记删除
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
func BatchDelete(ids []int64) bool {
	params := orm.Params{
		"status": Status_Del,
	}
	o := orm.NewOrm()
	err := o.Begin()
	for _, val := range ids {
		if Update(&params, val) == false {
			o.Rollback()
			return false
		}
	}
	err = o.Commit()
	check.Err(err)
	return true
}

/*
 *  更新数据库
    @param	data	Rule	要更新的数据内容
	@param	id		int64	主键的id
	@return		bool	更新是否成功
*/
/*
func Update(data Rule, id int64) bool {
	o := orm.NewOrm()
	where := Rule{Id: id}
	err := o.Read(&where)

	fmt.Println(orm.ErrNoRows)
	fmt.Println(orm.ErrMissPK)
	if err == orm.ErrNoRows {
		//fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		//fmt.Println("找不到主键")
	} else {
		data.Id = id
		fmt.Println(data)
		num, err := o.Update(&data)
		check.Err(err)
		fmt.Println(num)
		return true
	}
	return true
	//== nil
	//o.QueryTable(TableName).Filter("id", id).Update
}

*/
/**
* 从数据库里面 获取一条记录
 */
//func GetOne() Rule {
func GetOne() *Rule {
	o := orm.NewOrm()
	var data []*Rule
	qs := o.QueryTable(TableName)
	_, err := qs.Filter("id", 1).All(&data)
	check.Err(err)
	return data[0]
}
