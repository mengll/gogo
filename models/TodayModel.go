package models

import (
	//"dsp/dat/jrtt"
	"fmt"
	"reflect"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

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

//select data

//{
//    "requestId":"2017042611360317201704100613879E",
//    "apiVersion":"2.1",
//    "adslots":[
//        {
//            "id":"45549b62016140b0",
//            "banner":[
//                {
//                    "width":580,
//                    "height":240,
//                    "pos":"FEED",
//                    "sequence":"13"
//                }
//            ],
//            "adType":[
//                "TOUTIAO_FEED_APP_LARGE",
//                "TOUTIAO_FEED_LP_LARGE",
//                "TOUTIAO_FEED_APP_SMALL",
//                "TOUTIAO_FEED_LP_SMALL",
//                "TOUTIAO_FEED_LP_GROUP"
//            ],
//            "bidFloor":664,
//            "keywords":[

//            ]
//        }
//    ],
//    "app":{
//        "id":"13",
//        "name":"news_article",
//        "ver":"611"
//    },
//    "device":{
//        "ip":"113.57.183.196",
//        "geo":{
//            "lat":30.50482,
//            "lon":114.34165,
//            "city":"武汉"
//        },
//        "deviceId":"861451039768687",
//        "make":"unknown",
//        "model":"HUAWEI NXT-AL10",
//        "os":"android",
//        "osv":"7.0",
//        "connectionType":"NT_4G",
//        "deviceType":"PHONE",
//        "androidId":"51e64b24db09ef25"
//    },
//    "user":{
//        "id":"41341711578",
//        "yob":"31",
//        "gender":"MALE",
//        "data":[
//        ]
//    }
//}

func TTquery(bid *BidRequest) {
	//db := GetMysqlDb()

	//	data, err := db.Query("")
	//	if err {
	//	}

	//	for k, v := range data {
	//		fmt.Println(k, v)
	//	}
	fmt.Println("--<<><><>")
	ads := bid.GetAdslots()
	fmt.Println(ads)
	fmt.Println("<<<--->><")

	//get the device
	device := bid.GetDevice()
	device_geo := device.GetGeo()
	city := device_geo.GetCity()

	fmt.Println(city)
	fmt.Println(bid.Device.ConnectionType)

	db := GetMysqlDb()
	sql := fmt.Sprintf("select * from tf_plan where is_off = 0 and address regexp %q", city)
	fmt.Println(sql)
	dat, _ := db.Query(sql)
	fmt.Println(dat)

	//get the plan dat
	GetPlan(bid)
	getCastmoney()

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

//get
func getCastmoney() {
	mongo := GetMongoSession().Copy()
	//	var totalMoney int //the cast is fen
	defer mongo.Close()
	//	fmt.Println(planid)
	//	fmt.Println(totalMoney)
	//mongo.DB("channel").C("bid_request_dat").Insert(&indat)Z
	//	diter := mongo.DB("channel").C("click_Notify").Find({})

	//fmt.Println(diter)
	fmt.Println("=======><<><><<<>><><>---<<<><><")
	//	rea := Cnotify{} //the cast total money
	//	for diter.Next(&rea) {
	//	}
	//	money <- totalMoney

	job := &mgo.MapReduce{
		Map:    "function() { emit(this.did,this.bidprce) }",
		Reduce: "function(key, values) { return Array.sum(values) }",
	}
	var result []struct {
		Id    int "_id"
		Value int
	}

	_, err := mongo.DB("channel").C("win_Notify").Find(bson.M{"adid": '5'}).MapReduce(job, &result)

	if err != nil {
		fmt.Println("There had some error")
	}

	for _, item := range result {
		fmt.Println(item.Value)
	}

	fmt.Println("=-=-=-=-=--=-=-=-=-=-=-=-=-=0==-=--=-=-=-=-")
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

	sql += "status = 0"

	db := GetMysqlDb()

	db.Query(sql)

	fmt.Println(sql)

}

var USER_GENDER map[string]int = map[string]int{"UNKNOWN": 3, "FEMALE": 1, "MALE": 2}

var NT_ENUM map[string]int = map[string]int{"Honeycomb": 1, "WIFI": 2, "UNKNOWN": 3, "NT_2G": 4, "NT_4G": 5} //网络编辑

var Weeknum map[string]int = map[string]int{"Monday": 0, "Tuesday": 1, "Wednesday": 2, "Thursday": 3, "Friday": 4, "Saturday": 5, "Sunday": 6}
