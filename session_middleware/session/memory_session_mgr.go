package session

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"sync"
)

// 对象
type MemorySessionMgr struct {
	sessionMap map[string]Session // 存session的map
	rwLock     sync.RWMutex       // 读写锁
}

// 构造函数
func NewMemorySessionMgr() SessionMgr{
	sr := &MemorySessionMgr{
		sessionMap: make(map[string]Session , 1024),
	}
	return sr
}

func (m *MemorySessionMgr) Init(addr string, options ...string) (err error){
	return
}   // 初始化

func (m *MemorySessionMgr) CreateSession() (session Session, err error)  {
	// 加锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	// 用uuid做session的id
	uuid := uuid.NewV4()
	// 将uuid转成string类型
	sessionId := uuid.String()
	// 创建单个session
	session = NewMemorySession(sessionId)
	// 将单个session添加到MemorySessionMgr里的map
	m.sessionMap[sessionId] = session

	return
}     // 创建单独的session对象

func (m *MemorySessionMgr) Get(sessionId string) (session Session, err error){
	// 加锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	var ok bool
	session , ok = m.sessionMap[sessionId]
	if !ok {
		err = errors.New("session not exists.")
		return
	}
	return
} // 根据sessionId获取session