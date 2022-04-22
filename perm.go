package xmux

import (
	"net/http"
)

func DefaultPermissionTemplate(w http.ResponseWriter, r *http.Request) (post bool) {
	// 如果是管理员的，直接就过
	// if uid == <adminId> {
	// 	retrun false
	// }

	// roles := []string{"env", "important"}
	// 内置的方法最大支持8种权限，如果想要更多可以自己实现
	var pl = []string{"Read", "Create", "Update", "Delete"}
	// map 的key 对应页面的value  value 对应二进制位置(从右到左)
	permissionMap := make(map[string]int)
	for k, v := range pl {
		permissionMap[v] = k
	}
	// 假如权限拿到二进制对应的10进制数据是下面
	perm := make(map[string]uint8)
	perm["env"] = 14       // 00001110   {"Delete", "Create", "Update"}
	perm["important"] = 10 // 00001010   {"Create", "Delete"}
	perm["project"] = 4    // 00000100   {"Update"}

	//
	pages := GetInstance(r).Get(PAGES).(map[string]struct{})
	// 如果长度为0的话，说明任何人都可以访问
	if len(pages) == 0 {
		return false
	}
	//  请求/project/read     map[admin:{} project:{}]
	// 判断 pages 是否存在 perm
	// 注意点： 这里的页面权限本应该只会匹配到一个， 这个是对于的页面权限的值
	page := ""
	// 判断页面权限的
	hasPerm := false
	for role := range perm {
		if _, ok := pages[role]; ok {
			hasPerm = true
			page = role
			break
		}
	}
	if !hasPerm {
		w.Write([]byte("没有页面权限"))
		return true
	}
	// permMap := make(map[string]bool)
	result := GetPerm(pl, perm[page])
	handleName := GetInstance(r).GetFuncName()
	// 这个值就是判断有没有这个操作权限
	if !result[permissionMap[handleName]] {
		w.Write([]byte("没有权限"))
		return true
	}
	// 先拿到pl 对应名称的 索引
	//         8        4        2          1
	//		 delete	 update	 create		read
	//  bit   0        0       0         0
	/*
		用户表
		id
		1
		权限表
		id      uid   roles                       perm
		1       1     "env"                       0-15
		2       1     "important"
	*/
	return false
}

// 给定一个权限组， 顺序对应2进制的值必须是 1 << index,
// 最后返回对应位置 是不是 1 的 bool类型的切片
// 如果传入的切片大于8，只获取8个
func GetPerm(permList []string, flag uint8) []bool {
	length := len(permList)
	if length > 8 {
		length = 8
	}
	res := make([]bool, length)
	x := ToBinaryString(flag)
	for index := range permList {
		res[index] = x[7-index:8-index] == "1"
	}
	return res
}
