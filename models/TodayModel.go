package models

import (
	//"dsp/dat/jrtt"
	"fmt"
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

	db := GetMysqlDb()
	sql := fmt.Sprintf("select * from tf_plan where is_off = 0 and address regexp %q", city)
	fmt.Println(sql)
	dat, _ := db.Query(sql)
	fmt.Println(dat)

}

//input Data

func TTinputData(bid *BidRequest) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("This is the best you config!")
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
