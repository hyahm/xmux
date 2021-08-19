package xmux

import (
	"reflect"
	"runtime"
	"unsafe"
)

// 获取函数名
func GetFuncName(f interface{}) string {
	rv := reflect.ValueOf(f)
	if rv.Type().Kind() != reflect.Func {
		panic("not a func")
	}
	return runtime.FuncForPC(rv.Pointer()).Name()
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
