package models

import (
	"Tdsp/data/today"
	"fmt"
	"reflect"
	"strconv"
	//"sync"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Today struct {
	BID_REQUEST *today.BidRequest
	GET_PLANS   []map[string]interface{} //根据条件获取相关的广告计划
}

var (
	TODAY_USER_GENDER map[string]int = map[string]int{"UNKNOWN": 3, "FEMALE": 1, "MALE": 2}
	TODAY_NT_ENUM     map[string]int = map[string]int{"Honeycomb": 1, "WIFI": 2, "UNKNOWN": 3, "NT_2G": 4, "NT_4G": 5} //网络编辑
	TODAY_WEEK_NUM    map[string]int = map[string]int{"Monday": 0, "Tuesday": 1, "Wednesday": 2, "Thursday": 3, "Friday": 4, "Saturday": 5, "Sunday": 6}
)

const (
	TF_STYLE_ALL  = 0
	TF_STYLE_TIME = 1
)

//查询当前订单的信息
type BackData struct {
	Adid    string `bson:"adid"`
	Bidprce int    `bson:"bidprce"`
}

func (self *Today) Checkplan(dat map[string]string, planids chan string) {
	//bid := self.BID_REQUEST
	//检查显示查询创意
	//创建通道传递当前的
	padt := make(chan bool, 3)
	//defer close(padt)
	//1 检查广告组是否开启 判断当前的是否超过了组的限额
	dsp_group := dat["group_id"]
	go groupLimit(dsp_group, padt) //只判断是否满足条件

	//2 检查投放的时间，是否存在限制
	go timeLimit(padt, dat) //检查是否满足投放的需求

	//3广告计划的预算的控制方式
	ptype, _ := strconv.Atoi(dat["account_type"]) //投放的方式0 日预算 1 总预算
	//限制的金额
	go getTimemoney(dat, ptype, padt)

	//3，检查限额的控
	//CpmMoneyLimit(dat, ptype, padt)
	//4，检查关键词的投放的控制

	//遍历当前的条 如果当前的条件失败，那当前的计划是不能写入到
	var end_pv bool = true

	for i := 0; i < 3; i++ {
		pv := <-padt
		if pv == false {
			end_pv = false
		}
	}
	if end_pv == true {
		planids <- dat["id"]
	} else {
		planids <- ""
	}

}

//获取当前的广告计划

func (self *Today) GetPlans(bid *today.BidRequest) []byte {
	go self.SaveRquest()
	self.BID_REQUEST = bid

	if bid == nil {
		return []byte{}
	}

	selectMap := make(map[string]interface{})
	lk := *bid.GetDevice().Geo.City
	if lk != "" {
		selectMap["address"] = lk
	}
	ugender := bid.User.GetGender()
	if ugender.String() != "" {
		selectMap["gender"] = TODAY_USER_GENDER[bid.User.Gender.String()]
	}

	if *bid.Device.Os != "" {
		selectMap["app_type"] = *bid.Device.Os
	}
	yodstr := bid.User.GetYob()
	if yodstr != "" {
		selectMap["yob"] = yodstr
	}
	contype := bid.Device.ConnectionType.String()
	if contype != "" {
		selectMap["connection_type"] = TODAY_NT_ENUM[contype]
	}
	devicestr := bid.Device.GetCarrier()
	if devicestr != "" {
		selectMap["operator"] = devicestr
	}

	sql := "select * from tf_plan where "
	for k, v := range selectMap {
		ty := reflect.TypeOf(v)
		nn := ty.Name()
		if k != "app_type" {
			switch nn {
			case "string":
				sql += fmt.Sprintf("if(%s=-1,1,FIND_IN_SET(%q,%s)) and ", k, v, k)
			case "int":
				sql += fmt.Sprintf("if(%s=-1,1,FIND_IN_SET(%d,%s)) and ", k, v, k)
			}
		} else {
			sql += fmt.Sprintf("if(%s='',1,FIND_IN_SET(%q,app_type)) and ", k, v)
		}

	}

	jk := time.Now()
	sql += " start_time < " + strconv.FormatInt(jk.Unix(), 10) + " and "

	sql += "is_off = 0"

	fmt.Println("Getpaln", sql)

	//插叙当前满足条件的广告计划
	mysqldb := GetMysqlDb()

	fmt.Println("查询计划的sql", sql)

	dat, merr := mysqldb.Query(sql) //查询到满足条件的所有的广告计划
	if merr == false {
		//fmt.Println(merr)
	}

	//waitGroup := sync.WaitGroup
	//waitGroup.Add(len(dat)) //并发查询

	var planids chan string = make(chan string, len(dat))
	for _, v := range dat { //遍历当前的广告计划，查询是否满足当前的需求
		//		pid, _ := strconv.Atoi(v["id"])
		go self.Checkplan(v, planids) //检测当前广告计划是否满足条件，满足条件 ，当满足条件的时候添加当前的计划到任务数组中，查询满足条件的广告创意 查询到
	}

	pans := []string{}
	for i := 0; i < len(dat); i++ {
		pans = append(pans, <-planids)
	}

	//lks := self.Chuangyi(pans)
	lks := self.Backads(pans)
	return lks

}

