package xmux

import (
	"reflect"
	"runtime"
)

func compareFunc(f1 interface{}, f2 interface{}) bool {
	if reflect.TypeOf(f1).Kind() != reflect.Func {
		return false
	}
	if reflect.TypeOf(f2).Kind() != reflect.Func {
		return false
	}

	tm1 := reflect.ValueOf(f1).Pointer()

	tm2 := reflect.ValueOf(f2).Pointer()

	name1 := runtime.FuncForPC(tm1).Name()
	name2 := runtime.FuncForPC(tm2).Name()
	return name1 == name2
}
