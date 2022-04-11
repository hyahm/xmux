package xmux

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// 初始化临时使用， 最后会合并到 router
type Route struct {
	// 组里面也包括路由 后面的其实还是patter和handle, 还没到handle， 这里的key是个method
	isRoot bool         // 是否是router直接挂载的路由
	handle http.Handler // handle

	module    module    // 增加的 modules
	delmodule delModule // 删除的modules

	pagekeys    map[string]struct{} // 页面权限
	delPageKeys []string            // 删除的权限

	header    map[string]string // 请求头
	delheader []string          // 删除的请求头

	describe                         string // 接口描述
	request                          string // 请求的请求示例
	response                         string
	responseData                     interface{}       // 接口返回示例
	st_request                       interface{}       // api 请求示例
	params_request                   map[string]string // get请求参数
	st_response                      interface{}       // api 返回示例
	reqHeader                        map[string]string // api请求头
	supplement                       string            // api附录
	codeMsg                          map[string]string // api请求返回信息
	codeField                        string            // api 文档请求字段
	groupKey, groupLabel, groupTitle string            // 组路由的key， label， title
	apiDelReqHeader                  []string

	bindType   bindType    // 数据绑定格式
	dataSource interface{} // 数据源
	// perms map[int]
	midware    func(w http.ResponseWriter, r *http.Request)
	delmidware func(w http.ResponseWriter, r *http.Request)
}

func (rt *Route) GetHeader() map[string]string {
	return rt.header
}

func (rt *Route) BindResponse(response interface{}) *Route {
	rt.responseData = response
	return rt
}

func (rt *Route) GetMidwareName() string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	return GetFuncName(rt.midware)
}

func (rt *Route) ApiExitGroup() *Route {
	// 退出文档的组
	rt.codeField = ""
	return rt
}

func (rt *Route) AddPageKeys(pagekeys ...string) *Route {
	// 退出文档的组
	for _, v := range pagekeys {
		if rt.pagekeys == nil {
			rt.pagekeys = make(map[string]struct{})
		}
		rt.pagekeys[v] = struct{}{}
	}
	return rt
}

func (rt *Route) DelPageKeys(pagekeys ...string) *Route {
	if rt.delPageKeys == nil {
		if len(pagekeys) == 0 {
			return rt
		} else {
			rt.delPageKeys = pagekeys
		}

	}
	rt.delPageKeys = append(rt.delPageKeys, pagekeys...)
	return rt
}

// func (rt *Route) ApiAddGroup(key string) *Route {
// 	// 退出文档的组
// 	rt.groupKey = key
// 	return rt
// }

func (rt *Route) ApiCreateGroup(key, title, lable string) *Route {
	// 创建文档的组
	rt.groupKey = key
	rt.groupLabel = lable
	rt.groupTitle = title

	return rt
}

func (rt *Route) MiddleWare(midware func(w http.ResponseWriter, r *http.Request)) *Route {
	// 创建文档的组
	rt.midware = midware
	return rt
}

func (rt *Route) DelMiddleWare(midware func(w http.ResponseWriter, r *http.Request)) *Route {
	// 创建文档的组
	if rt.isRoot {
		// 那么直接就删除
		if rt.midware == nil {
			return rt
		}
		if rt.midware != nil {
			if midware != nil && GetFuncName(rt.midware) == GetFuncName(midware) {
				rt.midware = nil
			}
		}
		return rt
	} else {
		rt.delmidware = midware
	}

	return rt
}

func (rt *Route) ApiCodeField(s string) *Route {
	// 文档的 错误码字段的 key

	rt.codeField = s
	return rt
}

func (rt *Route) ApiCodeMsg(code string, msg string) *Route {
	// 文档的 错误码值及其含义
	//

	if rt.codeMsg == nil {
		rt.codeMsg = make(map[string]string)
	}
	rt.codeMsg[code] = msg
	return rt
}

// 数据绑定
func (rt *Route) Bind(s interface{}) *Route {
	rt.dataSource = s
	return rt
}

// json数据绑定
func (rt *Route) BindJson(s interface{}) *Route {
	// 接口补充说明
	rt.dataSource = s
	rt.bindType = jsonT
	return rt
}

