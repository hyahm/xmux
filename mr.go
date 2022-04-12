package xmux

import (
	"net/http"
)

// string 对应的是method
type MethodsRoute map[string]*Route

func (mr MethodsRoute) SetHeader(key, value string) MethodsRoute {
	for _, route := range mr {
		route.SetHeader(key, value)
	}
	return mr
}

func (mr MethodsRoute) MiddleWare(midware http.Handler) MethodsRoute {
	for _, route := range mr {
		route.MiddleWare(midware)
	}
	return mr
}

func (mr MethodsRoute) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) MethodsRoute {
	for _, route := range mr {
		route.AddModule(handles...)
	}
	return mr
}

func (mr MethodsRoute) AddPageKeys(pagekeys ...string) MethodsRoute {
	for _, route := range mr {
		route.AddPageKeys(pagekeys...)
	}
	return mr
}

func (mr MethodsRoute) DelHeader(key string) MethodsRoute {
	for _, route := range mr {
		route.DelHeader(key)
	}
	return mr
}

func (mr MethodsRoute) DelMiddleWare(midware http.Handler) MethodsRoute {
	for _, route := range mr {
		route.DelMiddleWare(midware)
	}
	return mr
}

func (mr MethodsRoute) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) MethodsRoute {
	for _, route := range mr {
		route.DelModule(handles...)
	}
	return mr
}

func (mr MethodsRoute) DelPageKeys(pagekeys ...string) MethodsRoute {
	for _, route := range mr {
		route.DelPageKeys(pagekeys...)
	}
	return mr
}
