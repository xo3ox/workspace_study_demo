// 向终端写日志相关内容
package mylogger

import (
	"fmt"
	"time"
)

// ConsoleLogger 日志类型的结构体
type ConsoleLogger struct {
	Level LogLevel
}

// NewLog 构造函数
func NewConsoleLogger(leveStr string) ConsoleLogger{
	level , err := parseLogLevel(leveStr)
	if err != nil{
		panic(err)
	}
	return ConsoleLogger{
		Level: level,
	}
}


// 方法--------------------------------

// 判断日志级别
func (c ConsoleLogger) enable(level LogLevel) bool {
	return  level >= c.Level	// 当传入的日志等级大于等于日志等级时返回true
}

// 写日志
func (c ConsoleLogger)log(lv LogLevel , format string ,a ...interface{}){
	if c.enable(lv) {
		//now := time.Now().Format("2006-01-02T15:04:05.000+0800")
		now := time.Now().Format("2006年01月02日 15:04:05")

		funcName, fileName, lineNo := getInfo(3)

		msg := fmt.Sprintf(format, a...)

		fmt.Printf("[%s] [%s] [%s:%s:%d]  %s\n",
			now, getLogString(lv), funcName, fileName, lineNo, msg)
	}
}


func (c ConsoleLogger) Debug(format string, a ...interface{})  {
		c.log(DEBUG,format, a...)
}

func (c ConsoleLogger) Trace(format string, a ...interface{})  {
		c.log(TRACE,format, a...)
}

func (c ConsoleLogger) Info(format string, a ...interface{})  {
		c.log(INFO,format, a...)
}

func (c ConsoleLogger) Warning(format string, a ...interface{})  {
		c.log(WARNING,format, a...)
}

func (c ConsoleLogger) Error(format string, a ...interface{})  {
		c.log(ERROR,format, a...)
}

func (c ConsoleLogger) Fatal(format string, a ...interface{})  {
		c.log(FATAL,format, a...)
}






