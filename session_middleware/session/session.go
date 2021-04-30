// 定义一些规范接口
package session

// session接口
type Session interface {
	Set(key string, value interface{}) (err error) // 设置session
	Get(key string) (value interface{},err error)     // 获取session
	Del(key string)(err error)                   // 删除session
	Save() (err error)                            // 保存
}
