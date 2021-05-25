package main

import (
	"fmt"
	"net/rpc"
)

func main(){
	// 1 用RPC连接服务器
	conn, err := rpc.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("rpc.Dial 错误：",err)
		return
	}
	defer conn.Close()

	// 2 调用远程函数
	var replay string	// 接收调用函数的返回值,通过传输参数
	err = conn.Call("hello.HelloWorld", "大帅哥", &replay)
	if err != nil {
		fmt.Println("conn.Call 错误：",err)
		return
	}
	fmt.Println(replay)
}


