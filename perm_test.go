package xmux

import "testing"

func TestPerm(t *testing.T) {
	pl := GetPerm([]string{"C", "U", "R", "D"}, 8)
	t.Log(pl)
}
