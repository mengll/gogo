package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Mysqldb struct {
	DB *sql.DB
}

var Mydb Mysqldb

func GetMysqlDb() Mysqldb {
	defer func() {
		if conerr := recover(); conerr != nil {
			fmt.Println("mysql connect error!")
		}
	}()
	if Mydb.DB == nil {
		var err error
		cont := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", MysqlConf.User, MysqlConf.PassWord, MysqlConf.Host, MysqlConf.Port, MysqlConf.DataBase)

		Mydb.DB, err = sql.Open("mysql", cont)
		Mydb.DB.SetMaxIdleConns(2000)
		Mydb.DB.SetMaxOpenConns(1000)                //设置请求的连接池
		Mydb.DB.SetConnMaxLifetime(60 * time.Second) // set the life time

		if err != nil {
			fmt.Errorf("open mysql database failed.", err)
		}
	}
	return Mydb
}

//	创建数据库的查询

func (this *Mysqldb) Query(sql_q string, args ...interface{}) ([]map[string]string, bool) {
	connect := Mydb.DB.Ping()
	if connect != nil {
		//panic("mysql connect error")
		fmt.Errorf("mysql connect %s", "can't connect")
		fmt.Println("mysql connect error")
		return []map[string]string{}, false
	}

	rows, err := this.DB.Query(sql_q, args...)
	if err != nil {
		return []map[string]string{}, false
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		//return "", false
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var dat []map[string]string = []map[string]string{}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			//			return "", false
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		mpp := make(map[string]string)
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			mpp[columns[i]] = string(value)
		}
		dat = append(dat, mpp)

	}
	if err = rows.Err(); err != nil {
		//		return "", false
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return dat, true //返回当前查询的结果 ，当前的查询的状态
}

//执行写入操作
func (this *Mysqldb) Insert(sql_q string) interface{} {
	rows, err := this.DB.Query(sql_q)
	if err != nil {
		//return "", false
	}
	return rows
}

func initTest() {
	mydb := GetMysqlDb()
	sql := fmt.Sprintf("select rid,pid from ucusers where uid = %q or mobile =%q", "18827092404", "18827092404")
	dat, err := mydb.Query(sql)
	if err {
		fmt.Println("This Connnect is error!")
	}
	//println the Data
	fmt.Println(dat[0]["rid"])

}
