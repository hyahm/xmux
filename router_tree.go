package xmux

import "encoding/json"

type Meta struct {
	Url     string   `json:"url"`
	Methods []string `json:"methods"`
}

type RouterTree struct {
	Metas    []Meta      `json:"metas"`
	Children *RouterTree `json:"children,omitempty"`
}

var enableRouterTree bool

var routerTrees RouterTree

func initRouterTree() {
	routerTrees = RouterTree{
		Metas: make([]Meta, 0),
	}
}

func (r *RouterTree) AddChild(child *RouterTree) {
	// 如果当前节点没有 children，直接退出
	if len(child.Metas) == 0 && child.Children == nil {
		return
	}
	if len(child.Metas) > 0 {
		if r.Children == nil {
			r.Children = &RouterTree{
				Metas: make([]Meta, 0),
			}
		}
		r.Children.Metas = append(r.Children.Metas, child.Metas...)
	}
	if child.Children != nil {
		if r.Children.Children == nil {
			r.Children.Children = &RouterTree{
				Metas: make([]Meta, 0),
			}
		}
		// 否则递归往下找，直到最后一个节点
		r.Children.AddChild(child.Children)
	}
}

func ToJson() []byte {
	b, _ := json.MarshalIndent(routerTrees, "", "  ")
	return b
}
