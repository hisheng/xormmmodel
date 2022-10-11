package xorm

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB = nil

//CreateMysqlDb 创建db
func CreateMysqlDb(dsn string) (*sql.DB, error) {
	if Db != nil {
		return Db, nil
	}
	_db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil
	} else {
		return _db, nil
	}
}
