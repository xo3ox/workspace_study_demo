package main

// 测试写的日志库

import (
	"rizhiku/mylogger"
)

var log mylogger.Logger

func main(){
	log = mylogger.NewConsoleLogger("Info")
	log = mylogger.NewFileLogger("debug","./mylogger/","log.log",10*1024*1024)
	for  {
		id := 1001
		name := "wangbad"
		log.Debug("这是一条Debug级别的日志,id:%d,name:%s",id,name)
		log.Trace("这是一条Trace级别的日志")
		log.Info("这是一条Info级别的日志")
		log.Warning("这是一条Waring级别的日志")
		log.Fatal("这是一条Fatal级别的日志")
		log.Error("这是一条Error级别的日志")
		//time.Sleep(1*time.Second)
	}
}