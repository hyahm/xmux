package xmux

import (
	"net/http"
	"strings"
)

type mr map[string]*Route

func (mr mr) Add(url string, rt *Route) {
	mr[url] = rt
}

func (mr mr) AppendTo(pattern string, doc *Doc) {
	for url, v := range mr {
		if url == pattern {
			continue
		}
		if strings.Contains(url, "/-/css/") {
			continue
		}
		if strings.Contains(url, "/-/js/") {
			continue
		}

		document := v.makeDoc()
		document.Url = url
		document.Supplement = v.supplement
		for mt, _ := range v.method {
			document.Method = mt
			if mt == http.MethodGet {
				if v.params_request != nil {
					document.Url += GetOpt(v.params_request)
				}
			} else {
				if v.st_request != nil {
					document.Opt = PostOpt(v.st_request)
				}
			}
			doc.Add(document)
		}
	}
}
