package xmux

import (
	"path"
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
// func prettySlash(s string) string {
// 	sl := strings.Split(s, "/")
// 	n := make([]string, 0, len(sl))
// 	for _, v := range sl {
// 		if v != "" {
// 			n = append(n, v)
// 		}
// 	}
// 	return "/" + strings.Join(n, "/")
// }

func PrettySlash(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}

// func cleanPath(p string) string {
// 	if p == "" {
// 		return "/"
// 	}
// 	if p[0] != '/' {
// 		p = "/" + p
// 	}
// 	np := path.Clean(p)
// 	// path.Clean removes trailing slash except for root;
// 	// put the trailing slash back if necessary.
// 	if p[len(p)-1] == '/' && np != "/" {
// 		// Fast path for common case of p being the string we want:
// 		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
// 			np = p
// 		} else {
// 			np += "/"
// 		}
// 	}
// 	return np
// }

// 返回正则表达式的url 和 params(key)
func match(path string) (string, []string) {
	// 如果是空的，直接报错
	if strings.Trim(path, " ") == "" {
		panic("pattern empty")
	}
	// 返回三个参数，  （正则）路径， 正则的参数， 是否是正则
	// 按照/分段计算
	pl := strings.Split(path, "/")
	var (
		pathlist []string
		varlist  []string
	)
	for _, v := range pl {
		block, vl := macheOne(v, "", []string{})
		pathlist = append(pathlist, block)
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

// path 是完整的url  newpath 是未处理的 path
func macheOne(path, newpath string, varlist []string) (string, []string) {
	// var varlist []string
	// 找第一个{
	firstPrev := strings.Index(path, "{")
	if firstPrev == -1 {
		// 找不到就是完全匹配, 也就不是正则
		return newpath + path, varlist
	}

	// 保存头部可能存在的字符串
	head := path[:firstPrev]
	// 保存已经去掉开头字符串的后面字符串
	newPath := path[firstPrev+1:]

	if path[firstPrev+1:] == "" {
		return newpath + path, varlist
	}
	if newPath[:3] == "re:" {
		// 如果后面是正则，那么先匹配到下一个:
		nextColon := strings.Index(newPath[3:], ":")
		firstSuffix := strings.Index(newPath[3+nextColon:], "}")
		if firstSuffix <= 0 {
			panic("路径: " + path + "有问题，{}之间必须要有别名来获取值或者括号没匹配")
			// 找不到就是完全匹配, 也就不是正则
		}
		vars := newPath[4+nextColon : 3+nextColon+firstSuffix]

		if strings.Count(vars, ",") == 0 {
			varlist = append(varlist, vars)
		} else {
			varlist = append(varlist, strings.Split(vars, ",")...)
		}
		newpath += head + newPath[3:3+nextColon]
		path = newPath[4+nextColon+firstSuffix:]
	} else {
		// 第一个 } 的位置
		firstSuffix := strings.Index(path, "}")
		if firstSuffix <= 0 {
			panic("路径: " + path + "有问题，{}之间必须要有别名来获取值或者括号没匹配")
			// 找不到就是完全匹配, 也就不是正则
		}
		// 按照:切割，分离标识和key
		count := strings.Count(path[firstPrev+1:firstSuffix], ":")
		if count > 1 {
			panic("路径: " + path + "有问题，非正则的{}里面最多只能出现一个:")
		}
		re, opt := normal(path[firstPrev+1 : firstSuffix])
		varlist = append(varlist, opt)
		newpath += head + re
		path = path[firstSuffix+1:]
	}
	if strings.Trim(path, " ") == "" {
		return newpath, varlist
	}
	return macheOne(path, newpath, varlist)
	// 找
	// if firstPrev > 0 {
	// 	// 找不到就是完全匹配, 也就不是正则
	// 	return path, varlist
	// }
	// 找第一个}
	// end = strings.Index(path, "}")
	// if firstSuffix == -1 {
	// 	// 找不到就是完全匹配,只是路径带了{
	// 	return path, varlist
	// }
	// 过来的都是有正则规则
	// if firstPrev > end {
	// 	// }{ 这样的也是完全匹配
	// 	return path, varlist
	// }
	// // 保存尾部可能存在的字符串
	// tail := path[end+1:]
	// //   ==========  进入正则匹配区
	// // 去掉{} 和 2边的空格
	// re := strings.Trim(path[firstprev+1:end], " ")
	// if re == "" {
	// 	// /{}  类似这样的会
	// 	panic("invalid uri " + path)
	// } else {
	// 	//判断: 目前只支持
	// 	// 没有:, {name}, {int:id}, {re:(.khjk)dfdf([a|b]):path,word}
	// 	// 一个:
	// 	// 二个 :
	// 	// 其他的全是错误的匹配
	// 	ts := strings.Split(re, ":")
	// 	switch len(ts) {

	// 	case 1:
	// 		// 没有:
	// 		opt := strings.Trim(ts[0], " ")
	// 		varlist = append(varlist, opt)
	// 		return head + word + tail, varlist
	// 	case 2:
	// 		// 一个:
	// 		// 判断类型
	// 		typ := strings.Trim(ts[0], " ")
	// 		typ = strings.ToLower(typ)
	// 		opt := strings.Trim(ts[1], " ")
	// 		varlist = append(varlist, opt)
	// 		switch typ {
	// 		case "int":
	// 			return head + interger + tail, varlist
	// 		case "word":
	// 			return head + word + tail, varlist
	// 		case "all":
	// 			return head + all + tail, varlist
	// 		case "string":
	// 			return head + str + tail, varlist
	// 		default:
	// 			// 默认使用path匹配
	// 			return head + word + tail, varlist
	// 		}

	// 	case 3:
	// 		// 二个:
	// 		// 参数必须是re， 如果不是。 默认改成re
	// 		typ := strings.Trim(ts[0], " ")
	// 		typ = strings.ToLower(typ)
	// 		opts := strings.Split(ts[2], sep)
	// 		for _, opt := range opts {
	// 			opt = strings.Trim(opt, " ")
	// 			varlist = append(varlist, opt)
	// 		}

	// 		if typ != "re" {
	// 			// panic
	// 			panic("pattern not support" + path)
	// 		}
	// 		// 正则2边不能有空格
	// 		ts[1] = strings.Trim(ts[1], " ")

	// 		// 正则必须与参数个数匹配
	// 		// 参数必须是, 分割
	// 		// 查找有多少对小括号
	// 		pc := parenthesesCount(ts[1], 0)
	// 		if pc != len(varlist) {
	// 			panic("pattern not support" + path)
	// 		}
	// 		return head + ts[1] + tail, varlist
	// 	default:
	// 		panic("pattern not support" + path)
	// 	}
	// }

}

// 非正则匹配
func normal(path string) (string, string) {
	ts := strings.Split(path, ":")
	switch len(ts) {
	case 1:
		// 没有:
		opt := strings.Trim(ts[0], " ")
		return word, opt
	case 2:
		// 一个:
		// 判断类型
		typ := strings.Trim(ts[0], " ")
		typ = strings.ToLower(typ)
		opt := strings.Trim(ts[1], " ")
		switch typ {
		case "int":
			return interger, opt
		case "word":
			return word, opt
		case "all":
			return all, opt
		case "string":
			return str, opt
		default:
			// 默认使用path匹配
			return word, opt
		}
	default:
		panic("")
	}
}

// func parenthesesCount(s string, c int) int {
// 	// 计算有多少对小括号
// 	start := strings.Index(s, "(")
// 	if start == -1 {
// 		return c
// 	}

// 	end := strings.Index(s, ")")
// 	if end == -1 {
// 		return c
// 	}
// 	if end < start {
// 		s = s[end+1:]
// 		return parenthesesCount(s, c)
// 	}
// 	s = s[end+1:]
// 	c++
// 	return parenthesesCount(s, c)
// }
