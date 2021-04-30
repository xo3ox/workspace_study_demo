package session

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"sync"
)

// 对象
type RedisSession struct {
	sessionId string
	pool      *redis.Pool // 连接redis的连接池
	// 设置session，可以先放在内存的map中，批量导入redis，提升性能，
	sessionMap map[string]interface{} // 先加载到内存中
	rwLock     sync.RWMutex           // 读写锁
	flag       int                    // 记录内存中map是否被操作
}

// 用常量定义状态
const (
	SessionFlagNone   = iota // 内存数据没有变化
	SessionFlagModify        // 有变化
)

// 构造函数
func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	s := &RedisSession{
		sessionId:  id,
		sessionMap: make(map[string]interface{}, 16),
		pool:       pool,
		flag:       SessionFlagNone,
	}
	return s
}

// 将session存储到内存的map中
func (r *RedisSession) Set(key string, value interface{}) (err error) {
	r.rwLock.Lock() // 加锁
	defer r.rwLock.Unlock()
	// 设置值
	r.sessionMap[key] = value
	// 标记记录
	r.flag = SessionFlagModify
	return
}
// 从redis中再次加载
func (r *RedisSession) loadFromRedis() (err error){
	conn := r.pool.Get()
	reply, err := conn.Do("GET", r.sessionId)
	if err != nil{
		return
	}
	// 转字符串
	data, err := redis.String(reply, err)
	if err != nil{
		return
	}
	// 取到东西，反序列化到内存map
	err = json.Unmarshal([]byte(data), &r.sessionMap)
	if err != nil{
		return
	}
	return
}
func (r *RedisSession) Get(key string) (value interface{}, err error) {
	r.rwLock.Lock() // 加锁
	defer r.rwLock.Unlock()
	// 先判断内存
	var ok bool
	value, ok = r.sessionMap[key]
	if !ok {
		err = errors.New("key not exists")
	}

	return

}
func (r *RedisSession) Del(key string) (err error) {
	r.rwLock.Lock() // 加锁
	defer r.rwLock.Unlock()
	r.flag = SessionFlagModify
	delete(r.sessionMap,key)
	return
}

// 将session存储到redis
func (r *RedisSession) Save() (err error) {
	r.rwLock.Lock() // 加锁
	defer r.rwLock.Unlock()
	// 若数据没有改变，不需要存
	if r.flag != SessionFlagModify {
		return
	}

	// 内存中的sessionMap进行序列化
	var data []byte
	data, err = json.Marshal(r.sessionMap)
	if err != nil {
		return
	}
	// 获取redis连接
	conn := r.pool.Get()
	// 保存k,v到redis
	_, err = conn.Do("SET", r.sessionId, string(data))
	r.flag = SessionFlagNone // 将session状态改回去
	if err != nil {
		return
	}

	return
}
