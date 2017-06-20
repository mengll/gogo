// Tdsp project main.go
package main

import (
	"Tdsp/controllers"

	"net/http"
	"runtime"

	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //开启多核

}

func main() {
	router := httprouter.New()
	router.GET("/index", controllers.Index)
	router.POST("/bit/req", controllers.TodayBidRequest)
	router.GET("/win/notify", controllers.WinRequest)
	router.GET("/click/notify", controllers.ClickRequest)
	router.GET("/show/notify", controllers.ShowRequest)

	http.ListenAndServe(":9090", router)

}
