package xmux

import "reflect"

func Clone(src interface{}) interface{} {
	g := reflect.TypeOf(src)
	if g.Kind() != reflect.Ptr {
		panic("bindReponse must be a pointer")
	}
	if g.Elem().Kind() != reflect.Struct {
		panic("bindReponse must be a pointer struct")
	}
	gv := reflect.ValueOf(src).Elem()
	v := reflect.New(g.Elem())
	for i := 0; i < gv.NumField(); i++ {
		v.Elem().Field(i).Set(gv.Field(i))
	}
	return v.Interface()
}
