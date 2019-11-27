package xmux

import "net/http"

type Route struct {
	// 组里面也包括路由 后面的其实还是patter和handle, 还没到handle， 这里的key是个method
	method map[string]http.Handler
	//allHandle http.Handler
}

func (rt *Route) Post(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodPost] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Get(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodGet] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Delete(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodDelete] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Head(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodHead] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Options(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodOptions] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Connect(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodConnect] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Patch(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodPatch] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Trace(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodTrace] = http.HandlerFunc(handler)
	return rt
}

func (rt *Route) Put(handler func(http.ResponseWriter, *http.Request)) *Route {
	rt.method[http.MethodPut] = http.HandlerFunc(handler)
	return rt
}

//func (rt *Route) All(handler func(http.ResponseWriter, *http.Request)) *Route {
//	rt.allHandle = http.HandlerFunc(handler)
//	return rt
//}