//保存保存传递过来的数据
func (self *Today) SaveRquest() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("This is the best you config!") // this will be write bidrequest log
		}
	}()

	bid := self.BID_REQUEST
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
	indat["time"] = time.Now().Unix()

	adslosts := bid.GetAdslots()
	adsdat := []string{}

	for _, v := range adslosts {
		adsdat = append(adsdat, v.String())
	}

	indat["adslots"] = adsdat
	mongo := GetMongoSession().Copy()
	defer mongo.Close()
	mongo.DB(MongodbConf.DataBase).C("bid_request_dat").Insert(&indat)

}

//获取组的限额和时候否超出当前的限额 最后使用通道的方式传递数据

func groupLimit(groupid string, bb chan bool) {
	group_id, _ := strconv.Atoi(groupid)
	groupSql := fmt.Sprintf("select g.account_price,p.id ,g.account_method,p.account_style,g.name,g.id from tf_ads_group as g inner join tf_plan as p on g.id = p.group_id where g.id = %d and g.status =0 and g.is_off = 0", group_id)
	fmt.Println("groupsql", groupSql)
	db := GetMysqlDb()
	dat, _ := db.Query(groupSql)
	//返回当前的文件
	if len(dat) > 0 {
		//返回当前广告组的限额
		group_money := dat[0]["account_price"]
		//检查当前是否开启了限额
		if dat[0]["account_method"] == "1" {
			//广告组现在的总金额
			group_now_money := 0
			for _, v := range dat {
				//广告计划的ID
				m := getQdat(v["id"], v["account_style"])
				group_now_money += m
			}

			//glimitmoney, _ := strconv.Atoi(group_money)
			glimitmoney := SftoI(group_money)
			if group_now_money > glimitmoney {

				content := fmt.Sprintf("http://dspadmin.qcwan.com/index.php?d=mp&c=MessageAction&m=sendMsg&msg=当前的广告组:%s\n广告组ID：%s\n超过了当前的限额：%d 元 今日的投放结束", dat[0]["name"], dat[0]["id"], glimitmoney)
				go Httprequest(content, "GET", "")
				bb <- false //超过了当前的限额
			} else {
				bb <- true
			}
		}
		bb <- true
		//return group_now_money

	} else {
		bb <- false //当数据没有查询到的时候使用默认值
		//当前没有匹配到相关的广告组信息 meiyou uanggaiio
		fmt.Println("dang")
		//return -2 //当出现-1的时候表示当前的是不成立的
	}

}

