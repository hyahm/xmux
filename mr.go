package xmux

import "net/http"

// string 对应的是method
type MethodsRoute map[string]*Route

func (mr MethodsRoute) SetHeader(key, value string) MethodsRoute {
	for _, route := range mr {
		route.SetHeader(key, value)
	}
	return mr
}

func (mr MethodsRoute) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) MethodsRoute {

	for method := range mr {
		mr[method].AddModule(handles...)
	}
	return mr
}

func (mr MethodsRoute) Bind(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].Bind(Clone(dest))
	}
	return mr
}

func (mr MethodsRoute) BindResponse(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].BindResponse(Clone(dest))
	}
	return mr
}

func (mr MethodsRoute) BindByContentType(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].BindByContentType(Clone(dest))
	}
	return mr
}

func (mr MethodsRoute) BindForm(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].BindForm(Clone(dest))
	}
	return mr
}

func (mr MethodsRoute) BindJson(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].BindJson(Clone(dest))
	}
	return mr
}

func (mr MethodsRoute) BindXml(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].BindXml(Clone(dest))
	}
	return mr
}

func (mr MethodsRoute) BindYaml(dest interface{}) MethodsRoute {
	for method := range mr {
		mr[method].BindYaml(Clone(dest))
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
