package helper

import (
	"reflect"
	"runtime"
)

// 获取函数名
func GetFuncName(f interface{}) string {
	rv := reflect.ValueOf(f)
	if rv.Type().Kind() != reflect.Func {
		panic("not a func")
	}
	return runtime.FuncForPC(rv.Pointer()).Name()
}