//time limit
func timeLimit(bbq chan bool, dat map[string]string) {

	//当前投放时间
	if dat["tf_time"] == "0" && dat["tf_style"] == "0" {
		//fmt.Println("时间出问题")
		//完全的满足不用考虑时间
		bbq <- true
		return
	}

	//投放时间限制
	if dat["tf_time"] == "0" && dat["tf_style"] == "1" {
		//判断办小时的时间
		JK := Stlimit(dat["tf_range"])
		if JK == false {
			//时间检测未能通过
			bbq <- false
			return
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
			bbq <- false
			return
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
			bbq <- false
			return
		}
		JK := Stlimit(dat["tf_range"])
		//时间点的判断
		if JK == false {
			//时间检测未能通过
			bbq <- false
			return
		}
	}
	bbq <- true
}

//获取当前传递的是时间限制 0 表示按天的限额 1总的限额 计划的限额的
func getTimemoney(dat map[string]string, tp int, bbq chan bool) {

	plstr := dat["id"] //当前的计划的id

	account_type := dat["account_type"]
	var mongotable string = "win_Notify"
	switch account_type {
	case "CPC":
		mongotable = "click_Notify"
	case "CPM":
		mongotable = "win_Notify"

	}

	var LimitMoney int
	//plstr := strconv.Itoa(planid)
	if tp == 0 { //每天的限额

		nowt := time.Now()
		timestamp := time.Date(nowt.Year(), nowt.Month(), nowt.Day(), 0, 0, 0, 0, time.Local) //当日的时间戳
		start_unix_time := timestamp.Unix()
		sd := strconv.FormatInt(start_unix_time, 10)
		mongdb := GetMongoSession().Copy()
		defer mongdb.Close()
		///, "timestamp": bson.M{"$lt": strconv.FormatInt(start_unix_time, 10)}
		result := BackData{}
		itler := mongdb.DB(MongodbConf.DataBase).C(mongotable).Find(bson.M{"planid": plstr, "timestamp": bson.M{"$gte": sd}}).Iter() // 当前投放的时间的偷偷放的金额的 今天 大于今日的零

		for itler.Next(&result) {
			LimitMoney += result.Bidprce
		}

	} else {
		LimitMoney = getQdat(plstr, account_type) //传递的参数，是当前计划的ID
	}

	limitmoneyOne := dat["account_price"]
	money, cerr := strconv.ParseFloat(limitmoneyOne, 10) //计划的限额
	if cerr != nil {
		bbq <- false
		return
	}
	account_money := money * 100

	t := strconv.FormatFloat(account_money, 'f', 0, 64)
	tlimi, _ := strconv.Atoi(t)

	//当前的消费的金额加上这次消费的金额
	fmoney := SftoI(dat["first_price"])
	Usemoney := LimitMoney
	LimitMoney += fmoney
	fmt.Println("现在已经使用金额为", LimitMoney)
	fmt.Println("计划的限额为", tlimi)

	if tlimi < LimitMoney {
		//超过了
		//go SendMsgToDsp(content)

		bbq <- false
		return
	}
	content := fmt.Sprintf("http://dspadmin.qcwan.com/index.php?d=mp&c=MessageAction&m=sendMsg&msg=当前的广告计划:%s-计划ID：%s- 当前的限额：%s 元 今日消费了%d", dat["name"], dat["id"], dat["account_price"], Usemoney/100.0)
	go Httprequest(content, "GET", "")

	bbq <- true
}

//获取当前的登录信息
func getQdat(planid, tb string) int {

	var mongotable string = "win_Notify"
	switch tb {
	case "CPC":
		mongotable = "click_Notify"
	case "CPM":
		mongotable = "show_Notify"

	}

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
	_, err := mongo.DB(MongodbConf.DataBase).C(mongotable).Find(bson.M{"planid": planid}).MapReduce(job, &result)

	if err != nil {
		fmt.Println("There had some error")
	}
	var total int
	for _, item := range result {
		total = item.Value
	}
	return total
}

