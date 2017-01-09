package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	//继承封装sql.Db，方便二次开发
	sql.DB
	conn *sql.DB
}
