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
