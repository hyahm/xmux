package xmux

import (
	"testing"
)

func G2[T any]() *RouteGroup {
	g2 := NewRouteGroup().AddPageKeys("g2")
	g2.Get("/g2", nil)
	return g2
}

func G3[T any]() *RouteGroup {
	g3 := NewRouteGroup()
	g3.Get("/g3", nil).AddPageKeys("g3")
	return g3
}

func G4[T any]() *RouteGroup {
	g4 := NewRouteGroup().DelPageKeys("g4")
	g4.Get("/g4", nil)
	g4.Get("/g4_1", nil).AddPageKeys("g4_1")
	return g4
}

func G5[T any]() *RouteGroup {
	g5 := NewRouteGroup().AddPageKeys("g5")
	g5.Get("/g5", nil).AddPageKeys("g2")
	return g5
}

func G6[T any]() *RouteGroup {
	g6 := NewRouteGroup().AddPageKeys("g6")
	g6.Get("/g6", nil).DelPageKeys("g3")
	g6.Get("/g6_1", nil).DelPageKeys("g1")
	return g6
}

func G1[T any]() *RouteGroup {

	g1 := NewRouteGroup().AddPageKeys("g1")
	g1.AddGroup(G2[T]())
	return g1
}

// 嵌套带中间件的组路由，组路由的中间件覆盖最外层的router的中间件
func TestGroupPagesCoveredRouterPages(t *testing.T) {
	router := NewRouter()
	router.AddGroup(G2[any]())
	router.DebugAssignRoute("/g2")

}

// 路由的中间件覆盖路由组的中间件
func TestRoutePagesCoveredGroupPages(t *testing.T) {
	router := NewRouter()
	router.AddGroup(G5[any]())
	router.DebugAssignRoute("/g5")

}

// 组路由删除中间，此组的中间件都将删除，但是外层的没影响
func TestDeleteGroupPagesWillDeleteRouterPages(t *testing.T) {
	router := NewRouter().AddPageKeys("g1")
	router.Get("/r1", nil)
	router.AddGroup(G4[any]())
	router.Get("/r2", nil)
	// router.DebugAssignRoute("/r1")
	// router.DebugAssignRoute("/r2")
	router.DebugAssignRoute("/g4")
	// router.DebugAssignRoute("/g4_1")

}

// router和组里面都带了中间件，组里面的中间件会覆盖最外层
// 删除的话，只有删除组里面的才会删除
// 路由删除中间件，将会删除对应Router或者组的中间件， 如果外层的的不一样
func TestOnlyOnePagesBeStayRouteDelPagesWillDeleteTheLastestPages(t *testing.T) {
	router := NewRouter()
	router.AddGroup(G6[any]())

}

// router里面挂载了全局的中间件， 测试下面的url是否都挂载了
func TestOnlyRouter(t *testing.T) {
	router := NewRouter()
	router.Get("/r1", nil)
	router.Get("/r2", nil)
}
