package xmux

// 这是 pattern 对应的
// 这是路径对应的各种方法
// api文档的操作
// type PatternRoute map[string]*Route

// func (pmr PatternRoute) Add(url string, rt *Route) {
// 	pmr[url] = rt
// }

// func (ptn PatternRoute) AppendTo(count *int) {

// 	for url, v := range ptn {
// 		// 初始化document
// 		for md, rt := range v {
// 			doc := Document{
// 				Opt:     make([]option, 0),
// 				Callbak: make([]option, 0),
// 			}
// 			doc.Url = url

// 			doc.Describe = rt.describe
// 			doc.Header = rt.reqHeader
// 			if rt.apiDelReqHeader != nil {
// 				for _, v := range rt.apiDelReqHeader {
// 					delete(doc.Header, v)
// 				}
// 			}
// 			if rt.st_response != nil {
// 				doc.Callbak = postOpt(rt.st_response)
// 			}

// 			doc.Request = rt.request
// 			doc.Response = rt.response
// 			doc.CodeField = rt.codeField
// 			doc.CodeMsg = rt.codeMsg
// 			if doc.CodeField == "" {
// 				doc.CodeField = "code"
// 			}
// 			doc.Supplement = rt.supplement
// 			doc.Method = md
// 			if md == http.MethodGet {
// 				if rt.params_request != nil {
// 					doc.Url += getOpt(rt.params_request)
// 				}
// 			} else {
// 				if rt.st_request != nil {
// 					doc.Opt = postOpt(rt.st_request)
// 				}
// 			}
// 			rt.makeDoc(url, count, &doc)
// 		}

// 	}
// }
