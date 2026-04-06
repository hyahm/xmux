package xmux

import (
	"reflect"
	"time"
)

// func Clone(src interface{}) interface{} {
// 	if src == nil {
// 		return nil
// 	}
// 	g := reflect.TypeOf(src)
// 	if g.Kind() != reflect.Ptr {
// 		panic("bindReponse must be a pointer")
// 	}
// 	if g.Elem().Kind() != reflect.Struct {
// 		panic("bindReponse must be a pointer struct")
// 	}
// 	gv := reflect.ValueOf(src).Elem()
// 	v := reflect.New(g.Elem())
// 	for i := 0; i < gv.NumField(); i++ {
// 		v.Elem().Field(i).Set(gv.Field(i))
// 	}
// 	return v.Interface()
// }

// DeepCopy 真正的递归深拷贝，支持任意类型
func DeepCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	original := reflect.ValueOf(src)
	return deepCopyRecursive(original).Interface()
}

func deepCopyRecursive(original reflect.Value) reflect.Value {
	// 处理可导出性判断
	if !original.CanSet() {
		return original
	}

	// 处理各种类型
	switch original.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		// 值类型直接复制
		return original

	case reflect.Ptr:
		if original.IsNil() {
			return reflect.Zero(original.Type())
		}
		// 递归拷贝指针指向的值
		elemCopy := deepCopyRecursive(original.Elem())
		copyPtr := reflect.New(elemCopy.Type())
		copyPtr.Elem().Set(elemCopy)
		return copyPtr

	case reflect.Struct:
		// 特殊处理 time.Time（不可导出字段，直接值复制）
		if original.Type() == reflect.TypeOf(time.Time{}) {
			return original
		}
		// 结构体字段递归拷贝
		copyStruct := reflect.New(original.Type()).Elem()
		for i := 0; i < original.NumField(); i++ {
			field := original.Field(i)
			if !field.CanSet() {
				continue
			}
			copyField := deepCopyRecursive(field)
			copyStruct.Field(i).Set(copyField)
		}
		return copyStruct

	case reflect.Slice:
		if original.IsNil() {
			return reflect.Zero(original.Type())
		}
		// 新建切片，逐个元素深拷贝
		copySlice := reflect.MakeSlice(original.Type(), original.Len(), original.Cap())
		for i := 0; i < original.Len(); i++ {
			elem := deepCopyRecursive(original.Index(i))
			copySlice.Index(i).Set(elem)
		}
		return copySlice

	case reflect.Map:
		if original.IsNil() {
			return reflect.Zero(original.Type())
		}
		// 新建 map，逐个 key/value 深拷贝
		copyMap := reflect.MakeMap(original.Type())
		for _, key := range original.MapKeys() {
			k := deepCopyRecursive(key)
			v := deepCopyRecursive(original.MapIndex(key))
			copyMap.SetMapIndex(k, v)
		}
		return copyMap

	case reflect.Array:
		copyArray := reflect.New(original.Type()).Elem()
		for i := 0; i < original.Len(); i++ {
			elem := deepCopyRecursive(original.Index(i))
			copyArray.Index(i).Set(elem)
		}
		return copyArray

	default:
		// 函数、通道等不支持拷贝，直接返回原对象
		return original
	}
}
