package xmux

import (
	"fmt"
	"net/http"
)

// 初始化临时使用， 最后会合并到 router
type Route struct {
	// 组里面也包括路由 后面的其实还是patter和handle, 还没到handle， 这里的key是个method
	handle     http.Handler
	header     map[string]string
	args       []string // 保存正则的变量名
	midware    []func(http.ResponseWriter, *http.Request) bool
	delmidware []func(http.ResponseWriter, *http.Request) bool

	describe                         string      // 接口描述
	request                          string      // 请求的请求示例
	dataSource                       interface{} // 数据源
	response                         string      // 接口返回示例
	st_request                       interface{}
	params_request                   map[string]string
	st_response                      interface{}
	reqHeader                        map[string]string
	supplement                       string
	delheader                        []string
	codeMsg                          map[string]string
	codeField                        string
	groupKey, groupLable, groupTitle string
	delReqHeader                     []string
}

func (rt *Route) ApiExitGroup() *Route {
	// 退出文档的组
	rt.codeField = ""
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
	rt.groupLable = lable
	rt.groupTitle = title

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

func (rt *Route) Bind(s interface{}) *Route {
	// 接口补充说明

	rt.dataSource = s
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
			d := ApiDocument[id]
			d.Api = append(d.Api, *doc)
			// ApiDocument[id].Api = append(ApiDocument[id].Api, *doc)
			ApiDocument[id] = d

		} else {
			keys[rt.groupKey] = *count
			d := Doc{
				Title: rt.groupTitle,
				Api:   make([]Document, 0),
			}
			d.Api = append(d.Api, *doc)
			ApiDocument[*count] = d

			sideUrl := fmt.Sprintf("/-/api/%d.html", *count)
			sidebar[sideUrl] = rt.groupLable
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

	if rt.delReqHeader == nil {
		rt.delReqHeader = make([]string, 0)
	}
	rt.delReqHeader = append(rt.delReqHeader, k)
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
	if rt.delheader == nil {
		rt.delheader = make([]string, 0)
	}
	rt.delheader = append(rt.delheader, k)
	return rt
}

func (rt *Route) AddMidware(handle func(http.ResponseWriter, *http.Request) bool) *Route {
	if rt.midware == nil {
		rt.midware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	rt.midware = append(rt.midware, handle)
	return rt
}

func (rt *Route) DelMidware(handle func(http.ResponseWriter, *http.Request) bool) *Route {
	if rt.midware == nil {
		rt.delmidware = make([]func(http.ResponseWriter, *http.Request) bool, 0)
	}
	rt.delmidware = append(rt.delmidware, handle)
	return rt
}
