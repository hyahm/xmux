package xmux

import (
	"strings"
)

const (
	str      = "([^\\/]+)"
	word     = "(\\w+)"
	interger = "(\\d+)"
	all      = "(.*?)"
	sep      = ","
)

// 将多个连续斜杠合成一个， 去掉末尾的斜杠，
// 例如   /asdf/sadf//asdfsadf/asdfsdaf////as///, 转为-》 /asdf/sadf/asdfsadf/asdfsdaf/as
func prettySlash(s string) string {
	sl := strings.Split(s, "/")
	n := make([]string, 0, len(sl))
	for _, v := range sl {
		if v != "" {
			n = append(n, v)
		}
	}
	return "/" + strings.Join(n, "/")
}

// 返回正则表达式 和 参数
func match(path string) (string, []string) {
	// 如果是空的，直接报错
	if strings.Trim(path, " ") == "" {
		panic("pattern empty")
	}
	// 返回三个参数，  （正则）路径， 正则的参数， 是否是正则
	// 分段
	pl := strings.Split(path, "/")
	pathlist := make([]string, 0)
	varlist := make([]string, 0)
	for _, v := range pl {
		newpath, vl := macheOne(v)
		pathlist = append(pathlist, newpath)
		varlist = append(varlist, vl...)
	}
	var newpath string
	// 合拼路径
	if len(varlist) == 0 {
		// 完全匹配
		newpath = strings.Join(pathlist, "/")
	} else {
		// 正则匹配
		newpath = "^" + strings.Join(pathlist, "/") + "$"
	}
	return newpath, varlist
}

func macheOne(path string) (string, []string) {
	varlist := make([]string, 0)
	// 找第一个{
	start := strings.Index(path, "{")
	if start == -1 {
		// 找不到就是完全匹配
		return path, varlist
	}
	// 保存头部可能存在的字符串
	head := path[:start]
	end := -1

	// 找最后一个}
	end = strings.LastIndex(path, "}")
	if end == -1 {
		// 找不到就是完全匹配,只是路径带了{
		return path, varlist
	}
	// 过来的都是有正则规则
	if start > end {
		// }{ 这样的也是完全匹配
		return path, varlist
	}
	// 保存尾部可能存在的字符串
	tail := path[end+1:]
	//   ==========  进入正则匹配区
	// 去掉{} 和 2边的空格
	re := strings.Trim(path[start+1:end], " ")
	if re == "" {
		// /{}  类似这样的会
		panic("invalid uri " + path)
	} else {
		//判断: 目前只支持
		// 没有:, {name}, {int:id}, {re:(.khjk)dfdf([a|b]):path,word}
		// 一个:
		// 二个 :
		// 其他的全是错误的匹配
		ts := strings.Split(re, ":")
		switch len(ts) {

		case 1:
			// 没有:
			opt := strings.Trim(ts[0], " ")
			varlist = append(varlist, opt)
			return head + word + tail, varlist
		case 2:
			// 一个:
			// 判断类型
			typ := strings.Trim(ts[0], " ")
			typ = strings.ToLower(typ)
			opt := strings.Trim(ts[1], " ")
			varlist = append(varlist, opt)
			switch typ {
			case "int":
				return head + interger + tail, varlist
			case "word":
				return head + word + tail, varlist
			case "all":
				return head + all + tail, varlist
			case "string":
				return head + str + tail, varlist
			default:
				// 默认使用path匹配
				return head + word + tail, varlist
			}

		case 3:
			// 二个:
			// 参数必须是re， 如果不是。 默认改成re
			typ := strings.Trim(ts[0], " ")
			typ = strings.ToLower(typ)
			opts := strings.Split(ts[2], sep)
			for _, opt := range opts {
				opt = strings.Trim(opt, " ")
				varlist = append(varlist, opt)
			}

			if typ != "re" {
				// panic
				panic("pattern not support" + path)
			}
			// 正则2边不能有空格
			ts[1] = strings.Trim(ts[1], " ")

			// 正则必须与参数个数匹配
			// 参数必须是, 分割
			// 查找有多少对小括号
			pc := parenthesesCount(ts[1], 0)
			if pc != len(varlist) {
				panic("pattern not support" + path)
			}
			return head + ts[1] + tail, varlist
		default:
			panic("pattern not support" + path)
		}
	}

}

func parenthesesCount(s string, c int) int {
	// 计算有多少对小括号
	start := strings.Index(s, "(")
	if start == -1 {
		return c
	}

	end := strings.Index(s, ")")
	if end == -1 {
		return c
	}
	if end < start {
		s = s[end+1:]
		return parenthesesCount(s, c)
	}
	s = s[end+1:]
	c++
	return parenthesesCount(s, c)
}
