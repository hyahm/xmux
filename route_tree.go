package xmux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Tree struct {
	Url         string                                          `json:"url"`
	Modules     []func(http.ResponseWriter, *http.Request) bool `json:"modules"`
	PostModules []func(http.ResponseWriter, *http.Request) bool `json:"post_modules"`
	Method      []string                                        `json:"methods"`
	Roles       []string                                        `json:"roles"`
	Children    []*Tree                                         `json:"children"`
}

var routeTree []*Tree

func initRouteTree() {
	routeTree = make([]*Tree, 0)
}

// addgroup 的操作
func addGroupRouteTree(parent, children []*Tree) {
	if len(children) == 0 {
		return
	}
	if parent == nil {
		parent = make([]*Tree, 0)
	}
	parent = append(parent, &Tree{
		Children: children,
	})
	fmt.Println(*parent[0])
}

func GetRouteTreeJson() {
	b, _ := json.MarshalIndent(routeTree, "", "  ")
	fmt.Println(b)
}
