package xmux

type mr map[string]*Route

func (mr mr) Add(url string, rt *Route) {
	mr[url] = rt
}
