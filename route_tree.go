package xmux

import (
	"fmt"

	"github.com/hyahm/xmux/auth"
)

type Meta struct {
	MenuType string `json:"menu_type"`
	Name     string `json:"name"`
	// URL      string `json:"url"`
	// Method   string `json:"method"`
	// UUID       string `json:"uuid"`
	// ParentUUID string `json:"parent_uuid"`
	// Icon string `json:"icon"`
}

type MenuTree struct {
	// 节点id，唯一标识， 根据 url. method, MenuType, name 生成
	MenuId     string      `json:"menu_id"`
	Uuid       string      `json:"uuid"`
	URL        string      `json:"url"`
	Method     []string    `json:"method"`
	ParentUUID string      `json:"parent_uuid"`
	Meta       Meta        `json:"meta"`
	Roles      []string    `json:"-"`
	Children   []*MenuTree `json:"children"`
}

func (m *MenuTree) makeMenuId() {
	m.MenuId = auth.Md5([]byte(fmt.Sprintf("%s-%s-%s-%s", m.URL, m.Method, m.Meta.MenuType, m.Meta.Name)))
}

// 扁平化菜单树， 方便权限校验， 方便插入数据库
func FlattenMenuTree(tree []*MenuTree) []*MenuTree {
	var list []*MenuTree

	var dfs func(items []*MenuTree)
	dfs = func(items []*MenuTree) {
		for _, item := range items {
			if item == nil {
				continue
			}

			// 重点：创建一个副本，清空 children！！
			flatItem := *item
			flatItem.Children = nil // 清空嵌套

			// 加入一维切片
			list = append(list, &flatItem)

			// 继续递归子节点
			dfs(item.Children)
		}
	}

	dfs(tree)
	return list
}

// BuildMenuTree 构建路由树, 拿到菜单树
func BuildMenuTree(list []MenuTree) []*MenuTree {
	// 1. 初始化 Map，预分配空间以提高性能
	nodeMap := make(map[string]*MenuTree, len(list))
	for i := range list {
		item := list[i]
		if item.Uuid == "" {
			continue
		}
		mt := &MenuTree{
			Uuid:       item.Uuid,
			Meta:       item.Meta,
			URL:        item.URL,
			Method:     item.Method,
			ParentUUID: item.ParentUUID,
			Children:   make([]*MenuTree, 0), // 初始化切片，避免前端收到 null
		}
		if mt.Meta.MenuType == "" || mt.Meta.Name == "" {
			continue // 跳过无效节点
		}
		mt.makeMenuId()
		nodeMap[item.Uuid] = mt
	}

	rootNodes := make([]*MenuTree, 0)

	// 2. 构建父子关系
	for i := range list {
		item := list[i]
		node, ok := nodeMap[item.Uuid]
		if !ok || node == nil { // 关键：过滤 nil
			continue
		}
		// 判断是否为顶级节点
		if item.ParentUUID == "root" {
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

// BuildRouteTree 构建路由树, 拿到菜单树
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
		mt.makeMenuId()
		nodeMap[item.Uuid] = mt
	}

	rootNodes := make([]*MenuTree, 0)

	// 2. 构建父子关系
	for i := range list {
		item := list[i]
		node := nodeMap[item.Uuid]

		// 判断是否为顶级节点
		if item.ParentUUID == "root" {
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

// FilterMenuTree 根据 menu_id 列表过滤菜单树
// 规则：
// 1. 存在于 menuIdSet 中的节点，其父节点 menu_type != a 才保留
// 2. 下级节点 menu_type == a 且 name == "" 保留
func FilterMenuTree(tree []*MenuTree, menuIdSet map[string]bool, filterType string) []*MenuTree {
	var result []*MenuTree

	for _, node := range tree {
		if node == nil {
			continue
		}

		// 递归处理子节点
		filteredChildren := FilterMenuTree(node.Children, menuIdSet, filterType)

		// ===================== 规则3：强制保留条件 =====================
		// menu_type == filterType 并且 name == "" → 必须保留
		forceKeep := node.Meta.MenuType == filterType && node.Meta.Name == ""

		// ===================== 规则1：在白名单里 =====================
		inWhiteList := menuIdSet[node.MenuId]

		// ===================== 规则2：子节点有保留，父节点必须保留 =====================
		hasValidChild := len(filteredChildren) > 0

		// ===================== 最终判断：满足任意一个就保留 =====================
		if forceKeep || inWhiteList || hasValidChild {
			newNode := *node
			newNode.Children = filteredChildren
			result = append(result, &newNode)
		}
	}

	return result
}
