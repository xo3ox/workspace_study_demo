package session

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type RedisSessionMgr struct {
	// redis 地址
	addr string
	// 密码
	pass string
	// 连接池
	pool *redis.Pool
	// 锁
	rwLock sync.RWMutex
	// 大map
	sessionMap map[string]Session
}

// 构造函数
func NewRedisSessionMgr() SessionMgr {
	sr := &RedisSessionMgr{
		sessionMap: make(map[string]Session, 32),
	}
	return sr
}

// 创建连接池
func myPool(addr, pass string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     64,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			// 如果有密码，判断
			if _, err := conn.Do("AUTH", pass); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, err
		},
		// 连接测试,开发写，上线注释
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
func (r *RedisSessionMgr) Init(addr string, options ...string) (err error) {
	// 若有其他参数
	if len(options) > 0 {
		r.pass = options[0]
	}
	// 创建连接池
	r.pool = myPool(addr, r.pass)
	r.addr = addr
	return
} // 初始化

func (r *RedisSessionMgr) CreateSession() (session Session, err error) {
	// 加锁
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	// 用uuid做session的id
	uuid := uuid.NewV4()
	// 将uuid转成string类型
	sessionId := uuid.String()
	// 创建单个session
	session = NewRedisSession(sessionId, r.pool)
	// 将单个session添加到MemorySessionMgr里的map
	r.sessionMap[sessionId] = session

	return
} // 创建单独的session对象

func (r *RedisSessionMgr) Get(sessionId string) (session Session, err error) {
	// 加锁
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	var ok bool
	session , ok = r.sessionMap[sessionId]
	if !ok {
		err = errors.New("session not exists.")
		return
	}
	return
} // 根据sessionId获取session
