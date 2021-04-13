package demo_tests

import (
	"reflect"
	"testing"
)
// 测试组
func TestSplit(t *testing.T) {
	// 定义一个测试用例类型
	type test struct {
		input string
		sep   string
		want  []string
	}
	// 定义一个存储测试用例的切片
	tests := []test{
		{input: "a:b:c", sep: ":", want: []string{"a", "b", "c"}},
		{input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		{input: "abcd", sep: "bc", want: []string{"a", "d"}},

		{input: "aaazaaa", sep: "z", want: []string{"aaa", "aaa"}},
	}
	// 遍历切片，逐一执行测试用例
	for _, tc := range tests {
		got := Split(tc.input, tc.sep)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("excepted:%v, got:%v", tc.want, got)
		}
	}
}

func Test1Split(t *testing.T) {
	got := Split("abac,asdk,aloa,lkj",",")	// 程序输出的结果
	want := []string{"abac","asdk","aloa","lkj"}	// 期望的结果
	if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("excepted:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

func Test2Split(t *testing.T) {
	got := Split("abac,acsdk,acloa,lkacj","ac")	// 程序输出的结果
	want := []string{"ab",",","sdk,","loa,lk","j"}	// 期望的结果
	if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("excepted:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}