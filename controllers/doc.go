// controllers project doc.go

/*
controllers document
*/
package controllers

import (
	"os"

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
