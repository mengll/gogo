// controllers project controllers.go
package controllers

import (
	"dsp/models"
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/golang/protobuf/proto"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("dsp controllers show ")

}

func init() {
	//	bytesa := []byte{10, 32, 50, 48, 49, 55, 48, 53, 48, 52, 48, 55, 51, 54, 53, 52, 49, 55, 50, 48, 49, 55, 49, 53, 50, 48, 48, 52, 54, 56, 54, 53, 68, 70, 18, 151, 6, 10, 148, 6, 10, 24, 53, 57, 48, 97, 54, 57, 57, 54, 50, 57, 101, 101, 100, 102, 50, 57, 53, 48, 52, 100, 102, 54, 56, 100, 18, 16, 102, 48, 102, 52, 52, 57, 56, 97, 101, 56, 99, 99, 52, 57, 57, 57, 24, 217, 4, 32, 150, 211, 169, 200, 149, 211, 193, 246, 2, 42, 196, 5, 8, 11, 18, 151, 1, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 119, 105, 110, 47, 110, 111, 116, 105, 102, 121, 63, 117, 115, 101, 114, 95, 105, 100, 61, 123, 117, 115, 101, 114, 95, 105, 100, 125, 38, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 61, 123, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 125, 38, 97, 100, 105, 100, 61, 123, 97, 100, 105, 100, 125, 38, 98, 105, 100, 95, 112, 114, 105, 99, 101, 61, 123, 98, 105, 100, 95, 112, 114, 105, 99, 101, 125, 38, 105, 112, 61, 123, 105, 112, 125, 38, 116, 105, 109, 101, 115, 116, 97, 109, 112, 61, 123, 116, 105, 109, 101, 115, 116, 97, 109, 112, 125, 38, 100, 105, 100, 61, 123, 100, 105, 100, 125, 26, 49, 231, 153, 189, 233, 135, 142, 231, 140, 170, 231, 136, 134, 228, 186, 134, 228, 184, 170, 232, 163, 133, 229, 164, 135, 46, 230, 141, 162, 228, 186, 134, 229, 176, 143, 229, 141, 138, 228, 184, 170, 230, 156, 136, 229, 183, 165, 232, 181, 132, 34, 12, 228, 188, 160, 229, 165, 135, 230, 151, 160, 229, 143, 140, 42, 134, 1, 16, 228, 1, 24, 150, 1, 34, 30, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 49, 46, 106, 112, 103, 42, 30, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 49, 46, 106, 112, 103, 42, 30, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 50, 46, 106, 112, 103, 42, 30, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 51, 46, 106, 112, 103, 50, 37, 10, 35, 104, 116, 116, 112, 58, 47, 47, 109, 46, 97, 110, 102, 101, 110, 103, 46, 99, 110, 47, 99, 113, 119, 115, 95, 98, 98, 107, 45, 105, 111, 115, 47, 49, 50, 47, 74, 152, 1, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 115, 104, 111, 119, 47, 110, 111, 116, 105, 102, 121, 63, 117, 115, 101, 114, 95, 105, 100, 61, 123, 117, 115, 101, 114, 95, 105, 100, 125, 38, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 61, 123, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 125, 38, 97, 100, 105, 100, 61, 123, 97, 100, 105, 100, 125, 38, 98, 105, 100, 95, 112, 114, 105, 99, 101, 61, 123, 98, 105, 100, 95, 112, 114, 105, 99, 101, 125, 38, 105, 112, 61, 123, 105, 112, 125, 38, 116, 105, 109, 101, 115, 116, 97, 109, 112, 61, 123, 116, 105, 109, 101, 115, 116, 97, 109, 112, 125, 38, 100, 105, 100, 61, 123, 100, 105, 100, 125, 82, 153, 1, 104, 116, 116, 112, 58, 47, 47, 106, 114, 116, 116, 46, 113, 99, 119, 97, 110, 119, 97, 110, 46, 99, 111, 109, 47, 99, 108, 105, 99, 107, 47, 110, 111, 116, 105, 102, 121, 63, 117, 115, 101, 114, 95, 105, 100, 61, 123, 117, 115, 101, 114, 95, 105, 100, 125, 38, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 61, 123, 114, 101, 113, 117, 101, 115, 116, 95, 105, 100, 125, 38, 97, 100, 105, 100, 61, 123, 97, 100, 105, 100, 125, 38, 98, 105, 100, 95, 112, 114, 105, 99, 101, 61, 123, 98, 105, 100, 95, 112, 114, 105, 99, 101, 125, 38, 105, 112, 61, 123, 105, 112, 125, 38, 116, 105, 109, 101, 115, 116, 97, 109, 112, 61, 123, 116, 105, 109, 101, 115, 116, 97, 109, 112, 125, 38, 100, 105, 100, 61, 123, 100, 105, 100, 125, 58, 18, 50, 49, 48, 56, 51, 50, 48, 49, 51, 50, 52, 57, 56, 57, 56, 57, 48, 50}

	//	dec := &models.BidResponse{}
	//	er := proto.Unmarshal(bytesa, dec)
	//	if er != nil {
	//		fmt.Println("转化错误")
	//	}

	//	//bid := &models.SeatBid{}
	//	bid := dec.GetSeatbids()
	//	for _, v := range bid { //
	//		m := v.GetAds()
	//		bid := m[0]

	//		fmt.Println(bid)
	//		for _, bn := range m {
	//			fmt.Println(bn.Adid)
	//		}
	//		fmt.Println(v.GetSeat())
	//	}

	//hja := proto.UnmarshalText(dec.Seatbids, bid)

	//fmt.Println(bid)
	fmt.Println("<<------>")
	dat, err := ioutil.ReadFile("bit.req")
	if err != nil {
		fmt.Println(err)
	}

	req := &models.BidRequest{}
	era := proto.Unmarshal(dat, req)
	if era != nil {
		fmt.Println("转化错误")
	}
	fmt.Println(req)
	adts := req.GetAdslots()
	fmt.Println(adts)
	fmt.Println("------>")

	TuiHead()
}

