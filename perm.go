package xmux

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
