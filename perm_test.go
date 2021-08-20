package xmux

import "testing"

// func TestPerm(t *testing.T) {
// 	t.Log(Read)
// 	t.Log(Create)
// 	t.Log(Update)
// 	t.Log(Delete)
// 	t.Log(SetPerm(Create | Update | Delete | Read))
// }

func TestPerm(t *testing.T) {
	pl := GetPerm([]string{"C", "U", "R", "D"}, 8)
	t.Log(pl)
}
