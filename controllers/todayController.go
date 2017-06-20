package controllers

import (
	"Tdsp/data/today"
	"Tdsp/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"

	"time"

	"github.com/golang/protobuf/proto"

	"github.com/julienschmidt/httprouter"
)

func init() {

}

const (
	TJ_FEED_DSP_ID  string = "1756165498"
	TJ_FEED_DSP_KEY string = "a74e696576394976bc694fbd58a2b0d6"
	XQ_FEED_DSP_ID  string = "1756165504"
	XQ_FEED_DSP_KEY string = "c8c2df73229b47378f8eceb9cc12ea1a"
	DZ_FEED_DSP_ID  string = "1756165502"
	DZ_FEED_DSP_KEY string = "f49916e2447f490d93822dff5c345aa3"
)

type NotifyDat struct {
	UserId    string
	RequestID string
	Adid      string
	BidPrce   uint64
	Ip        string
	TimesTamp string
	Did       string
	PlanId    string
	Ads       string
	Imei      string
	Idfa      string
	Os        string
	Win_num   int
}

//宏替换的数据

func TodayBidRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("now request")
	st_time := time.Now().UnixNano()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("出错了")
			fmt.Println(err)
		}

	}()

	httputil.DumpRequest(r, true)
	bydata, _ := ioutil.ReadAll(r.Body)

	reqa := &today.BidRequest{}

	era := proto.Unmarshal(bydata, reqa)
	if era != nil {
		fmt.Println("转化错误iiiii")
		panic("解析错误")
	}

	todaymodel := models.Today{}
	data := todaymodel.GetPlans(reqa) //回去相关的广告计划

	//	data := models.GetPlan(reqa)
	w.Write(data)

	newTest := &today.BidResponse{}
	err := proto.Unmarshal(data, newTest)
	if err != nil {
		fmt.Println("----->>>")
	}

	fmt.Println(newTest)

	end_time := time.Now().UnixNano()

	fmt.Println("\r\n useTime:", (end_time-st_time)/1e6)
	btext := newTest.String()
	if len(btext) != 0 {
		go saveReponseDat((end_time-st_time)/1e6, btext)
	}

}

//监测地址数据

var JCdat []map[string]string = make([]map[string]string, 1)

func saveData(collectionName string, ps url.Values) {

	dat := NotifyDat{}
	dat.Adid = ps.Get("adid")
	dat.BidPrce = Decprice(ps.Get("bid_price"))
	dat.Did = ps.Get("did")
	dat.Ip = ps.Get("ip")
	dat.RequestID = ps.Get("request_id")
	dat.TimesTamp = ps.Get("timestamp")
	dat.UserId = ps.Get("user_id")
	dat.PlanId = ps.Get("plan_id")
	dat.Ads = ps.Get("ads")
	dat.Imei = ps.Get("imei")
	dat.Idfa = ps.Get("idfa")
	dat.Os = ps.Get("os")
	win_num, _ := strconv.Atoi(ps.Get("g_pos"))
	dat.Win_num = win_num
	fmt.Println(dat)
	//models.MongoDb.C("win_Notify").Insert(dat)

	ads_sql := fmt.Sprintf("select show_url,click_url from tf_ads where id = %s", ps.Get("ads"))

	fmt.Println("ads_sql", ads_sql)
	mydb := models.GetMysqlDb()
	JCdat, _ = mydb.Query(ads_sql)

	mongo := models.GetMongoSession().Copy()
	defer mongo.Close()
	mongo.DB(models.MongodbConf.DataBase).C(collectionName).Insert(dat)
}

//获胜后获取的数据github.com/djimenez/iconv-go"

func WinRequest(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	fmt.Println("WIN")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER win")
			fmt.Println(err)
		}
	}()
	ps := r.URL.Query()
	go saveData("win_Notify", ps)
}

//展示的时候请求的数据

func ShowRequest(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	fmt.Println("Show ")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER show")
			fmt.Println(err)
		}
	}()
	ps := r.URL.Query()
	//fmt.Println(ps)
	go saveData("show_Notify", ps)

	if len(JCdat) > 0 {
		showurl := JCdat[0]["show_url"]
		//Hreaplace(showurl)
		go senddata(psa, showurl)
	}

}

//点击的时候展示

func ClickRequest(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER click")
			fmt.Println(err)
		}
	}()
	fmt.Println("click")
	ps := r.URL.Query()
	saveData("click_Notify", ps)
	clickurl := JCdat[0]["click_url"]
	//Hreaplace(showurl)
	go senddata(psa, clickurl)
	//go func(ps httprouter.Params) {}(ps)

}

func senddata(ps httprouter.Params, sendurl string) {
	if len(sendurl) == 0 {
		return
	}
	hdat := map[string]string{}
	reg, _ := regexp.Compile(`__\w+__`)
	hdat["__IDFA__"] = ps.ByName("idfa")
	hdat["__IMEI__"] = ps.ByName("imei")
	hdat["__IP__"] = ps.ByName("ip")
	hdat["__TS__"] = ps.ByName("ts")
	hdat["__OS__"] = ps.ByName("os")
	sendstr := reg.ReplaceAllStringFunc(sendurl,
		func(b string) string {
			dh := hdat[b]
			return dh
		})
	Httprequest(sendstr, "GET", "")
}

func Index(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	fmt.Println("This is dsp Index page")
	w.Write([]byte("Index func"))
}

//保存响应的数据

func saveReponseDat(utime int64, content string) {
	dt := map[string]interface{}{"useTime": utime, "content": content, "addtime": time.Now().Unix()}
	mongo := models.GetMongoSession().Copy()
	mongo.DB(models.MongodbConf.DataBase).C("backdat").Insert(dt)
}
