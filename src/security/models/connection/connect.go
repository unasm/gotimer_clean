/*
	connection	名字取的不太合理，用于处理各个model 的基础方法
*/
package connection

import (
	"database/sql"
	"security/conf"
	"security/lib/check"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*
 * 得到一个数据库的连接
	@params		string	name	数据库的名称
	@params		map		config	数据库的配置，规则包含 key 为 username,host, password,dbname,port, dbname,port默认为3306

*/
func GetLink(name string, config map[string]string) *sql.DB {
	if _, ok := config["port"]; !ok {
		//默认port 3306
		config["port"] = "3306"
	}
	db, err := sql.Open("mysql", config["username"]+":"+
		config["password"]+"@tcp("+config["host"]+":"+config["port"]+")/"+config["dbname"])
	check.Err(err)
	return db
}

func GetDsn(name string, config map[string]string) string {
	if _, ok := config["port"]; !ok {
		//默认port 3306
		config["port"] = "3306"
	}
	dsn := config["username"] + ":" +
		config["password"] + "@tcp(" + config["host"] + ":" + config["port"] + ")/" + config["dbname"] + "?charset=utf8"
	if value, ok := config["dbTimeout"]; ok {
		//默认port 3306
		//config["port"] = "3306"
		dsn = dsn + "&timeout=" + value
	}
	return dsn
}

/*
	初始化数据库连接
*/
func Init() {
	Config := conf.New()
	dsn := GetDsn("default", Config.Dict)
	orm.RegisterDataBase("default", "mysql", dsn, 30)

	//orm.RegisterDataBase("default", "mysql", "root:@tcp(100.73.13.11:3306)/security?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
	//orm.RunSyncdb("default", false, true)
}
