// models project doc.go

/*
models document
*/
package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//mysql config

var MysqlConf struct {
	Host     string
	User     string
	PassWord string
	Port     string
	DataBase string
}

//redis config
var RedisConf struct {
	Host     string
	User     string
	PassWord string
	Port     string
}

//mongodb config
var MongodbConf struct {
	Host     string
	User     string
	PassWord string
	Port     string
	DataBase string
}

type Mysqldb struct {
	DB *sql.DB
}

var Mydb Mysqldb

func init() {
	initDatabase() //初始化配置信息
	initMysql()    //初始化mysql
	initMongo()    //初始化mongo
}

/*
获取程序运行路径
*/
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("This is an error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//初始化mysql
func initDatabase() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("数据库初始化")
			return
		}
	}()
	path := getCurrentDirectory()
	fmt.Println("配置初始化")
	//获取当前文件的配置
	dat, err := ioutil.ReadFile(path + "/config/config.json")

	if err != nil {
		fmt.Println("为找到配置文件")
		return
	}

	dt := DataChange(string(dat))

	//mysql
	dta := dt["mysql"]
	dap := DataChange(JsonEncodeString(dta))
	MysqlConf.Host = dap["host"].(string)
	MysqlConf.PassWord = dap["password"].(string)
	MysqlConf.Port = dap["port"].(string)
	MysqlConf.User = dap["user"].(string)
	MysqlConf.DataBase = dap["database"].(string)

	//redis
	redisdt := dt["redis"]
	redisdap := DataChange(JsonEncodeString(redisdt))
	RedisConf.Host = redisdap["host"].(string)
	RedisConf.PassWord = redisdap["password"].(string)
	RedisConf.Port = redisdap["port"].(string)
	RedisConf.User = redisdap["user"].(string)

	//mongodb
	mongodt := dt["mongodb"]
	mongodap := DataChange(JsonEncodeString(mongodt))
	MongodbConf.Host = mongodap["host"].(string)
	MongodbConf.PassWord = mongodap["password"].(string)
	MongodbConf.Port = mongodap["port"].(string)
	MongodbConf.User = mongodap["user"].(string)
	MongodbConf.DataBase = mongodap["database"].(string)

}

//数据格式转化的操作

func DataChange(data string) map[string]interface{} {
	var dat map[string]interface{}
	json.Unmarshal([]byte(data), &dat)
	return dat
}

// 结构转换成json对象
func JsonEncodeString(data interface{}) string {
	back, err := json.Marshal(data)
	if err != nil {
		return "encode error"
	}
	return string(back)
}

//map的类型转换成！

func MaptoJson(data map[string]interface{}) string {
	configJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return ""
	}
	return string(configJSON) //返回格式化后的字符串的内容0
}

//mysql  初始化操作连接

func initMysql() {
	var err error
	o := Mysqldb{}
	cont := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", MysqlConf.User, MysqlConf.PassWord, MysqlConf.Host, MysqlConf.Port, MysqlConf.DataBase)

	o.DB, err = sql.Open("mysql", cont)
	o.DB.SetMaxIdleConns(2000)
	o.DB.SetMaxOpenConns(1000) //设置请求的连接池
	if err != nil {
		fmt.Errorf("open oracle database failed.", err)
	}
	Mydb = o

}

//初始化
func initRedis() {

}

var MongoDb *mgo.Database

type Adt struct {
	Gameid   string
	App_type string
	Channel  int
	Imei     string
}

//初始化mongodb
func initMongo() {
	//con_str := fmt.Sprintf("%s:%s@%s:%s", MongodbConf.User, MongodbConf.PassWord, MongodbConf.Host, MongodbConf.Port)
	session, err := mgo.Dial("106.75.146.174:4077")
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	result := Adt{}
	session.SetMode(mgo.Monotonic, true)
	MongoDb = session.DB(MongodbConf.DataBase)
	err = MongoDb.C("test_channel").Find(bson.M{"muid": "b4c2a07a94bfd56651dd89c5d92664f8"}).One(&result)

	if err != nil {
		fmt.Println("There have some error!")
	}

	fmt.Println(result)
}

//	创建数据库的查询

func (this *Mysqldb) Query(sql_q string, args ...interface{}) ([]map[string]string, bool) {
	rows, err := this.DB.Query(sql_q, args...)
	if err != nil {
		//return "", false
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
