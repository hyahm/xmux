package xmux

import (
	"net/http"
	"sync"
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
	router.AddGroup(sub3group())
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

	{
		sub4post := router.route["/sub4/post"].module.funcOrder
		if len(sub4post) != 4 {
			t.Fail()
		}
		if GetFuncName(sub4post[0]) != "github.com/hyahm/xmux.m1" {
			t.Fail()
		}
		if GetFuncName(sub4post[1]) != "github.com/hyahm/xmux.m2" {
			t.Fail()
		}
		if GetFuncName(sub4post[2]) != "github.com/hyahm/xmux.m3" {
			t.Fail()
		}
		if GetFuncName(sub4post[3]) != "github.com/hyahm/xmux.m4" {
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
	sub1.AddGroup(sub2group())
	return sub1
}

func sub2group() *GroupRoute {
	sub1 := NewGroupRoute().AddModule(m5)
	sub1.Get("/sub2/get", nil)
	sub1.Post("/sub2/post", nil).DelModule(m3)
	sub1.Any("/sub2/any", nil)
	return sub1
}

func sub3group() *GroupRoute {
	sub1 := NewGroupRoute()
	sub1.Get("/sub3/get", nil)
	sub1.Post("/sub3/post", nil)
	sub1.Any("/sub3/any", nil)
	sub1.AddGroup(sub4group())
	return sub1
}

func sub4group() *GroupRoute {
	sub1 := NewGroupRoute().AddModule(m3)
	sub1.Get("/sub4/get", nil)
	sub1.Post("/sub4/post", nil).AddModule(m4)
	sub1.Any("/sub4/any", nil)
	return sub1
}

func TestClonseModule(t *testing.T) {
	m := &module{
		filter:    make(map[string]struct{}),
		funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		mu:        sync.RWMutex{},
	}

	m.add(m1, m2)
	mc := m.cloneMudule()
	mc.add(m3, m4)
	t.Log(m.filter)
}

func TestClonseModule2(t *testing.T) {
	m := &module{
		filter:    make(map[string]struct{}),
		funcOrder: make([]func(w http.ResponseWriter, r *http.Request) bool, 0),
		mu:        sync.RWMutex{},
	}

	m.add(m1, m2)
	mc := &module{
		filter:    m.filter,
		funcOrder: m.funcOrder,
		mu:        sync.RWMutex{},
	}
	mc.add(m3, m4)
	t.Log(m.filter)
}
