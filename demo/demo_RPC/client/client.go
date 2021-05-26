package main

import (
	"demo_RPC/design"
	"fmt"
)

// 结合design测试
func main(){
	myClient := design.InitMyClient("127.0.0.1:8080")
	var resp string
	err := myClient.HelloWorld("王八蛋", &resp)
	if err != nil{
		fmt.Println("HelloWorld 调用失败：",err)
		return
	}
	fmt.Println(resp,err)

}

//func main(){
//	// 1 用RPC连接服务器
//	//conn, err := rpc.Dial("tcp", "127.0.0.1:8080")
//	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8080") // 防止网络传输乱码问题
//	if err != nil {
//		fmt.Println("rpc.Dial 错误：",err)
//		return
//	}
//	defer conn.Close()
//
//	// 2 调用远程函数
//	var replay string	// 接收调用函数的返回值,通过传输参数
//	err = conn.Call("hello.HelloWorld", "大帅哥", &replay)
//	if err != nil {
//		fmt.Println("conn.Call 错误：",err)
//		return
//	}
//	fmt.Println(replay)
//}


