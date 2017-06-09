package models

import (
	//"dsp/dat/jrtt"
	"encoding/json"
	"fmt"
	"math/rand"

	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
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

//保存计划的数组
var Planids chan int

//get the access plan

type BBQ struct {
	Bid *BidRequest
}

func GetPlan(bid *BidRequest) {
	if bid == nil {
		fmt.Println("there is empty for that !")
		return
	}
	selectMap := make(map[string]interface{})
	lk := bid.GetDevice().Geo.City
	fmt.Println("city", lk)
	if lk != "" {
		selectMap["address"] = lk
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
	fmt.Println("select_plan_sql=>", sql)
	db := GetMysqlDb()
	dat, merr := db.Query(sql)
	if merr == false {
		fmt.Println(merr)
	}
	//fmt.Println(dat)
	//get limit money
	//moneyLimit := make(chan int, 1)
	Planids = make(chan int, len(dat))
	pans := []string{}

	//数据长度
	fmt.Println("data length=>", len(dat))

	for _, v := range dat { //遍历当前的广告计划，查询是否满足当前的需求

		pans = append(pans, v["id"])
		pid, _ := strconv.Atoi(v["id"])
		Planids <- pid
		go CheckPlan(v, bid) //检测当前广告计划是否满足条件，满足条件 ，当满足条件的时候添加当前的计划到任务数组中，查询满足条件的广告创意 查询到
	}
	//plan string
	pids := []int{}
	for i := 0; i < len(dat); i++ {
		//	fmt.Println("当前的计划ID：", <-Planids)
		pids = append(pids, <-Planids)
	}
	//ads_id := []string{"10", "11"}
	bid.Chuangyi(pans)

	//查询当前的广告创意广告创意的ID的实现
	//out put the data
	//fmt.Println(sql)
	bid.BbqRequest()
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

func CheckPlan(dat map[string]string, bid *BidRequest) {
	//检查显示查询创意
	//创建通道传递当前的
	padt := make(chan bool, 1)
	//defer close(padt)
	//1 检查广告组是否开启 判断当前的是否超过了组的限额
	dsp_group := dat["group_id"]
	go groupLimit(dsp_group, padt) //只判断是否满足条件
	//fmt.Println("this is dsp_group ", dsp_group) //当前的组
	//2 检查投放的时间，是否存在限制
	fmt.Println("时间控制的测试")
	go timeLimit(padt, dat) //检查是否满足投放的需求

	//3广告计划的预算的控制方式
	ptype, _ := strconv.Atoi(dat["account_type"]) //投放的方式0 日预算 1 总预算
	//限制的金额
	go getTimemoney(dat, ptype, padt)
	//3，检查限额的控
	//4，检查关键词的投放的控制
	select {
	case <-time.After(time.Second * 3):
		fmt.Println("This is the file you can see in 90") //投放刚超时
		break
	case <-padt:
		Planids <- 1 //
		fmt.Println("This is ok the config!")
	}
}

//获取组的限额和时候否超出当前的限额 最后使用通道的方式传递数据

func groupLimit(groupid string, bb chan bool) {
	group_id, _ := strconv.Atoi(groupid)
	groupSql := fmt.Sprintf("select g.account_price,p.id ,g.account_method from tf_ads_group as g inner join tf_plan as p on g.id = p.group_id where g.id = %d and g.status =0 and g.is_off = 0", group_id)
	fmt.Println("查询相关额度组", groupSql)
	db := GetMysqlDb()
	dat, _ := db.Query(groupSql)
	//返回当前的文件
	if len(dat) > 0 {
		//返回当前广告组的限额
		group_money := dat[0]["account_price"]
		fmt.Println(group_money)
		fmt.Println("account_method=>", dat[0]["account_method"])
		//检查当前是否开启了限额
		if dat[0]["account_method"] == "1" {
			//广告组现在的总金额
			group_now_money := 0
			for _, v := range dat {
				//广告计划的ID
				m := getQdat(v["id"])
				group_now_money += m
			}
			fmt.Println("This is group :", group_now_money)
			glimitmoney, _ := strconv.Atoi(group_money)
			if group_now_money > glimitmoney {
				fmt.Println("超过了组限额")
				bb <- false //超过了当前的限额
			} else {
				bb <- true
			}
		}
		bb <- true
		//return group_now_money

	} else {
		//当前没有匹配到相关的广告组信息 meiyou uanggaiio
		fmt.Println("dang")
		//return -2 //当出现-1的时候表示当前的是不成立的
	}

}

//time limit
func timeLimit(bbq chan bool, dat map[string]string) {

	//当前投放时间
	if dat["tf_time"] == "0" && dat["tf_style"] == "0" {
		//完全的满足不用考虑时间
	}

	//投放时间限制
	if dat["tf_time"] == "0" && dat["tf_style"] == "1" {
		//判断办小时的时间
		JK := Stlimit(dat["tf_range"])
		if JK == false {
			//时间检测未能通过
		}
	}
	nowt := time.Now()
	timestamp := time.Date(nowt.Year(), nowt.Month(), nowt.Day(), 0, 0, 0, 0, time.Local) //当日的时间戳

	st_unix := strconv.FormatInt(timestamp.Unix(), 10)
	start_unix_time, _ := strconv.Atoi(st_unix)
	//时间对象的控制S
	if dat["tf_time"] == "1" && dat["tf_style"] == "0" {
		end_time, _ := strconv.Atoi(dat["end_time"])
		st_time, _ := strconv.Atoi(dat["start_time"])

		if start_unix_time > end_time || start_unix_time < st_time {
			//出错了少年 不在投放的时间段呢
		}
		//判断开始时间，结束时间
	}

	//时间段控制的处理
	if dat["tf_time"] == "1" && dat["tf_style"] == "1" {
		//时间段的判断办小时
		end_time, _ := strconv.Atoi(dat["end_time"])
		st_time, _ := strconv.Atoi(dat["start_time"])
		if start_unix_time > end_time || start_unix_time < st_time {
			//出错了少年 不在投放的时间段呢
		}
		JK := Stlimit(dat["tf_range"])
		//时间点的判断
		if JK == false {
			//时间检测未能通过
			fmt.Println("This is error!")
		}
	}
}

//当前的分钟的级别的投放的控制
func Stlimit(selecTime string) bool {
	var jj [][]string
	json.Unmarshal([]byte(selecTime), &jj)
	nowt := time.Now()
	dl := nowt.Weekday().String()
	daytime := Weeknum[dl]
	htime := nowt.Hour()                             //获取当前时
	mintime := nowt.Minute()                         //获取当前的分钟
	skip := mintime / 30                             //当前的时偏移量
	hanftindex := htime*2 + skip                     //半小时索引
	fmt.Println("是否处于投放时间", jj[daytime][hanftindex]) //获取当前是否投放广告 tf_style 投放的方式全天
	//当前的如果不满组条件的时候创建一个新的文件
	//	tp, _ := strconv.Atoi( jj[daytime][hanftindex]))
	//检查当前的投放是不是全天的投放的模式
	if jj[daytime][hanftindex] == "0" {
		return false
	} else {
		return true
	}
}

//查询当前订单的信息
type BackData struct {
	Adid    string `bson:"adid"`
	Bidprce int    `bson:"bidprce"`
}

//获取当前传递的是时间限制 0 表示按天的限额 1总的限额 计划的限额的
func getTimemoney(dat map[string]string, tp int, bbq chan bool) {
	plstr := dat["id"]
	var LimitMoney int
	//plstr := strconv.Itoa(planid)
	if tp == 0 {
		nowt := time.Now()
		timestamp := time.Date(nowt.Year(), nowt.Month(), nowt.Day(), 0, 0, 0, 0, time.Local) //当日的时间戳
		start_unix_time := timestamp.Unix()
		sd := strconv.FormatInt(start_unix_time, 10)
		mongdb := GetMongoSession().Copy()
		defer mongdb.Close()
		///, "timestamp": bson.M{"$lt": strconv.FormatInt(start_unix_time, 10)}
		result := BackData{}
		itler := mongdb.DB("channel").C("win_Notify").Find(bson.M{"adid": plstr, "timestamp": bson.M{"$gte": sd}}).Iter() // 当前投放的时间的偷偷放的金额的 今天 大于今日的零

		for itler.Next(&result) {
			LimitMoney += result.Bidprce
		}
		fmt.Println("Total 总金额", LimitMoney)

	} else {
		LimitMoney = getQdat(plstr)
		fmt.Println(LimitMoney)
	}

	limitmoneyOne := dat["account_price"]
	money, cerr := strconv.ParseFloat(limitmoneyOne, 10) //计划的限额
	if cerr != nil {
		fmt.Println("cha")
	}
	account_money := money * 100
	fmt.Println(account_money)
	t := strconv.FormatFloat(account_money, 'f', 0, 64)
	tlimi, _ := strconv.Atoi(t)

	fmt.Println("moneyLimit=->", LimitMoney)
	if tlimi < LimitMoney {
		//超过了
		fmt.Println("超过===》那就是失败")
		bbq <- true
	}
}

//查满足条件的广告创意

func (bid *BidRequest) Chuangyi(pids []string) {
	fmt.Println("This is chuangyi func !")
	pidst := strings.Join(pids, ",")
	fmt.Println("广告创意的组合====》", pidst)
	sql := fmt.Sprintf("select s.* from tf_plan as p inner join tf_ads as s on p.id = s.plan_id where p.id in(%s) and p.is_off = 0 and s.is_off = 0 ", pidst)
	fmt.Println("创意sql =", sql)
	mydb := GetMysqlDb()
	dat, isdat := mydb.Query(sql)

	if isdat {
		//当前的时间限制的
		fmt.Println("查询出错了。。。", dat)
	}
	//遍历当前的广告创意
	for _, dv := range dat {
		fmt.Println("用户中心==》", dv)
		//获取当前的投放的位置
		ad_type := dv["ad_type"]

		type_s := strings.Split(ad_type, ",")

	}

	//请求的广告创意
	adsNum := bid.GetAdslots()
	for _, ads := range adsNum {
		cc := []int32{}
		for _, vm := range ads.AdType {
			kv := vm.String()
			fmt.Println(vm)
			knum := AdType_value[kv]
			//	vl := AdType_value[v]
			fmt.Println(knum) //当前的广告位置
			cc = append(cc, knum)
		}
		//获取当前广告位的大小
		//获取当前的他又放的
		fmt.Println(cc)
	}
	//查询满足条件的创意信息

	fmt.Println("This is a gob ")
	fmt.Println(sql)
}

func (bid *BidRequest) BbqRequest() []byte {
	var domain string = "http://jrtt.qcwanwan.com"
	var param string = "notify?user_id={user_id}&request_id={request_id}&adid={adid}&bid_price={bid_price}&ip={ip}&timestamp={timestamp}&did={did}"

	//响应的对象的数据
	rback := rand.New(rand.NewSource(time.Now().UnixNano()))

	adslots := bid.GetAdslots()[0] //获取最上曾数据
	biddata := &Bid{}
	biddata.Id = strconv.Itoa(rback.Intn(20)) //生成唯一的商品的信息
	biddata.Adid = uint64(rback.Intn(18))     //参与竞价的广告ID
	biddata.Price = adslots.BidFloor + 1
	biddata.AdslotId = adslots.Id
	biddata.Cid = strconv.Itoa(rback.Intn(18)) //扩展的广告ID

	modelste := &MaterialMeta{} //广告素材对象
	modelste.AdType = AdType_TOUTIAO_FEED_LP_GROUP
	modelste.Nurl = fmt.Sprintf("%s/%s/%s", domain, "win", param)
	modelste.Title = "白野猪爆了个装备.换了小半个月工资NB"
	modelste.Source = "传奇无双"

	//banner的图片信息显示
	imgbanner := &MaterialMeta_ImageMeta{}
	imgbanner.Width = 228
	imgbanner.Height = 150
	imgbanner.Url = "http://jrtt.qcwanwan.com/1.jpg"
	imgbanner.Urls = []string{"http://jrtt.qcwanwan.com/1.jpg",
		"http://jrtt.qcwanwan.com/2.jpg",
		"http://jrtt.qcwanwan.com/3.jpg"}
	modelste.ImageBanner = imgbanner

	//设置当前的操作流程
	dsp_external := &MaterialMeta_ExternalMeta{}
	dsp_external.Url = "http://m.anfeng.cn/cqws_bbk-ios/12/"
	modelste.External = dsp_external
	biddata.Creative = modelste

	modelste.ShowUrl = []string{fmt.Sprintf("%s/%s/%s", domain, "show", param)}
	modelste.ClickUrl = []string{fmt.Sprintf("%s/%s/%s", domain, "click", param)}

	dsp := &SeatBid{}
	dsp.Ads = []*Bid{biddata}

	res := &BidResponse{}
	res.Seatbids = []*SeatBid{dsp}
	res.RequestId = bid.RequestId

	data, err := proto.Marshal(res)

	newTest := &BidResponse{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		fmt.Println("----->>>")
	}
	return data
}
