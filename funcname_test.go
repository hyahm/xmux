package xmux

import (
	"fmt"
	"testing"
)

func TestFunc(t *testing.T) {
	a := 12
	b := func() {
		fmt.Println(11)
	}
	t.Log(GetFuncName(a))
	t.Log(GetFuncName(b))
}
