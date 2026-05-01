package xmux

import (
	"fmt"
	"strings"

	"github.com/hyahm/xmux/auth"
)

type Meta struct {
	MenuType string `json:"menu_type"`
	Name     string `json:"name"`
	// URL      string `json:"url"`
	// Method   string `json:"method"`
	// UUID       string `json:"uuid"`
	// ParentUUID string `json:"parent_uuid"`
	Icon string `json:"icon"`
}

type MenuTree struct {
	// 节点id，唯一标识， 根据 url. method, role, name 生成
	MenuId     string      `json:"menu_id"`
	Uuid       string      `json:"uuid"`
	URL        string      `json:"url"`
	Method     string      `json:"method"`
	ParentUUID string      `json:"parent_uuid"`
	Meta       Meta        `json:"meta"`
	Roles      []string    `json:"-"`
	Children   []*MenuTree `json:"children"`
}

func (m *MenuTree) MakeMenuId() {
	m.MenuId = auth.Md5([]byte(fmt.Sprintf("%s-%s-%s-%s", m.URL, m.Method, strings.Join(m.Roles, ","), m.Meta.Name)))
}

func BuildRouteTree(list []MenuTree) []*MenuTree {
	// 1. 初始化 Map，预分配空间以提高性能
	nodeMap := make(map[string]*MenuTree, len(list))
	for i := range list {
		item := list[i]
		mt := &MenuTree{
			Uuid:       item.Uuid,
			Meta:       item.Meta,
			URL:        item.URL,
			Method:     item.Method,
			ParentUUID: item.ParentUUID,
			Children:   make([]*MenuTree, 0), // 初始化切片，避免前端收到 null
		}
		mt.MakeMenuId()
		nodeMap[item.Uuid] = mt
	}

	rootNodes := make([]*MenuTree, 0)

	// 2. 构建父子关系
	for i := range list {
		item := list[i]
		node := nodeMap[item.Uuid]

		// 判断是否为顶级节点
		// 建议增加判断：如果父 ID 为空字符串，也视作顶级节点（视你的业务逻辑而定）
		if item.ParentUUID == "root" || item.ParentUUID == "" {
			rootNodes = append(rootNodes, node)
			continue
		}

		// 找到父节点并挂载
		if parent, ok := nodeMap[item.ParentUUID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			// 【关键修复】：如果找不到父节点，说明它是孤儿节点
			// 做法 A：视作顶级节点（防止数据丢失）
			// 做法 B：直接忽略（取决于你的业务需求）
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes
}
