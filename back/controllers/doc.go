// controllers project doc.go

/*
controllers document
*/
package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	Logrus "github.com/sirupsen/logrus"
)

var Log *Logrus.Logger

func init() {
	Log = Logrus.New()
	//Log.SetFormatter(&Log.JSONFormatter{})
	file, err := os.OpenFile("logrus.log", os.O_APPEND|os.O_CREATE, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

}

func Httprequest(requestUrl, requestType, requestData string) (string, bool) {
	defer func() {
		if err := recover(); err != nil {
			//Logdebug("error", JsonEncodeString(err))
		}
	}()

	client := new(http.Client)
	reqest, err := http.NewRequest(requestType, requestUrl, strings.NewReader(requestData))

	if err != nil {
		//Logdebug("error", err.Error())
		return "网络请求出错", false
	}

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	resp, err := client.Do(reqest)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "err", false
	}
	bodyText := string(body)
	status := resp.StatusCode
	backContent := fmt.Sprintf("请求状态：%d 请求的响应时间: %s 请求响应的页面内容：%s", status, resp.Header.Get("Date"), bodyText)
	//	Logdebug("20170809", backContent)
	fmt.Println(backContent)
	if status == 200 {
		return bodyText, true
	} else {
		return bodyText, false
	}

}

//func Logfunc(str string, data interface{}) {
//	Log.WithFields(log.Fields{
//		"animal": "walrus",
//		"size":   10,
//	}).Info("A group of walrus emerges from the ocean")
//}

//日志的使用的案例
//  log.WithFields(log.Fields{
//    "animal": "walrus",
//    "size":   10,
//  }).Info("A group of walrus emerges from the ocean")

//  log.WithFields(log.Fields{
//    "omg":    true,
//    "number": 122,
//  }).Warn("The group's number increased tremendously!")

//  log.WithFields(log.Fields{
//    "omg":    true,
//    "number": 100,
//  }).Fatal("The ice breaks!")

//  // A common pattern is to re-use fields between logging statements by re-using
//  // the logrus.Entry returned from WithFields()
//  contextLogger := log.WithFields(log.Fields{
//    "common": "this is a common field",
//    "other": "I also should be logged always",
//  })
