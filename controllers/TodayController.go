package controllers

import (
	"dsp/models"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//接收第一次传递过来的数据
func RequestToday(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There had some error!")
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

}

//获胜后获取的数据

func WinRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	models.MongoDb.C("dsp_win").Insert(r.Header)

}

//展示的时候请求的数据

func ShowNotify(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

//点击的时候展示

func ClickNotify(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
