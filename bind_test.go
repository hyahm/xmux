package xmux

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	a := ""
	t.Log(strings.Split(a, ","))
	t.Log(len(strings.Split(a, ",")))
}
