package controllers

import (
	"dsp/models"
	"fmt"

	"net/http"

	"io/ioutil"
	"math/rand"
	"net/http/httputil"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/url"
	"strconv"
	"time"

	"encoding/json"

	"github.com/golang/protobuf/proto"

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

func saveData(collectionName string, ps url.Values) {

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

//保存到mysql数据库中

func sqlMysqlData(ntp int, ps url.Values) {
	sql := fmt.Sprintf("insert into tf_notify (`userid`,`requestid`,`adid`,`bidprce`,`ip`,`timestamp`,`did`,`type`)values(%s,%s,%s,%d,%s,%s,%s,%s,%d)", ps.Get("user_id"), ps.Get("request_id"), ps.Get("adid"), Decprice(ps.Get("bid_price")), ps.Get("ip"), ps.Get("timestamp"), ps.Get("did"), ntp)
	mydb := models.GetMysqlDb()
	mydb.Insert(sql)
	fmt.Println(sql)
}

//获胜后获取的数据github.com/djimenez/iconv-go"

func WinRequest(w http.ResponseWriter, r *http.Request, psa httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some ER")
			fmt.Println(err)
		}
	}()
	ps := r.URL.Query()
	saveData("win_Notify", ps)
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
	saveData("show_Notify", ps)

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
	saveData("click_Notify", ps)
}

const (
	TJ_FEED_DSP_ID  string = "1756165498"
	TJ_FEED_DSP_KEY string = "a74e696576394976bc694fbd58a2b0d6"
	XQ_FEED_DSP_ID  string = "1756165504"
	XQ_FEED_DSP_KEY string = "c8c2df73229b47378f8eceb9cc12ea1a"
	DZ_FEED_DSP_ID  string = "1756165502"
	DZ_FEED_DSP_KEY string = "f49916e2447f490d93822dff5c345aa3"
)

type QueryCreatives struct {
	dspid     string
	dspkey    string
	startDate string
	endDate   string
	timestamp string
}

//send request dat

func initas() {
	fmt.Println("1")
	//nowTime := time.Now()

	mm := map[string]string{"dspkey": TJ_FEED_DSP_KEY, "enddate": "2017-05-14", "startdate": "2017-05-12", "timestamp": strconv.FormatInt(time.Now().Unix(), 10), "dspid": TJ_FEED_DSP_ID}
	QueryTtNum(mm)
}

var Qudat map[string]string

//get the data
func QueryTtNum(qudat map[string]string) {

	var str string
	for k, v := range qudat {
		if k != "dspkey" {
			str += fmt.Sprintf("%s=%s&", k, v)
		}
	}
	tdt := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println(tdt)
	str = fmt.Sprintf("dspid=%s&timestamp=%s&startdate=%s&enddate=%s&", qudat["dspid"], tdt, qudat["startdate"], qudat["enddate"])
	var urls string = fmt.Sprintf("http://adx.toutiao.com/adxbuyer/api/v1.0/creatives/stat?%s", str)
	uu := fmt.Sprintf("http://adx.toutiao.com/adxbuyer/api/v1.0/creatives/stat?dspid=%s&timestamp=%s&startdate=%s&enddate=%s", qudat["dspid"], tdt, qudat["startdate"], qudat["enddate"])
	//var urls = "http://adx.toutiao.com/adxbuyer/api/v1.0/creatives/stat?dspid=1756165498&timestamp=1495459053&startdate=2017-04-01&enddate=2017-04-27&"

	mla := "http://adx.toutiao.com/adxbuyer/api/v1.0/creatives/stat?dspid=1756165498&timestamp=1495472852&startdate=2017-05-12&enddate=2017-05-14"

	fmt.Println(qudat["dspkey"])
	stra := CreatSign(qudat["dspkey"], uu)

	urls = urls + "signature=" + stra
	fmt.Println(urls)

	fmt.Println("---->><><<<")
	kl := CreatSign(qudat["dspkey"], mla)
	fmt.Println(kl)
	httpGet(urls)

	//get the data
}

func httpGet(urls string) {
	resp, err := http.Get(urls)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
	dat := models.DataChange(string(body))
	if dat["error"] != nil {
		return
	}
	dda := dat["result"]

	//byte change to utf8
	var hj string = dda.(string)
	bba := []byte(hj)
	//TtDecodecode(hj)
	fmt.Println("--->>", string(bba))
	aesEnc := AesEncrypt{}
	rs, _ := aesEnc.Decrypt(bba)

	var mm map[string]interface{}
	json.Unmarshal(rs, &mm)
	fmt.Println(mm)
	fmt.Println("-<><<<<<<<<")
	fmt.Println(string(rs))

	// {"date_stats"=>[{"click_cnt"=>0, "cost"=>0.02, "req_cnt"=>78, "date"=>"2017-04-26", "win_cnt"=>4, "bid_cnt"=>5, "show_cnt"=>7}, {"click_cnt"=>6, "cost"=>0.11, "req_cnt"=>56, "date"=>"2017-04-27", "win_cnt"=>27, "bid_cnt"=>28, "show_cnt"=>39}], "local_ads_stats"=>{"click_cnt"=>0, "bid"=>4.0, "show_cnt"=>0, "cost"=>0.0, "extra"=>{"daily_stats"=>{}}}}
	//bf := bytes.NewBuffer(rs)
	//ss := bf.String()
	var anydata QuerDat
	json.Unmarshal(rs, &anydata)
	fmt.Println("----------------------------")
	fmt.Println(anydata)

}

type TTDayValue struct {
	ClickCnt int     `json:"click_cnt"`
	Cost     float32 `json:"cost"`
	ReqCnt   int32   `json:"req_cnt"`
	Date     string  `json:"date"`
	WinCnt   int     `json:"win_cnt"`
	BidCnt   int     `json:"bid_cnt"`
	ShowCnt  int     `json:"show_cnt"`
}

type QuerDat struct {
	DateStatus    []TTDayValue
	LocalAdsStats TTDayValue
	Extra         map[string]interface{}
}

func (this *AesEncrypt) getKey() []byte {
	strKey := "a74e696576394976bc694fbd58a2b0d6"
	keyLen := len(strKey)
	fmt.Println(keyLen)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(strKey)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

type AesEncrypt struct {
}

//加密字符串
func (this *AesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))
	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(src []byte) (strDesc []byte, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	key := this.getKey()
	var iv = src[:aes.BlockSize]
	clipdt := src[aes.BlockSize:]
	//decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(clipdt, clipdt) //d out src input
	return clipdt, nil
}

//decode data

func TtDecodecode(dt string) {

	data, err := base64.StdEncoding.DecodeString(dt)
	if err != nil {
	}

	iv := data[0:16]
	//dat := data[16:]

	keys := []byte(TJ_FEED_DSP_KEY)

	decrypted := make([]byte, aes.BlockSize+len(data))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(keys)
	if err != nil {

	}

	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)

	aesDecrypter.XORKeyStream(decrypted, data)

	fmt.Println(decrypted)

	//	fmt.Println(keys)
	//	asa, _ := aes.NewCipher(keys)
	//	aas := cipher.NewCFBEncrypter(asa, iv)
	//	aas.XORKeyStream()
	//	fmt.Println(aas)
	//	fmt.Println(dat, iv)

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
