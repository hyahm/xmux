package xmux

// a 减去b的差集
func Subtract[T comparable](a, b []T) []T {
	// want []string{"aa"}
	bm := make(map[T]struct{})
	for _, v := range b {
		bm[v] = struct{}{}
	}
	temp := make([]T, 0, len(a))
	for _, v := range a {
		if _, ok := bm[v]; !ok {
			temp = append(temp, v)
		}
	}
	return temp
}

func SubtractSliceMap[T comparable](a []T, b map[T]struct{}) []T {
	// want []string{"aa"}
	temp := make([]T, 0, len(a))
	for _, v := range a {
		if _, ok := b[v]; !ok {
			temp = append(temp, v)
		}
	}
	return temp
}

// 重复的K  有多个key 也只返回一个
func SliceExsit(s1, s2 []string) (string, bool) {
	mm := make(map[string]struct{})
	for _, v := range s2 {
		mm[v] = struct{}{}
	}
	for _, v := range s1 {
		if _, ok := mm[v]; ok {
			return v, true
		}
	}
	return "", false
}
