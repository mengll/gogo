// controllers project controllers.go
package controllers

import (
	"fmt"

	"io/ioutil"

	"Tdsp/data/today"
	"Tdsp/models"

	"github.com/golang/protobuf/proto"
)

func initas() {

	dat, err := ioutil.ReadFile("bit.req")

	if err != nil {
		fmt.Println(err)
	}

	req := &today.BidRequest{}
	//	fmt.Println(dat)
	era := proto.Unmarshal(dat, req)
	if era != nil {
		fmt.Println(era)
		fmt.Println("转化错误")
	}
	//	fmt.Println(req)
	//	adts := req.GetAdslots()
	//	fmt.Println(adts)
	tod := models.Today{}
	tod.GetPlans(req)
	//TuiHead()

}
