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

//func TTquery(bid *jrtt.BidRequest) {
//	db := GetMysqlDb()

//	data, err := db.Query("")
//	if err {
//	}

//	for k, v := range data {
//		fmt.Println(k, v)
//	}

//}

type TTPlan struct {
}

func TTinit() {
	db := GetMysqlDb()
	sql := "select * from tf_plan where is_off = 0"
	dat, err := db.Query(sql)
	fmt.Println(dat)
	fmt.Println("\r\n------")
	fmt.Println(err)

}

//input Data

//func TTinputData(bid *jrtt.BidRequest) {
//	dat := TTnodat{}
//	dat.data = bid.String()
//	MongoDb.C("bid_request_dat").Insert(dat)
//}