//当前的分钟的级别的投放的控制
func Stlimit(selecTime string) bool {
	var jj [][]string
	json.Unmarshal([]byte(selecTime), &jj)
	nowt := time.Now()
	dl := nowt.Weekday().String()
	daytime := TODAY_WEEK_NUM[dl]
	htime := nowt.Hour()         //获取当前时
	mintime := nowt.Minute()     //获取当前的分钟
	skip := mintime / 30         //当前的时偏移量
	hanftindex := htime*2 + skip //半小时索引
	//fmt.Println("是否处于投放时间", jj[daytime][hanftindex]) //获取当前是否投放广告 tf_style 投放的方式全天
	//当前的如果不满组条件的时候创建一个新的文件
	//	tp, _ := strconv.Atoi( jj[daytime][hanftindex]))
	//检查当前的投放是不是全天的投放的模式
	if jj[daytime][hanftindex] == "0" {
		return false
	} else {
		return true
	}
}

//找到满足条件的广告创意
//查满足条件的广告创意

func (self *Today) Chuangyi(pids []string) []byte {
	bid := self.BID_REQUEST

	pidsta := []string{}
	for _, v := range pids {
		if v != "" {
			pidsta = append(pidsta, v)
		}
	}
	pidst := strings.Join(pidsta, ",")

	sql := fmt.Sprintf("select * from tf_plan as p inner join tf_ads as s on p.id = s.plan_id where p.id in(%s) and p.is_off = 0 and s.is_off = 0 ", pidst)

	mydb := GetMysqlDb()
	dat, isdat := mydb.Query(sql) //广告创意集合

	if isdat == false {
		fmt.Println(sql)
		//当前的时间限制的
		fmt.Println("mysql not found some message 123")
	}

	//满足条件的创意的合
	var Ads_Array []map[string]string

	adsNum := bid.GetAdslots()

	//广告位置的合计
	//var Arr_Seats []*SeatBid //广告计划
	var plan_ads_arr map[string][]map[string]string = make(map[string][]map[string]string, len(dat))
	//广告请求位置
	for _, ads := range adsNum {

		for _, vm := range ads.AdType {
			kv := vm.String()
			knum := today.AdType_value[kv]
			for _, vv := range dat { //广告创意
				//广告创意
				ch_type := vv["ad_type"] //创意 投放的位置
				ad_type_arr := strings.Split(ch_type, ",")
				for _, tv := range ad_type_arr {
					tv_num, _ := strconv.Atoi(tv)
					if tv_num == int(knum) { //满足条件的创意
						vv["ads_type"] = tv
						Ads_Array = append(Ads_Array, vv)
						if len(plan_ads_arr[vv["plan_id"]]) == 0 {
						}

						plan_ads_arr[vv["plan_id"]] = append(plan_ads_arr[vv["plan_id"]], vv)
					}
				}
			}

		}
	}

	if len(plan_ads_arr) == 0 {
		return []byte{}
	}
	return []byte{}
	//查询满足条件的创意信息
	//bydat := bid.BbqRequest(plan_ads_arr)
	//return bydat
}

//生成广告创意合计

func (self *Today) Backads(pids []string) []byte {
	//bid := self.BID_REQUEST
	if len(pids) == 0 {
		return []byte{}
	}

	pidsta := []string{}
	for _, v := range pids {
		if v != "" {
			pidsta = append(pidsta, v)
		}
	}

	pidst := strings.Join(pidsta, ",")
	sql := fmt.Sprintf("SELECT ads.id,ads.type,ads.ad_type,im.*,ads.source,ads.detail_url,ads.plan_id,plan.link_url from tf_plan as plan inner join tf_ads as ads on plan.id =ads.plan_id INNER JOIN tf_images as im ON ads.id = im.ads_id WHERE ads.plan_id in (%s) AND ads.is_off =0 AND im.status = %q", pidst, "normal")

	mydb := GetMysqlDb()
	dat, isdat := mydb.Query(sql) //满足条件的创意``

	if isdat == false {
		//当前的时间限制的
		fmt.Println("mysql not found some message 23")
	}
	//满足条件的创意的合
	//var Ads_Array []map[string]string

	return self.CreatProf(dat)
	//return []byte{}
}

//生成probuf格式数据

