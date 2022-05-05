package xmux

import (
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	a := []byte(`{ "aaa": "bbbb", is；了哈楼上的回复
	"cccc: 555"	`)
	b := make([]byte, 0, len(a))
	for _, v := range a {
		if v == '\n' || v == '\t' || v == '\r' || v == ' ' || v == '\v' || v == '\f' || v == 0x85 || v == 0xA0 {
			continue
		}
		b = append(b, v)
	}
	fmt.Println(b)
	fmt.Println(string(b))
}
