package xmux

// func TestMap(t *testing.T) {
// 	defer golog.Sync()
// 	aaa := &AAA{}

// 	mr := aaa.Add("shangshan", http.MethodConnect, http.MethodDelete, http.MethodGet)
// 	golog.Infof("%p", mr)
// 	t.Log(aaa.m["shangshan"])
// 	mr.Describe("7777")
// 	t.Log(aaa.m["shangshan"])
// }

// func (a *AAA) Add(url string, methods ...string) MethodsRoute {
// 	if a.m == nil {
// 		a.m = make(map[string]MethodsRoute)
// 	}
// 	mr := make(MethodsRoute)
// 	for _, method := range methods {
// 		mr[method] = &Route{}
// 	}
// 	a.m[url] = mr
// 	return a.m[url]
// }

// func (mr MethodsRoute) Describe(describe string) MethodsRoute {
// 	golog.Infof("%p", mr)
// 	for method := range mr {
// 		golog.Info(method)
// 		mr[method].describe = describe
// 	}
// 	golog.Infof("%p", mr)
// 	return mr
// }
