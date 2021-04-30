package session

import (
	"errors"
	"sync"
)

// 对象
type MemorySession struct {
	sessionId string
	data map[string]interface{} // 存key、value
	rwLock sync.RWMutex	// 读写锁
}

// 构造函数
func NewMemorySession(id string) *MemorySession{
	s := &MemorySession{
		sessionId: id,
		data: make(map[string]interface{} , 16),
	}
	return s
}

// 实现interface接口定义的方法
func (m *MemorySession) Set(key string, value interface{}) (err error)  {
	// 对session的map写入操作，首先加锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	// 设置值
	m.data[key]=value
	return
}

func (m *MemorySession) Get(key string) (value interface{}, err error)  {
	// 对session的map获取操作，首先加锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	// 根据key取值
	var ok bool
	value , ok  = m.data[key]
	if !ok{
		err = errors.New("key not exists in session.")
		return
	}
	return
}

func (m *MemorySession) Del(key string)(err error)  {
	// 对session的map删除操作，首先加锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.data,key) // map的删除
	return
}

func (m *MemorySession) Save() (err error) {
	return
}