func (self *Today) CreatProf(dat []map[string]string) []byte {

	if len(dat) == 0 {
		mk := []byte{}
		return mk
	}

	var domain string = "http://jrtt.qcwanwan.com"
	var param string = "notify?user_id={user_id}&request_id={request_id}&adid={adid}&bid_price={bid_price}&ip={ip}&timestamp={timestamp}&did={did}&imei={IMEI}&idfa={IDFA}&os={OS}&g_pos={g_pos}"
	var imgdom string = ""

	adslos := self.BID_REQUEST.Adslots
	var Arr_Seats []*today.SeatBid

	for _, adslocation := range adslos { //广告位每个广告位对象一个set读

		//生成一个新广告
		Seatbids := &today.SeatBid{}
		var Ads_arr []*today.Bid  //
		for _, ads := range dat { //广告创意数组
			ads_type := ads["ad_type"]
			ad_type_arr := strings.Split(ads_type, ",")
			//查询当前的计划id

			plansql := fmt.Sprintf("select * from tf_plan where id = %s", ads["plan_id"])
			mydb := GetMysqlDb()
			datplam, isdat := mydb.Query(plansql) //满足条件的创意

			if isdat == false {
				//当前的时间限制的
				fmt.Println("get plan message 1")
			}

			//广告创意
			for _, v := range ad_type_arr {

				//创意生成-------------------------------------------------------------
				biddata := &today.Bid{}
				ida_t := time.Now().Unix()
				d_yime := strconv.FormatInt(ida_t, 10)
				biddata.Id = &d_yime //生成唯一的商品的信息
				//创意的ID
				cid, _ := strconv.ParseUint(d_yime, 10, 64)

				//获取当前金额
				tf_money := SftoI(datplam[0]["first_price"])

				biddata.Adid = &cid //参与竞价的广告ID
				umprice := uint32(tf_money)

				//pp := adslocation.GetBidFloor() + 1
				//biddata.Price = &pp
				biddata.Price = &umprice
				//biddata.Price = adslocation.BidFloor + 1
				biddata.AdslotId = adslocation.Id
				cidstring := strconv.FormatUint(cid, 10)
				biddata.Cid = &cidstring //扩展的广告ID

				modelste := &today.MaterialMeta{} //广告素材对象
				modelste.AdType = getAdtype(v)
				apparam := fmt.Sprintf("%s&plan_id=%s&ads=%s", param, ads["plan_id"], ads["id"])
				nurlstr := fmt.Sprintf("%s/%s/%s", domain, "win", apparam)
				modelste.Nurl = &nurlstr
				titstr := ads["title"]
				modelste.Title = &titstr
				sourcestr := ads["source"]
				modelste.Source = &sourcestr

				//获取传递过来的banner

				//banner := adslots.Banner
				//banner_dat := banner[stk_v]

				//banner的图片信息显示
				imgbanner := &today.MaterialMeta_ImageMeta{}
				width_32, _ := strconv.ParseUint(ads["width"], 10, 32)
				height_32, _ := strconv.ParseUint(ads["height"], 10, 32)
				wid_addr := uint32(width_32)
				height_addr := uint32(height_32)
				imgbanner.Width = &wid_addr
				imgbanner.Height = &height_addr
				ads_title := ads["title"]

				imgbanner.Description = &ads_title

				if ads["type"] != "3" {
					im_url := fmt.Sprintf("%s%s", imgdom, ads["img_url"])
					imgbanner.Url = &im_url
					imgbanner.Urls = []string{fmt.Sprintf("%s%s", imgdom, ads["img_url"])}
				} else {
					var imgsA []string
					dsd := strings.Split(ads["img_url"], ",")
					for _, v := range dsd {
						imgsA = append(imgsA, fmt.Sprintf("%s%s", imgdom, v))
					}
					imgurl := imgsA[0]
					imgbanner.Url = &imgurl
					imgbanner.Urls = imgsA

				}
				modelste.ImageBanner = imgbanner
				//设置当前的操作流程

				dsp_external := &today.MaterialMeta_ExternalMeta{}
				detailurl := dat[0]["detail_url"]
				//exter_url := ads["link_url"]
				dsp_external.Url = &detailurl
				modelste.External = dsp_external
				biddata.Creative = modelste

				modelste.ShowUrl = []string{fmt.Sprintf("%s/%s/%s", domain, "show", apparam)}
				modelste.ClickUrl = []string{fmt.Sprintf("%s/%s/%s", domain, "click", apparam)}

				bjk := []string{"3", "4", "7", "9", "18", "4"}
				//获取当前系统
				for _, vtype := range bjk {
					if v == vtype {
						os := self.BID_REQUEST.Device.Os
						switch *os {
						case "ios":
							iosDat := &today.MaterialMeta_IosApp{}
							ios_downurl := datplam[0]["link_url"]
							iosDat.DownloadUrl = &ios_downurl
							ios_appname := datplam[0]["app_name"]
							iosDat.AppName = &ios_appname
							modelste.IosApp = iosDat
						case "android":

							android_url := datplam[0]["link_url"]
							android := &today.MaterialMeta_AndroidApp{}
							android_appname := datplam[0]["app_name"]
							android.AppName = &android_appname
							android.DownloadUrl = &android_url
							android.WebUrl = &detailurl
							modelste.AndroidApp = android
						}
					}
				}
				//创意生成结束------------------------------------------------------------

				Ads_arr = append(Ads_arr, biddata) //创意数组

			}

		}
		Seatbids.Ads = Ads_arr

		//生成新的广告位
		Arr_Seats = append(Arr_Seats, Seatbids)
	}

	res := &today.BidResponse{}
	res.Seatbids = Arr_Seats

	res.RequestId = self.BID_REQUEST.RequestId
	data, err := proto.Marshal(res)

	if err != nil {
		return []byte{}
	}
	return data

	//return []byte{}
}

