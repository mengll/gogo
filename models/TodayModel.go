package models

import (
	//"dsp/dat/jrtt"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

var USER_GENDER map[string]int = map[string]int{"UNKNOWN": 3, "FEMALE": 1, "MALE": 2}

var NT_ENUM map[string]int = map[string]int{"Honeycomb": 1, "WIFI": 2, "UNKNOWN": 3, "NT_2G": 4, "NT_4G": 5} //网络编辑

var Weeknum map[string]int = map[string]int{"Monday": 0, "Tuesday": 1, "Wednesday": 2, "Thursday": 3, "Friday": 4, "Saturday": 5, "Sunday": 6}

type TTdata struct {
	ip       string
	lat      string
	lon      string
	city     string
	deviceId string
	os       string
	osv      string
	yod      string
	gender   string
}

type TTnodat struct {
	data string
}

const (
	TF_STYLE_ALL  = 0
	TF_STYLE_TIME = 1
)

func TTquery(bid *BidRequest) {

	//	fmt.Println("--<<><><>")
	//	ads := bid.GetAdslots()
	//	fmt.Println(ads)
	//	fmt.Println("<<<--->><")

	//	//get the device
	//	device := bid.GetDevice()
	//	device_geo := device.GetGeo()
	//	city := device_geo.GetCity()

	//	fmt.Println(city)
	//	fmt.Println(bid.Device.ConnectionType)

	//	db := GetMysqlDb()
	//	sql := fmt.Sprintf("select * from tf_plan where is_off = 0 and address regexp %q", city)
	//	fmt.Println(sql)
	//	dat, _ := db.Query(sql)
	//	fmt.Println(dat)

	//get the plan dat
	GetPlan(bid)
	//get the limit money

}

//input Data

func TTinputData(bid *BidRequest) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("This is the best you config!") // this will be write bidrequest log
		}
	}()

	indat := make(map[string]interface{})
	indat["deviceid"] = bid.Device.DeviceId
	indat["os"] = bid.Device.Os
	indat["osv"] = bid.Device.Osv
	indat["carrier"] = bid.Device.Carrier
	indat["devicetype"] = bid.Device.DeviceType
	indat["model"] = bid.Device.Model
	indat["ip"] = bid.Device.Ip
	indat["requestid"] = bid.RequestId
	indat["city"] = bid.Device.Geo.City
	indat["Country"] = bid.Device.Geo.Country
	indat["lat"] = bid.Device.Geo.Lat
	indat["lon"] = bid.Device.Geo.Lon
	indat["region"] = bid.Device.Geo.Region

	indat["gender"] = bid.User.Gender
	indat["keywords"] = bid.User.Keywords
	indat["yob"] = bid.User.Yob

	adslosts := bid.GetAdslots()
	adsdat := []string{}

	for _, v := range adslosts {
		adsdat = append(adsdat, v.String())
	}

	indat["adslots"] = adsdat
	mongo := GetMongoSession().Copy()
	defer mongo.Close()
	mongo.DB("channel").C("bid_request_dat").Insert(&indat)

}

type Cnotify struct {
	did   uint32
	money int
}

//get the access plan

func GetPlan(bid *BidRequest) {
	if bid == nil {
		fmt.Println("there is empty for that !")
		return
	}
	selectMap := make(map[string]interface{})
	lk := bid.GetUser().GetGeo()
	if lk != nil {
		selectMap["address"] = lk.String()
	}
	ugender := bid.User.GetGender()
	if ugender.String() != "" {
		selectMap["gender"] = USER_GENDER[bid.User.Gender.String()]
	}

	if bid.Device.Os != "" {
		selectMap["platform"] = bid.Device.Os
	}

	if bid.User.Yob != "" {
		selectMap["yob"] = bid.User.Yob
	}
	selectMap["connection_type"] = NT_ENUM["NT_4G"]
	if bid.Device.Carrier != "" {
		selectMap["operator"] = bid.Device.Carrier
	}

	sql := "select * from tf_plan where "
	for k, v := range selectMap {
		ty := reflect.TypeOf(v)
		nn := ty.Name()
		switch nn {
		case "string":
			sql += fmt.Sprintf("if(%s=-1,1,FIND_IN_SET(%q,%s)) and ", k, v, k)
		case "int":
			sql += fmt.Sprintf("if(%s=-1,1,FIND_IN_SET(%d,%s)) and ", k, v, k)
		}

	}
	jk := time.Now()

	sql += " start_time < " + strconv.FormatInt(jk.Unix(), 10) + " and "

	sql += "status = 0"
	fmt.Println(sql)
	db := GetMysqlDb()
	dat, merr := db.Query(sql)
	if merr == false {
		fmt.Println(merr)
	}
	fmt.Println(dat)

	//get limit money
	moneyLimit := make(chan int, 1)

	for _, v := range dat {
		fmt.Println(v["id"])
		CheckPlan(v)
	}

	//out put the data
	fmt.Println(<-moneyLimit)
	//fmt.Println(sql)

}

//获取当前的登录信息
func getQdat(planid string) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("there had some error")
		}
	}()
	mongo := GetMongoSession().Copy()
	defer mongo.Close()

	job := &mgo.MapReduce{
		Map:    "function() { emit(this.adid,this.bidprce) }",
		Reduce: "function(key, values) { return Array.sum(values) }",
	}

	var result []struct {
		Id    int `bson:"adid"` //广告的ID
		Value int
	}
	_, err := mongo.DB("channel").C("win_Notify").Find(bson.M{"adid": planid}).MapReduce(job, &result)

	if err != nil {
		fmt.Println("There had some error")
	}
	var total int
	for _, item := range result {
		total = item.Value
	}
	return total
}

//select order id by plan dat

func CheckPlan(dat map[string]string) {
	limitmoney := dat["account_price"]
	money, cerr := strconv.ParseFloat(limitmoney, 10)
	if cerr != nil {
		fmt.Println("cha")
	}
	account_money := money * 100
	fmt.Println(account_money)
	//check money

	//check time
	tf_style, _ := strconv.Atoi(dat["tf_style"])

	//check tf_time
	if tf_style != TF_STYLE_ALL {
		fmt.Println("tf_toufang all day")
	}
	//bu jiancha shijian
	fmt.Println(tf_style)

	//

}
