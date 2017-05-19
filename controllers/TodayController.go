package controllers

import (
	"dsp/models"
	"fmt"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

//the Type of notify
type NotifyDat struct {
	UserId    string
	RequestID string
	Adid      string
	BidPrce   uint64
	Ip        string
	TimesTamp string
	Did       string
}

//接收第一次传递过来的数据
func RequestToday(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some error!")
			fmt.Println(err)
		}
	}()

	fmt.Println("This is request function1")
	sql := fmt.Sprintf("select rid,pid from ucusers where uid = %q or mobile =%q", "18827092404", "18827092404")
	fmt.Println(sql)
	dat, ok := models.Mydb.Query(sql)
	fmt.Println(dat)
	fmt.Println(ok)

	Log.Infof("用户信息", dat)

	smt, err := models.Mydb.DB.Prepare("select rid,pid from ucusers where uid =? or mobile = ? ")
	if err != nil {
		fmt.Println("数据库连接错误")
	}
	res, erra := smt.Exec("18818818801", "18818818801")
	if erra != nil {
		fmt.Println("执行出错了！")
	}
	fmt.Println(res.RowsAffected()) //只是单纯的执行的了查询的操作

	backData := TuiHead()

	//r.Response.Status = "200"
	//r.Response.Write(w)
	_, erro := fmt.Println(backData)
	if erro != nil {
		fmt.Println("This is go function")
	}

}

//获胜后获取的数据

func WinRequest(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER")
			fmt.Println(err)
		}
	}()
	ps := r.URL.Query()
	//fmt.Println("dsa", r.URL.Query().Get("adid"))
	//models.MongoDb.C("dsp_win").Insert(r.Header)
	dat := NotifyDat{}
	dat.Adid = ps.Get("adid")
	dat.BidPrce = Decprice(ps.Get("bid_price"))
	dat.Did = ps.Get("did")
	dat.Ip = ps.Get("ip")
	dat.RequestID = ps.Get("request_id")
	dat.TimesTamp = ps.Get("timestamp")
	dat.UserId = ps.Get("user_id")
	fmt.Println(dat)
	models.MongoDb.C("win_Notify").Insert(dat)

}

//展示的时候请求的数据

func ShowNotify(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER")
			fmt.Println(err)
		}
	}()
	ps := r.URL.Query()
	//fmt.Println("dsa", r.URL.Query().Get("adid"))
	//models.MongoDb.C("dsp_win").Insert(r.Header)
	dat := NotifyDat{}
	dat.Adid = ps.Get("adid")
	dat.BidPrce = Decprice(ps.Get("bid_price"))
	dat.Did = ps.Get("did")
	dat.Ip = ps.Get("ip")
	dat.RequestID = ps.Get("request_id")
	dat.TimesTamp = ps.Get("timestamp")
	dat.UserId = ps.Get("user_id")
	fmt.Println(dat)
	models.MongoDb.C("show_Notify").Insert(dat)

}

//点击的时候展示

func ClickNotify(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER")
			fmt.Println(err)
		}
	}()
	ps := r.URL.Query()
	//fmt.Println("dsa", r.URL.Query().Get("adid"))
	//models.MongoDb.C("dsp_win").Insert(r.Header)
	dat := NotifyDat{}
	dat.Adid = ps.Get("adid")
	dat.BidPrce = Decprice(ps.Get("bid_price"))
	dat.Did = ps.Get("did")
	dat.Ip = ps.Get("ip")
	dat.RequestID = ps.Get("request_id")
	dat.TimesTamp = ps.Get("timestamp")
	dat.UserId = ps.Get("user_id")
	fmt.Println(dat)
	models.MongoDb.C("click_Notify").Insert(dat)

}

//bidrequest

func BidRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
     defer func() {
                if err := recover(); err != nil {
                        fmt.Println("出错了")
                }

        }()

        dasa, _ := httputil.DumpRequest(r, true)
        fmt.Println(r.Body)
        fmt.Println(string(dasa))
        fmt.Println("--=====----")

                bydata ,_ :=ioutil.ReadAll(r.Body)


        reqa := &models.BidRequest{}
        era := proto.Unmarshal(bydata, reqa)
        if era != nil {
                        fmt.Println("转化错误iiiii")
        }
        fmt.Println(reqa)


	backData := TuiHead()

	//r.Response.Status = "200"
	//r.Response.Write(w)
	_, erro := fmt.Println(backData)
	if erro != nil {
		fmt.Println("This is go function")
	}

}

//这是响应今日头条的第一次请求的的调用的方法


func BbqRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("出错了")
		}

	}()

	httputil.DumpRequest(r, true)

	bydata, _ := ioutil.ReadAll(r.Body)

	reqa := &models.BidRequest{}
	era := proto.Unmarshal(bydata, reqa)
	if era != nil {
		fmt.Println("转化错误iiiii")
		panic("解析错误")
	}

	//响应的对象的数据
	rback := rand.New(rand.NewSource(time.Now().UnixNano()))

	adslots := reqa.GetAdslots()[0] //获取最上曾数据
	biddata := &models.Bid{}
	biddata.Id = strconv.Itoa(rback.Intn(20)) //生成唯一的商品的信息
	biddata.Adid = uint64(rback.Intn(18))
	biddata.Price = adslots.BidFloor + 1
	biddata.AdslotId = adslots.Id
	biddata.Cid = strconv.Itoa(rback.Intn(18))

	modelste := &models.MaterialMeta{}
	modelste.AdType = models.AdType_TOUTIAO_FEED_LP_GROUP
	modelste.Nurl = "http://jrtt.qcwanwan.com/win/notify?user_id={user_id}&request_id={request_id}&adid={adid}&bid_price={bid_price}&ip={ip}&timestamp={timestamp}&did={did}"
	modelste.Title = "白野猪爆了个装备.换了小半个月工资NB"
	modelste.Source = "传奇无双"

	//banner的图片信息显示
	imgbanner := &models.MaterialMeta_ImageMeta{}
	imgbanner.Width = 228
	imgbanner.Height = 150
	imgbanner.Url = "http://jrtt.qcwanwan.com/1.jpg"
	imgbanner.Urls = []string{"http://jrtt.qcwanwan.com/1.jpg",
		"http://jrtt.qcwanwan.com/2.jpg",
		"http://jrtt.qcwanwan.com/3.jpg"}
	modelste.ImageBanner = imgbanner

	//设置当前的操作流程
	dsp_external := &models.MaterialMeta_ExternalMeta{}
	dsp_external.Url = "http://m.anfeng.cn/cqws_bbk-ios/12/"
	modelste.External = dsp_external
	biddata.Creative = modelste

	modelste.ShowUrl = []string{"http://jrtt.qcwanwan.com/show/notify?user_id={user_id}&request_id={request_id}&adid={adid}&bid_price={bid_price}&ip={ip}&timestamp={timestamp}&did={did}"}

	modelste.ClickUrl = []string{"http://jrtt.qcwanwan.com/click/notify?user_id={user_id}&request_id={request_id}&adid={adid}&bid_price={bid_price}&ip={ip}&timestamp={timestamp}&did={did}"}

	dsp := &models.SeatBid{}
	dsp.Ads = []*models.Bid{biddata}

	res := &models.BidResponse{}
	res.Seatbids = []*models.SeatBid{dsp}
	res.RequestId = reqa.RequestId

	data, err := proto.Marshal(res)
	w.Write(data)

	newTest := &models.BidResponse{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		fmt.Println("----->>>")
	}
	fmt.Println(newTest)

}


