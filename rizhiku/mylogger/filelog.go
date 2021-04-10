package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件中写日志相关代码

type FileLogger struct {
	Level LogLevel
	filePath string	// 日志文件保存的路径
	fileName string	// 日志文件保存的文件名
	fileObj *os.File
	errFileObj *os.File
	maxFileSize int64	// 最大文件大小
}

// 构造函数
func NewFileLogger(levelStr,fp,fn string ,maxSize int64) *FileLogger{
	logLevel , err := parseLogLevel(levelStr)
	if err != nil{
		panic(err)
	}
	fl := &FileLogger{
		Level: logLevel,
		filePath: fp,
		fileName: fn,
		maxFileSize: maxSize,

	}
	err = fl.initFile()	// 按照文件路径和文件名将文件打开
	if err !=  nil{
		panic(err)
	}
	return fl
}
// 方法--------------------------------
// 根据指定的日志文件路径和文件名打开日志文件
func (f *FileLogger) initFile() error{
	// 记录普通级别的日志
	fullFileName := path.Join(f.filePath,f.fileName)
	fileObj , err := os.OpenFile(fullFileName, os.O_APPEND | os.O_CREATE |os.O_WRONLY,0644)
	if err != nil{
		fmt.Printf("打开日志文件出错，err : %v\n", err)
		return err
	}
	// 记录错误级别的日志
	errFileObj , err := os.OpenFile(fullFileName+".err", os.O_APPEND | os.O_CREATE |os.O_WRONLY,0644)
	if err != nil{
		fmt.Printf("打开错误日志文件出错，err : %v\n", err)
		return err
	}
	// 日志文件都打开了
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

// 判断是否需要记录该日志
func (f *FileLogger) enable(level LogLevel) bool {
	return  level >= f.Level	// 当传入的日志等级大于等于日志等级时返回true
}


// 根据文件大小判断是否需要切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo ,err := file.Stat()
	if err != nil{
		fmt.Printf("获取文件信息失败：err : %v\n",err)
		return false
	}
	// 如果当前文件大小 大于等于 日志文件的最大值 就应该返回true
	return 	fileInfo.Size() >= f.maxFileSize
}

// 切割文件
func (f *FileLogger)splitFile(file *os.File) (*os.File,error) {
	// 需要切割文件
	// 1 备份一下 rename
	nowStr := time.Now().Format("2006_01_02_15_04_05")
	fileInfo , err := file.Stat()
	if err != nil{
		fmt.Printf("获取文件信息失败，err : %v\n",err)
		return nil, err
	}

	logName := path.Join(f.filePath,fileInfo.Name())	// 拿到当前的日志文件完整路径
	newLogName := fmt.Sprintf("%s.bak%s",logName,nowStr)	// 拼接一个日志文件的备份名
	// 2 关闭当前日志文件
	file.Close()
	os.Rename(logName,newLogName)
	// 3 打开一个新的日志文件
	fileObj ,err := os.OpenFile(logName,os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)
	if err != nil{
		fmt.Printf("打开新的日志文件失败，err : %v\n",err)
		return nil,err
	}
	// 4 将打开的新的日志文件对象赋值给 f.fileObj
	return fileObj, nil
}

// 记录日志的方法
func (f *FileLogger) log(lv LogLevel , format string ,a ...interface{}){
	if f.enable(lv) {
		//now := time.Now().Format("2006-01-02T15:04:05.000+0800")
		now := time.Now().Format("2006年01月02日 15:04:05")

		funcName, fileName, lineNo := getInfo(3)

		msg := fmt.Sprintf(format, a...)

		if f.checkSize(f.fileObj) {
			newFile , err := f.splitFile(f.fileObj)	// 普通日志文件
			if err != nil {
				return
			}
			f.fileObj = newFile
		}
		_ ,_ =fmt.Fprintf(f.fileObj,"[%s] [%s] [%s:%s:%d]  %s\n",
			now, getLogString(lv), funcName, fileName, lineNo, msg)
		if lv >= ERROR{ // 如果要记录的日志级别大于等于Error级别，还要在err日志文件中再记录一遍
			if f.checkSize(f.errFileObj) {
				newFile , err := f.splitFile(f.errFileObj)	// 普通日志文件
				if err != nil {
					return
				}
				f.errFileObj = newFile
			}
			fmt.Fprintf(f.errFileObj,"[%s] [%s] [%s:%s:%d]  %s\n",
				now, getLogString(lv), funcName, fileName, lineNo, msg)
		}
	}
}


func (f *FileLogger) Close(){
	f.fileObj.Close()
	f.errFileObj.Close()
}

func (f *FileLogger) Debug(format string, a ...interface{})  {

		f.log(DEBUG,format, a...)

}

func (f *FileLogger) Trace(format string, a ...interface{})  {

		f.log(TRACE,format)

}

func (f *FileLogger) Info(format string, a ...interface{})  {

	f.log(INFO,format)

}

func (f *FileLogger) Warning(format string, a ...interface{})  {

	f.log(WARNING,format)

}

func (f *FileLogger) Error(format string, a ...interface{})  {

	f.log(ERROR,format)

}

func (f *FileLogger) Fatal(format string, a ...interface{})  {

	f.log(FATAL,format)

}



