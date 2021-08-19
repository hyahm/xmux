package xmux

import (
	"testing"
)

func TestNil(t *testing.T) {
	a := make([]string, 0)
	var b []string
	b = nil
	a = append(a, b...)
}
