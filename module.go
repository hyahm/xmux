package xmux

import (
	"net/http"
	"reflect"
	"runtime"
)

func DefaultModuleTemplate(w http.ResponseWriter, r *http.Request) bool {
	return false
}

type funcOrder []string
type filter map[string]func(w http.ResponseWriter, r *http.Request) bool // 函数名和函数值

type delModule struct {
	modules map[string]struct{}
}

// 在启动的时候才会执行， 所以不用加锁
type module struct {
	// order     []func(w http.ResponseWriter, r *http.Request) bool // 保存执行的顺序
	filter    filter    // 过滤重复的
	funcOrder funcOrder // 函数名排序
}

// 获取module数组
func (m module) getMuduleList() []func(w http.ResponseWriter, r *http.Request) bool {
	ml := make([]func(w http.ResponseWriter, r *http.Request) bool, 0)
	for _, name := range m.funcOrder {
		ml = append(ml, m.filter[name])
	}
	return ml
}

func (m module) add(mfs ...func(w http.ResponseWriter, r *http.Request) bool) module {
	// 添加, 不会重复添加 module
	temp := module{
		filter:    make(filter),
		funcOrder: make(funcOrder, len(m.funcOrder)),
	}
	for k, v := range m.filter {
		temp.filter[k] = v
	}
	for k, v := range m.funcOrder {
		temp.funcOrder[k] = v
	}

	for _, mf := range mfs {
		mn := runtime.FuncForPC(reflect.ValueOf(mf).Pointer()).Name()
		if _, ok := temp.filter[mn]; !ok {
			temp.funcOrder = append(temp.funcOrder, mn)
			temp.filter[mn] = mf
		}
	}
	return temp
}

func (m module) addModule(new module) module {

	// 添加, module
	temp := module{
		filter:    make(filter),
		funcOrder: make(funcOrder, len(m.funcOrder)),
	}
	for k, v := range m.filter {
		temp.filter[k] = v
	}
	for k, v := range m.funcOrder {
		temp.funcOrder[k] = v
	}
	if temp.filter == nil {
		temp.filter = make(filter)
	}
	for _, name := range new.funcOrder {
		if _, ok := temp.filter[name]; !ok {
			temp.funcOrder = append(temp.funcOrder, name)
			temp.filter[name] = new.filter[name]
		}
	}
	return temp
}

func (m module) deleteKey(names ...string) module {
	// 添加, 不会重复添加 module
	temp := module{
		filter:    make(filter),
		funcOrder: make(funcOrder, len(m.funcOrder)),
	}
	for k, v := range m.filter {
		temp.filter[k] = v
	}
	for k, v := range m.funcOrder {
		temp.funcOrder[k] = v
	}
	for _, delname := range names {
		if _, ok := temp.filter[delname]; ok {
			// 删除排序数组的值
			for index, name := range temp.funcOrder {
				if delname == name {
					temp.funcOrder = append(temp.funcOrder[:index], temp.funcOrder[index+1:]...)
					break
				}
			}
			// 删除map
			delete(temp.filter, delname)
		}
	}
	return temp
}

//
// 直接删除函数名， delmodule 里面使用的
func (dm delModule) addDeleteKey(mfs ...func(w http.ResponseWriter, r *http.Request) bool) delModule {
	temp := delModule{
		modules: make(map[string]struct{}),
	}
	for k, v := range dm.modules {
		temp.modules[k] = v
	}

	// 添加, 不会重复添加 module, 删除的时候， 是在程序运行的时候执行的，所以需要加锁

	for _, mf := range mfs {
		mn := runtime.FuncForPC(reflect.ValueOf(mf).Pointer()).Name()
		temp.modules[mn] = struct{}{}
	}
	return temp

}
