package xmux

import (
	"net/http"
	"testing"
)

func m1(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func m2(w http.ResponseWriter, r *http.Request) bool {
	return false
}
func m3(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func m4(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func m5(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func TestModule(t *testing.T) {
	router := NewRouter().AddModule(m1, m2)
	router.AddGroup(subgroup())
	router.Get("/get", nil)
	router.Post("/", nil).DelModule(m2)

	{
		s := router.route["/"].module.funcOrder
		if GetFuncName(s[0]) != "github.com/hyahm/xmux.m1" || len(s) != 1 {
			t.Fail()
		}
		get := router.route["/get"].module.funcOrder
		if len(get) != 2 {
			t.Fail()
		}
		if GetFuncName(get[0]) != "github.com/hyahm/xmux.m1" {
			t.Fail()
		}
		if GetFuncName(get[1]) != "github.com/hyahm/xmux.m2" {
			t.Fail()
		}
	}

	{
		subpost := router.route["/sub/post"].module.funcOrder
		if len(subpost) != 2 {
			t.Fail()
		}
		if GetFuncName(subpost[0]) != "github.com/hyahm/xmux.m2" {
			t.Fail()
		}
		if GetFuncName(subpost[1]) != "github.com/hyahm/xmux.m3" {
			t.Fail()
		}
	}

	{
		sub1get := router.route["/sub1/get"].module.funcOrder
		if len(sub1get) != 4 {
			t.Fail()
		}
		if GetFuncName(sub1get[0]) != "github.com/hyahm/xmux.m2" {
			t.Fail()
		}
		if GetFuncName(sub1get[1]) != "github.com/hyahm/xmux.m3" {
			t.Fail()
		}
		if GetFuncName(sub1get[2]) != "github.com/hyahm/xmux.m4" {
			t.Fail()
		}
		if GetFuncName(sub1get[3]) != "github.com/hyahm/xmux.m5" {
			t.Fail()
		}

	}

	{
		sub1post := router.route["/sub1/post"].module.funcOrder
		if len(sub1post) != 2 {
			t.Fail()
		}
		if GetFuncName(sub1post[0]) != "github.com/hyahm/xmux.m2" {
			t.Fail()
		}
		if GetFuncName(sub1post[1]) != "github.com/hyahm/xmux.m4" {
			t.Fail()
		}

	}

}

func subgroup() *GroupRoute {
	sub := NewGroupRoute().AddModule(m2, m3).DelModule(m1)
	sub.Get("/sub/get", nil)
	sub.Post("/sub/post", nil)
	sub.Any("/sub/any", nil)
	sub.AddGroup(sub1group())
	return sub
}

func sub1group() *GroupRoute {
	sub1 := NewGroupRoute().AddModule(m4)
	sub1.Get("/sub1/get", nil).AddModule(m5)
	sub1.Post("/sub1/post", nil).DelModule(m3)
	sub1.Any("/sub1/any", nil)
	return sub1
}