func getAdtype(num string) *today.AdType {

	//	return modelste
	n, _ := strconv.ParseUint(num, 10, 32)
	var th today.AdType = today.AdType(n)
	return &th
}

//查询当前广告组的限额

func GroupLimitMoney() {
	//
	groupsql := "SELECT id FROM tf_ads_group WHERE is_off = 0 AND `status` =0  "
	mysqldb := GetMysqlDb()
	gdat, _ := mysqldb.Query(groupsql) //获取当前的广告
	fmt.Println(gdat)
	for _, vb := range gdat {
		fmt.Println(vb)
	}
}

type Bidprce struct {
	Bprice int64 `"bson:bidprce"`
}

//获取当前计划的CPM的付费的方式的的计费总额 tp 0 表示按天的限额 1总的限额 计划的限额的
func CpmMoneyLimit(dat map[string]string, tp int, bbq chan bool) {

	nowt := time.Now()
	timestamp := time.Date(nowt.Year(), nowt.Month(), nowt.Day(), 0, 0, 0, 0, time.Local) //当日的时间戳
	start_unix_time := timestamp.Unix()
	sd := strconv.FormatInt(start_unix_time, 10)

	planid := dat["id"]
	var mongotable string = "win_Notify"
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("there had some error")
		}
	}()
	mongo := GetMongoSession().Copy()
	defer mongo.Close()

	num, err := mongo.DB(MongodbConf.DataBase).C(mongotable).Find(bson.M{"planid": planid, "timestamp": bson.M{"$gte": sd}}).Count()
	fmt.Println("查询到当前的计划总数", num, sd)

	if err != nil {
		fmt.Println("There had some error")
	}

	//如果当前没有相关的数据返回true 表示条件当前是满足的
	if num == 0 {
		bbq <- true
	}

	bidprce := Bidprce{}
	//获取当前的bidprce
	mongo.DB(MongodbConf.DataBase).C(mongotable).Find(bson.M{"planid": planid, "timestamp": bson.M{"$gte": sd}}).One(&bidprce)

	fmt.Println("bid prrice=====>", bidprce.Bprice, planid)
	bbprce := bidprce.Bprice
	oneprice := bbprce / 1000.0 //每个单价
	fmt.Println("每个展示的单价", oneprice)
	strconv.ParseFloat(strconv.FormatInt(bbprce, 10), 10)

	//return bidprce
}
