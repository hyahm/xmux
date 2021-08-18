package xmux

import (
	"unsafe"
)

// func compareFunc(f1 interface{}, f2 interface{}) bool {
// 	if reflect.TypeOf(f1).Kind() != reflect.Func {
// 		return false
// 	}
// 	if reflect.TypeOf(f2).Kind() != reflect.Func {
// 		return false
// 	}

// 	tm1 := reflect.ValueOf(f1).Pointer()

// 	tm2 := reflect.ValueOf(f2).Pointer()

// 	name1 := runtime.FuncForPC(tm1).Name()
// 	name2 := runtime.FuncForPC(tm2).Name()
// 	return name1 == name2
// }

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
