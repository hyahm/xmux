package xmux

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

func exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	var send []byte
	var err error
	if GetInstance(r).Response != nil && GetInstance(r).Get(STATUSCODE).(int) == 200 {

		ck := GetInstance(r).Get(CacheKey)

		if ck != nil {
			cacheKey := ck.(string)
			if IsUpdate(cacheKey) {
				// 如果没有设置缓存，还是以前的处理方法
				send, err = json.Marshal(GetInstance(r).Response)
				if err != nil {
					log.Println(err)
				}

				// 如果之前是更新的状态，那么就修改
				SetCache(cacheKey, send)
				w.Write(send)
			} else {
				// 如果不是更新的状态， 那么就不用更新，而是直接从缓存取值
				send = GetCache(cacheKey)
				w.Write(send)
			}
		} else {
			send, err = json.Marshal(GetInstance(r).Response)
			if err != nil {
				log.Println(err)
			}
			w.Write(send)
		}

	}
	log.Printf("connect_id: %d,method: %s\turl: %s\ttime: %f\t status_code: %v, body: %v\n",
		GetInstance(r).GetConnectId(),
		r.Method,
		r.URL.Path, time.Since(start).Seconds(),
		GetInstance(r).Get(STATUSCODE),
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
	mu        sync.RWMutex
}

// 获取module数组
func (m *module) cloneMudule() *module {
	newModule := &module{
		filter:    make(map[string]struct{}),
		funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0, len(m.funcOrder)),
		mu:        sync.RWMutex{},
	}
	for k := range m.filter {
		newModule.filter[k] = struct{}{}
	}
	newModule.funcOrder = append(newModule.funcOrder, m.funcOrder...)
	return newModule
}

// 删除模块， 返回新的模块
func (m *module) delete(delmodules map[string]struct{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for name := range delmodules {
		if _, ok := m.filter[name]; ok {
			// 说明存在
			for index, value := range m.funcOrder {
				if GetFuncName(value) == name {
					m.funcOrder = append(m.funcOrder[:index], m.funcOrder[index+1:]...)
					delete(m.filter, name)
					break
				}
			}
		}
	}
}

func (m *module) GetModules() []func(w http.ResponseWriter, r *http.Request) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.funcOrder
}

// 添加模块
func (m *module) add(mds ...func(w http.ResponseWriter, r *http.Request) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 添加, 不会重复添加 module
	for _, md := range mds {
		name := GetFuncName(md)
		if _, ok := m.filter[name]; !ok {
			m.funcOrder = append(m.funcOrder, md)
			m.filter[name] = struct{}{}
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
