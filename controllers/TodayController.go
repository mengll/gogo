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
