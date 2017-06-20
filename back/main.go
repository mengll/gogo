// dsp project main.go
package main

import (
	"net/http"

	"dsp/controllers"

	"runtime"

	"github.com/julienschmidt/httprouter"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //开启多核
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //设置静态资源的访问路径
	router := httprouter.New()
	router.GET("/", controllers.Index) //

	router.POST("/bit/req", controllers.RequestToday)
	router.GET("/win/notify", controllers.WinRequest)

	http.ListenAndServe(":9090", router)
}
