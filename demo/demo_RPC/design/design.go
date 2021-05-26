package design

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 要求服务端在注册RPC对象时，能让编译器检测出 注册对象是否合法

//-----------服务端------------
// 创建接口，在接口中定义方法的原型
type MyInterface interface {
	HelloWorld(string, *string) error
}

// 注册服务 , 调用该方法时，需要给 i 传参，参数应该是实现了HelloWorld接口方法的对象
func RegisterService(i MyInterface){
	rpc.RegisterName("hello",i)
}

// ---------客户端-----------
// 像调用本地函数一样，调用远程函数

// 定义类
type MyClient struct {
	c *rpc.Client
}

// 由于使用了c调用Call方法，因此需要初始化c
func InitMyClient(addr string)MyClient{
	conn , _ := jsonrpc.Dial("tcp",addr)
	return MyClient{
		c:conn,
	}
}

// 实现函数， 原型参照上面的interface接口来实现
func (m *MyClient) HelloWorld(a string , b *string) error{

	// 参数1：参照上面的interface ，RegisterName 而来  ， a : 传入参数，b : 传出参数
	return m.c.Call("hello.HelloWorld",a,b)
}