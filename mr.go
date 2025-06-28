package xmux

import "net/http"

// string 对应的是url
type UrlRoute map[string]*Route

func (mr UrlRoute) SetHeader(key, value string) UrlRoute {
	for _, route := range mr {
		route.SetHeader(key, value)
	}
	return mr
}

func (mr UrlRoute) AddModule(handles ...func(http.ResponseWriter, *http.Request) bool) UrlRoute {

	for method := range mr {
		mr[method].AddModule(handles...)
	}
	return mr
}

func (mr UrlRoute) Bind(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].Bind(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) BindResponse(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].BindResponse(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) BindByContentType(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].BindByContentType(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) BindForm(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].BindForm(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) BindJson(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].BindJson(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) BindXml(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].BindXml(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) BindYaml(dest interface{}) UrlRoute {
	for method := range mr {
		mr[method].BindYaml(Clone(dest))
	}
	return mr
}

func (mr UrlRoute) AddPageKeys(pagekeys ...string) UrlRoute {
	for _, route := range mr {
		route.AddPageKeys(pagekeys...)
	}
	return mr
}

func (mr UrlRoute) DelHeader(key string) UrlRoute {
	for _, route := range mr {
		route.DelHeader(key)
	}
	return mr
}

// get route by method. if not found will return nil
func (mr UrlRoute) GetRoute(method string) *Route {
	if _, ok := mr[method]; ok {
		return mr[method]
	}
	return nil
}

func (mr UrlRoute) DelModule(handles ...func(http.ResponseWriter, *http.Request) bool) UrlRoute {
	for _, route := range mr {
		route.DelModule(handles...)
	}
	return mr
}

func (mr UrlRoute) DelPageKeys(pagekeys ...string) UrlRoute {
	for _, route := range mr {
		route.DelPageKeys(pagekeys...)
	}
	return mr
}
