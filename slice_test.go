package xmux

import "testing"

func TestUnit(t *testing.T) {
	a := []string{"aa", "bb", "cc"}
	b := []string{"bb", "cc"}

	temp := Subtract(a, b)

	t.Log(temp)
}
