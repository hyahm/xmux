package xmux

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	var send []byte
	var err error
	if GetInstance(r).Response != nil {
		send, err = json.Marshal(GetInstance(r).Response)
		if err != nil {
			log.Println(err)
		}
		w.Write(send)
	}
	log.Printf("connect_id: %d,method: %s\turl: %s\ttime: %f\t status_code: %v, body: %v\n",
		GetInstance(r).Get(CONNECTID),
		r.Method,
		r.URL.Path, time.Since(start).Seconds(), GetInstance(r).Get(STATUSCODE),
		string(send))
}

func DefaultModuleTemplate(w http.ResponseWriter, r *http.Request) bool {
	return false
}

// 在启动的时候才会执行， 所以不用加锁
type module struct {
	// order     []func(w http.ResponseWriter, r *http.Request) bool // 保存执行的顺序
	filter    map[string]struct{}                                 // 过滤重复的,值是索引的位置
	funcOrder []func(w http.ResponseWriter, r *http.Request) bool // 函数名排序
}

// 获取module数组
func (m *module) cloneMudule() *module {
	if m == nil {
		return &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		}

	}
	return &module{
		filter:    m.filter,
		funcOrder: m.funcOrder,
	}
}

// 删除模块， 返回新的模块
func (m *module) delete(delmodules map[string]struct{}) {
	if m == nil {
		m = &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		}
	}
	for name := range delmodules {
		if _, ok := m.filter[name]; ok {
			// 说明存在
			exsit := false
			for index, value := range m.funcOrder {
				if GetFuncName(value) == name {
					m.funcOrder = append(m.funcOrder[:index], m.funcOrder[index+1:]...)
					delete(m.filter, name)
					break
				}
			}
			if !exsit {
				log.Println("xmux must be have somthing wrong, please make issue it to https://github.com/hyahm/xmux/issues")
			}
		}
	}
}

func (m *module) GetModules() []func(w http.ResponseWriter, r *http.Request) bool {
	return m.funcOrder
}

// 添加模块
func (m *module) add(mds ...func(w http.ResponseWriter, r *http.Request) bool) {
	// 添加, 不会重复添加 module
	if m == nil {
		m = &module{
			filter:    make(map[string]struct{}),
			funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		}
	}
	for _, md := range mds {
		if _, ok := m.filter[GetFuncName(md)]; !ok {
			m.funcOrder = append(m.funcOrder, md)
		}
	}

}

// func (m module) addModule(new module) module {
// 	// 添加, module
// 	temp := module{
// 		filter:    make(filter),
// 		funcOrder: make(funcOrder, len(m.funcOrder)),
// 	}
// 	for k, v := range m.filter {
// 		temp.filter[k] = v
// 	}
// 	for k, v := range m.funcOrder {
// 		temp.funcOrder[k] = v
// 	}
// 	if temp.filter == nil {
// 		temp.filter = make(filter)
// 	}
// 	for _, name := range new.funcOrder {
// 		if _, ok := temp.filter[name]; !ok {
// 			temp.funcOrder = append(temp.funcOrder, name)
// 			temp.filter[name] = new.filter[name]
// 		}
// 	}
// 	return temp
// }

// func (m module) deleteKey(names ...string) module {
// 	// 添加, 不会重复添加 module
// 	temp := module{
// 		filter:    make(filter),
// 		funcOrder: make(funcOrder, len(m.funcOrder)),
// 	}
// 	for k, v := range m.filter {
// 		temp.filter[k] = v
// 	}
// 	for k, v := range m.funcOrder {
// 		temp.funcOrder[k] = v
// 	}
// 	for _, delname := range names {
// 		if _, ok := temp.filter[delname]; ok {
// 			// 删除排序数组的值
// 			for index, name := range temp.funcOrder {
// 				if delname == name {
// 					temp.funcOrder = append(temp.funcOrder[:index], temp.funcOrder[index+1:]...)
// 					break
// 				}
// 			}
// 			// 删除map
// 			delete(temp.filter, delname)
// 		}
// 	}
// 	return temp
// }

//
// 直接删除函数名， delmodule 里面使用的
// func (dm delModule) addDeleteKey(mfs ...func(w http.ResponseWriter, r *http.Request) bool) delModule {
// 	temp := delModule{
// 		modules: make(map[string]struct{}),
// 	}
// 	for k, v := range dm.modules {
// 		temp.modules[k] = v
// 	}

// 	// 添加, 不会重复添加 module, 删除的时候， 是在程序运行的时候执行的，所以需要加锁

// 	for _, mf := range mfs {
// 		mn := runtime.FuncForPC(reflect.ValueOf(mf).Pointer()).Name()
// 		temp.modules[mn] = struct{}{}
// 	}
// 	return temp

// }