//响应返回的数据 推荐流
//models.BidResponse.Seatbids
//models.SeatBid.Ads
//models.Bid.Creative.ImageBanner.Urls

func TuiHead() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("This had some Error!")
		}
	}() //显示相关的错误信息

	biddata := &models.Bid{}
	biddata.Adid = 2323423423

	modelste := &models.MaterialMeta{}
	modelste.AdType = 11
	modelste.Nurl = ""
	modelste.Title = "怕冷的姑娘一定要看"
	modelste.Source = "阿里妈妈"

	//banner的图片信息显示
	imgbanner := &models.MaterialMeta_ImageMeta{}
	imgbanner.Width = 228
	imgbanner.Height = 150
	imgbanner.Url = "https://otrade.alicdn.com/8fb73fad38124adcbacf5019cf2c1941.jpg"
	imgbanner.Urls = []string{"https://otrade.alicdn.com/8fb73fad38124adcbacf5019cf2c1941.jpg",
		"https://otrade.alicdn.com/c8fcf15b9e1141ecae7a225b59bc4e2a.jpg",
		"https://otrade.alicdn.com/93a6ab4f6a9c4baeba1fb9e081aaf1dc.jpg"}
	modelste.ImageBanner = imgbanner

	//设置当前的操作流程
	dsp_external := &models.MaterialMeta_ExternalMeta{}
	dsp_external.Url = "http://click.tanx.com/oc?e=F31NInZSPJV%2fj7ecl7%2bKGxsQmLP5zomMnn7luM%2fSmJa2WyNXBSkEIFCLPjtHgr4%2bPhZRdbjT8909FfaWGEkZqG1npCuXOU5%2fs2kJvWidcxrb6tTVDmweIeRMyz8MsisZXPmCHlFobeSwTwo53V4MB8ytpBIWKRFsOTuoi7EExlejH3q7mm1A1lAVvHwxODoemZX4wPvZJ3Mgyh1K568EMDmwIyJ%2fLbi5\u0026u=https%3a%2f%2fhuodong.taobao.com%2fwow%2ftb-20161212%2fact%2fzhuhuichang%3ftanx_bid%3d0a67155200002c1e584948fb07e96695%26cid%3d45774%26spm%3da21bo.50862.523823.1.lQHX3w%26wh_weex%3dtrue%26wh_prefetch_id%3dact-zhuhuichang%26mm_gxbid%3d1_1220490_d320e681b5d06519d8e7ae56736e76b3\u0026k=225"
	modelste.External = dsp_external
	biddata.Creative = modelste

	modelste.ShowUrl = []string{"http://df.tanx.com/spf?e=A0O2bOF0juc6sENrlwpRJ8FbrlrFAVKGHkWt%2ByPw88T9kLR7EOzwHo0At2bA7sdfN%2BAUXkJkpDshivnk0X8uI5OxcwsuaDrs9RxAh35XNNdsHHiCO6zMt2iQt%2BqXgptzJMIjqlCuzWAExVQzNGErjsqChaHlR9AjNCbTXtTlWrxfW8f0H6Ym5w%2BWZiRUwAqCGWz8OQhOHfph58HvU8ZxTxsQmLP5zomM\u0026u=http%3A%2F%2Fef.tanx.com%2Fgateway%3Fch%3Dtanx%26p%3DAQpnHYEAA1hJSPthXwBYC3ejTPHXm%252Bdhzw%253D%253D%26e%3D3TAqTv%252fRTGA6sENrlwpRJ8FbrlrFAVKGHkWt%252byPw88RgkZDSsfWBsVNxl47FXkEs1p27wUp6NF%252fIQYP9SXPRLvWFWOKEOEmy5qb1xC5iX9ctP1muuE%252fAMyjQv%252fv1Fr3jXzssD%252fdf%252bKL1hVjihDhJsuoj4IhejEpQ67lY2eMh%252bwOU5Hb8ObgDK%252bycHbwkBplG5GZxuZhvxvnKuZKYeTfFsynYxSU7yoqtbmfNNBiTVW7zJRN5cYE1qO4M7fF34YaN%26k%3D257\u0026k=225\u0026p={bid_price}",
		"http://gxb.mmstat.com/gxb.gif?tanx_bid=0a67155200002c1e584948fb07e96695\u0026si=1220490\u0026ref=\u0026lang=undefined\u0026bw=0\u0026bh=0\u0026pu=\u0026ht=pageview\u0026di=\u0026dim=\u0026dud=0"}

	modelste.ClickUrl = []string{"http://rdstat.tanx.com/trd?f=https%3a%2f%2fhuodong.taobao.com%2fwow%2ftb-20161212%2fact%2fzhuhuichang\u0026k=748569918d3cfb0d\u0026p=mm_26632268_8510653_67720319\u0026pvid=0a671d810003584948fb615f00580b77\u0026s=228x150\u0026d=111072246\u0026did=2012876\u0026t=1481197819",
		"http://gxb.mmstat.com/gxb.gif?tanx_bid=0a67155200002c1e584948fb07e96695\u0026t=https%3A%2F%2Fhuodong.taobao.com%2Fwow%2Ftb-20161212%2Fact%2Fzhuhuichang%3Fspm%3Da21bo.50862.523823.1.lQHX3w%26wh_weex%3Dtrue%26wh_prefetch_id%3Dact-zhuhuichang%26mm_gxbid%3D1_1220490_d320e681b5d06519d8e7ae56736e76b3\u0026v=01962da65b36\u0026di=\u0026dim=\u0026dud=0"}

	modelste.IsInapp = false

	biddata.Dealid = "4444444444"
	biddata.Cid = "111072246_45774_15475717"

	dsp := &models.SeatBid{}
	dsp.Ads = []*models.Bid{biddata}

	res := &models.BidResponse{}
	res.Seatbids = []*models.SeatBid{dsp}
	res.RequestId = "2016120819501901000404220878378B"
	res.ErrorCode = 0

	data, err := proto.Marshal(res)

	if err != nil {
		fmt.Println("转化错误")
	}
	fmt.Println("---------------------m------------------")
	fmt.Println(data)
	fmt.Println("--------------------e-m-----------------")

	newTest := &models.BidResponse{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {

	}

	fmt.Println("================================m========================")
	fmt.Println(newTest)
	fmt.Println("=========================e-m=============================")
	fmt.Println(res)
	//fmt.Println(biddata)
}
