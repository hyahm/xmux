package xmux

type mstringstring map[string]string

func (ss mstringstring) deleteHeader(su mstringstruct) {
	for k := range su {
		delete(ss, k)
	}
}

func (ss mstringstring) addHeader(su mstringstring) {
	for k, v := range su {
		ss[k] = v
	}
}

type mstringstruct map[string]struct{}
