package xmux

import (
	"testing"
)

func TestMatch(t *testing.T) {
	l, v := match("/download/aa{re:v([1-9]+)-([1-9]+)-([1-9]+):v1,v2,v3}bb{word}")
	// l, v := match("/download/{name}")
	t.Log(l)
	t.Log(v)
}

func TestCleanPath(t *testing.T) {
	path := PrettySlash("/asdf/sadf//asdfsadf/asdfsdaf////as///")
	t.Log(path)
}

// func TestCleanPath1(t *testing.T) {
// 	path := cleanPath("/asdf/sadf//asdfsadf/asdfsdaf////as///")
// 	t.Log(path)
// }

func TestMatch1(t *testing.T) {
	l, v := match("/download/aa{re:v([1-9]+)-([1-9]+)-([1-9]+):v1,v2,v3}bb{word}")
	// l, v := match("/download/{name}")
	t.Log(l)
	t.Log(v)
}
