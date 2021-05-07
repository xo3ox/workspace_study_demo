package registry

import "context"

type Registry interface {
	Name() string                                                              // 插件名字，例如传etcd
	Init(ctx context.Context, opts ...Option) (err error)                      // 初始化
	Register(ctx context.Context, service *Service) (err error)                // 服务注册
	UnRegister(ctx context.Context, service *Service) (err error)              // 服务反注册
	GetService(ctx context.Context, name string) (service *Service, err error) // 服务发现

}
