// session 的管理者
package session

// 定义管理者管理所有session
type SessionMgr interface {
	Init(addr string, options ...string) (err error)   // 初始化
	CreateSession() (session Session, err error)       // 创建单独的session对象
	Get(sessionId string) (session Session, err error) // 根据sessionId获取session
}
