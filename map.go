package xmux

type mstringstring map[string]string

func (ss mstringstring) clone() mstringstring {
	nss := make(mstringstring)
	for k, v := range ss {
		nss[k] = v
	}
	return nss
}

func (ss mstringstring) delete(su mstringstruct) {
	for k := range su {
		delete(ss, k)
	}
}

func (ss mstringstring) add(su mstringstring) {
	for k, v := range su {
		ss[k] = v
	}
}

type mstringstruct map[string]struct{}

func (su mstringstruct) clone() mstringstruct {
	nsu := make(mstringstruct)
	for k := range su {
		nsu[k] = struct{}{}
	}
	return nsu
}

func (su mstringstruct) add(msu mstringstruct) {
	for k := range msu {
		su[k] = struct{}{}
	}
}

func (su mstringstruct) delete(msu mstringstruct) {
	for k := range msu {
		delete(su, k)
	}
}
