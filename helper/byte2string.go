package helper

import "unsafe"

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 去掉多余的符号
func CompressBytes(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	for _, v := range src {
		if v == '\n' || v == '\t' || v == '\r' || v == ' ' || v == '\v' || v == '\f' || v == 0x85 || v == 0xA0 {
			continue
		}
		dst = append(dst, v)
	}
	return dst
}
