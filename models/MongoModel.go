package models

import (
	"time"

	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var MongoDb *mgo.Database

type Adt struct {
	Gameid   string
	App_type string
	Channel  int
	Imei     string
}

var session *mgo.Session

func GetMongoSession() *mgo.Session {

	if session == nil {
		var err error
		var con_str string
		if MongodbConf.User == "" && MongodbConf.PassWord == "" {
			con_str = fmt.Sprintf("%s:%s", MongodbConf.Host, MongodbConf.Port)
		} else {
			con_str = fmt.Sprintf("%s:%s@%s:%s", MongodbConf.User, MongodbConf.PassWord, MongodbConf.Host, MongodbConf.Port)
		}

		session, err = mgo.Dial(con_str)
		mgo.DialWithTimeout(con_str, time.Second*60)

		if err != nil {
			fmt.Println("mongo connect err!!")
		}
	}

	return session
}

//create the data base
func createDbSession() {
	var err error
	session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{MongodbConf.Host},
		Username: MongodbConf.User,
		Password: MongodbConf.PassWord,
		Timeout:  60 * time.Second,
	})

	if err != nil {

	}
}

func MongoTest() {
	results := Adt{}
	mon := GetMongoSession().Copy()
	defer mon.Close()

	db := mon.DB("channel").C("test_channel")
	if db != nil {
		err := db.Find(bson.M{"muid": "b4c2a07a94bfd56651dd89c5d92664f8"}).One(&results)
		if err != nil {
		}
		fmt.Println(results)
	} else {
		fmt.Println("Mongo connect error!")
	}

}
