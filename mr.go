package xmux

import (
	"net/http"
)

type mr map[string]*Route

func (mr mr) Add(url string, rt *Route) {
	mr[url] = rt
}

func (mr mr) AppendTo(count *int) {

	for url, v := range mr {
		// 初始化document

		doc := Document{
			Opt:     make([]Opt, 0),
			Callbak: make([]Opt, 0),
		}
		doc.Describe = v.describe
		doc.Header = v.reqHeader
		if v.delReqHeader != nil {
			for _, v := range v.delReqHeader {
				delete(doc.Header, v)
			}
		}
		if v.st_response != nil {
			doc.Callbak = PostOpt(v.st_response)
		}
		doc.Url = url
		doc.Request = v.request
		doc.Response = v.response
		doc.CodeField = v.codeField
		doc.CodeMsg = v.codeMsg
		if doc.CodeField == "" {
			doc.CodeField = "code"
		}
		doc.Supplement = v.supplement
		for mt, _ := range v.method {
			doc.Method = mt
			if mt == http.MethodGet {
				if v.params_request != nil {
					doc.Url += GetOpt(v.params_request)
				}
			} else {
				if v.st_request != nil {
					doc.Opt = PostOpt(v.st_request)
				}
			}
		}
		v.makeDoc(url, count, &doc)
	}
}
