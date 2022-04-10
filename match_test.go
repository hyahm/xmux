package xmux

import (
	"testing"
)

func TestMatch(t *testing.T) {
	l, v := match("/download/aa{re:v([1-9]+)-([1-9]+)-([1-9]+):v1,v2,v3}bb{word}")
	t.Log(l)
	t.Log(v)
}