// yaml数据绑定
func (rt *Route) BindYaml(s interface{}) *Route {
	// 接口补充说明
	rt.dataSource = s
	rt.bindType = yamlT
	return rt
}

// xml数据绑定
func (rt *Route) BindXml(s interface{}) *Route {
	// 接口补充说明

	rt.dataSource = s
	rt.bindType = xmlT
	return rt
}

func (rt *Route) ApiSupplement(s string) *Route {
	// 接口补充说明

	rt.supplement = s
	return rt
}

func (rt *Route) ApiReqStruct(s interface{}) *Route {
	// 接口返回数据的结构

	rt.st_request = s
	return rt
}

func (rt *Route) ApiReqParams(k, v string) *Route {
	// 接口返回数据的结构
	if rt.params_request == nil {
		rt.params_request = make(map[string]string)
	}
	rt.params_request[k] = v
	return rt
}

func (rt *Route) ApiResStruct(s interface{}) *Route {
	// 接口接收数据的结构

	rt.st_response = s
	return rt
}

func (rt *Route) makeDoc(url string, count *int, doc *Document) {
	// 生成侧边栏
	if rt.groupKey != "" {
		// 组路由
		// 判断key 是否存在

		if id, ok := keys[rt.groupKey]; ok {
			// 存在的话
			// 添加文档就好了
			d := apiDocument[id]
			d.Api = append(d.Api, *doc)
			// apiDocument[id].Api = append(apiDocument[id].Api, *doc)
			apiDocument[id] = d

		} else {
			keys[rt.groupKey] = *count
			d := Doc{
				Title: rt.groupTitle,
				Api:   make([]Document, 0),
			}
			d.Api = append(d.Api, *doc)
			apiDocument[*count] = d

			sideUrl := fmt.Sprintf("/-/api/%d.html", *count)
			sidebar[sideUrl] = rt.groupLabel
			*count++
		}

	}

}

func (rt *Route) ApiDescribe(s string) *Route {
	// 接口的简单描述

	rt.describe = s
	return rt
}

func (rt *Route) ApiReqHeader(k, v string) *Route {
	// 接口的请求头
	if rt.reqHeader == nil {
		rt.reqHeader = make(map[string]string)
	}
	rt.reqHeader[k] = v
	return rt
}

func (rt *Route) ApiDelReqHeader(k string) *Route {
	// 接口的请求头

	if rt.apiDelReqHeader == nil {
		rt.apiDelReqHeader = make([]string, 0)
	}
	rt.apiDelReqHeader = append(rt.apiDelReqHeader, k)
	return rt
}

func (rt *Route) ApiRequestTemplate(s JsonStr) *Route {
	// 接口的请求实例， 一般是json的字符串
	rt.request = s.Json()
	return rt
}

func (rt *Route) ApiResponseTemplate(s JsonStr) *Route {
	// 接口的返回实例， 一般是json的字符串

	rt.response = s.Json()
	return rt
}

// 组里面也包括路由 后面的其实还是patter和handle

func (rt *Route) SetHeader(k, v string) *Route {
	if rt.header == nil {
		rt.header = map[string]string{}
	}
	rt.header[k] = v
	return rt
}

func (rt *Route) DelHeader(k string) *Route {
	if rt.isRoot {
		for k := range rt.header {
			delete(rt.header, k)
		}
	} else {
		if rt.delheader == nil {
			rt.delheader = make([]string, 0)
		}
		rt.delheader = append(rt.delheader, k)
	}
	return rt
}

func (rt *Route) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) *Route {
	rt.module = rt.module.add(handles...)
	return rt
}

func (rt *Route) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) *Route {
	if rt.isRoot {
		for _, mf := range handles {
			mn := runtime.FuncForPC(reflect.ValueOf(mf).Pointer()).Name()
			rt.module = rt.module.deleteKey(mn)
		}
	} else {
		rt.delmodule = rt.delmodule.addDeleteKey(handles...)
	}
	return rt
}

func (rt *Route) GetModule() []string {

	return rt.module.funcOrder
}
