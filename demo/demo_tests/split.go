package demo_tests

import "strings"

// 根据传入的sep 分割字符串
func Split(s ,sep string) (res []string){
	i := strings.Index(s, sep)	// 获取sep在s中的位置
	for i >= 0{
		res = append(res , s[:i])	// 根据位置追加到res
		s = s[i+len(sep) : ]		// 从len(sep)的位置开始继续切分字符串
		i = strings.Index(s, sep)	// 重置i的位置
	}
	res = append(res, s)	// 最后再追加上最后一个切割的字符串
	return
}
