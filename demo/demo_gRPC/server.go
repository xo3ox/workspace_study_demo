package main

// 1 需要监听
// 2 需要实例化gRPC服务端
// 3 在gRPC上注册微服务
// 4 启动服务端

import (
	"context"
	pb "demo_gRPC/proto"
	"fmt"
	"google.golang.org/grpc"

	"net"
)

// 定义空结构体，对应.proto文件中的接口message
type UserInfoService struct {
}
var u = UserInfoService{}
// 实现方法，对应.proto文件中的方法service
func (u *UserInfoService)GetUserInfo(ctx context.Context,req *pb.UserRequest)(resp *pb.UserResponse,err error){
	// 通过用户名查询用户信息
	name := req.Name
	// 数据里查用户信息
	if name == "zs"{
		resp = &pb.UserResponse{
			Id: 1,
			Name: name,
			Age: 24,
			Hobby: []string{"sing","song","run"},
		}
	}
	return
}


func main(){
	// 地址
	addr := "127.0.0.1:8080"
	// 1. 监听
	listen, err := net.Listen("tcp", addr)
	if err != nil{
		fmt.Println("监听异常；",err)
	}
	fmt.Println("监听端口：",addr)
	// 2. 实例化gRPC
	s := grpc.NewServer()
	// 3. 在gRPC上注册微服务
	pb.RegisterUserInfoServiceServer(s,&u)
	// 4. 启动服务端
	s.Serve(listen)

}