package registry

import (
	"context"
	"fmt"
	"sync"
)

// 插件管理

// 管理者结构体
type PluginMgr struct {
	// map 维护所有的插件
	plugins map[string]Registry
	lock sync.Mutex
}

var (
	// 初始化
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

// 插件注册
func RegisterPlugin(registry Registry) (err error) {
	return pluginMgr.registerPlugin(registry)
}
// 注册插件
func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// 先去看里面有没有
	_ , ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("registry plugin exist")
		return
	}
	p.plugins[plugin.Name()] = plugin
	return
}

// 进行初始化注册中心
func InitRegistry(ctx context.Context, name string,opts ...Option) (registry Registry ,err error){
	return pluginMgr.initRegistry(ctx,name,opts ...)
}
func (p *PluginMgr)initRegistry(ctx context.Context, name string,opts ...Option) (registry Registry ,err error){
	p.lock.Lock()
	defer p.lock.Unlock()
	// 先查看服务列表，服务是否存在，若不存在，就没有初始化，报错
	plugin , ok := p.plugins[name]
	if !ok{
		err = fmt.Errorf("plugin %s not exist",name)
		return
	}
	// 存在， 返回值赋值
	registry = plugin
	// 进行组件初始化
	err = plugin.Init(ctx, opts...)
	return
}