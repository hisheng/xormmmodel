// Package xorm
// @author 张海生<zhanghaisheng@qimao.com>
// @dateTime   : 2021/11/8 4:30 下午
package xorm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var xormTable string

func InitStruct(xormDsn string, table string) {
	xormTable = table
	//用系统orm，这样可以兼容以后的gorm等
	mysqlDb, err := CreateMysqlDb(xormDsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mysqlDb.Close()
	if len(xormTable) > 0 {
		//生成单表
		initTableStruct(mysqlDb)
	} else {
		//生成整个数据库
		tables, err := mysqlDb.Query("SELECT table_name FROM information_schema.TABLES WHERE table_schema=DATABASE () AND table_type='BASE TABLE'; ")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer tables.Close()
		for tables.Next() {
			err = tables.Scan(&xormTable)
			if err != nil {
				fmt.Println(err)
				continue
			}
			initTableStruct(mysqlDb)
		}

	}

}

func initTableStruct(mysqlDb *sql.DB) {
	columns, err := mysqlDb.Query("SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT,COLUMN_TYPE ,COLUMN_DEFAULT FROM information_schema.COLUMNS WHERE table_schema=DATABASE() AND table_name=?;", xormTable)
	if err != nil {
		fmt.Println(err)
	}
	defer columns.Close()

	row := mysqlDb.QueryRow("SELECT k.column_name FROM information_schema.table_constraints t JOIN information_schema.key_column_usage k USING(constraint_name,table_schema,table_name) WHERE t.constraint_type='PRIMARY KEY' AND t.table_schema= DATABASE() AND t.table_name= ?;", xormTable)
	type PK struct {
		ColumnName string
	}
	var pk PK
	err = row.Scan(&pk.ColumnName)
	fmt.Println(pk, err)

	structStrArr := make([]string, 0, 0)
	for columns.Next() {
		columnName := ""
		dataType := ""
		isNullable := ""
		tableName := ""
		columnComment := ""
		columnType := ""
		var defaultValue interface{}
		err = columns.Scan(&columnName, &dataType, &isNullable, &tableName, &columnComment, &columnType, &defaultValue)
		if err != nil {
			fmt.Println(err)
		}
		null := "not null"
		if isNullable == "YES" {
			null = "null"
		}
		comment := ""
		if len(columnComment) > 0 {
			comment = "comment('"
			comment += columnComment
			comment += "')"
		}
		defaultValueString := ""
		if defaultValue != nil {
			defaultValueString = fmt.Sprintf("default('%s')", defaultValue)
		}

		_type, ok := typeForMysqlToGo[dataType]
		if !ok {
			_type = "[]byte"
		}

		//主键
		pkString := ""
		if columnName == pk.ColumnName {
			pkString = "pk"
		}
		rowXorm := fmt.Sprintf("	%s %s `json:\"%s\" xorm:\"%s %s %s %s %s %s\"` \n", upperCamelCase(columnName), _type, columnName, "'"+columnName+"'", columnType, null, pkString, defaultValueString, comment)

		structStrArr = append(structStrArr, rowXorm)
	}
	saveToFile(xormTable, structStrArr)
}

//map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int64",
	"integer":            "int64",
	"tinyint":            "int64",
	"smallint":           "int64",
	"mediumint":          "int64",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
	"json":               "string",
}
