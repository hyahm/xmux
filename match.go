package xmux

import (
	"log"
	"path"
	"strings"
)

const (
	str      = "([^\\/]+)"
	word     = `([a-zA-Z0-9_]+)`
	interger = "(\\d+)"
	all      = "(.*?)"
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

func cleanPath(p string) string {
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

// 返回正则表达式的url 和 params(key)
func match(path string) (string, []string) {
	// 如果是空的，直接报错
	if strings.Trim(path, " ") == "" {
		panic("pattern empty")
	}
	// 返回三个参数，  （正则）路径， 正则的参数， 是否是正则
	// 按照/分段计算, 计算的时候需要计算后面的斜杠， 如果有的长度会多一个
	pl := strings.Split(path, "/")
	var (
		pathlist []string
		varlist  []string
	)
	for _, v := range pl {
		block, vl := macheOne(v, "", []string{})
		if block != "" {
			pathlist = append(pathlist, block)
		}
		if len(vl) > 0 {
			varlist = append(varlist, vl...)
		}

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

// macheOne 解析带{}占位符的URL路径，返回处理后的路径和提取的变量名列表
// path: 待解析的完整URL路径模板
// newpath: 拼接中的处理后路径
// varlist: 收集的变量名列表
func macheOne(path, newpath string, varlist []string) (string, []string) {
	// 找第一个{的位置
	firstPrev := strings.Index(path, "{")
	if firstPrev == -1 {
		// 找不到{，直接拼接剩余路径并返回
		return newpath + path, varlist
	}

	// 提取{之前的头部字符串
	head := path[:firstPrev]
	// 剩余需要处理的路径（去掉{）
	remainingPath := path[firstPrev+1:]

	// 边界：{是最后一个字符（无闭合}）
	if remainingPath == "" {
		log.Fatal("路径: " + path + " 有问题，{后无内容，缺少闭合}和变量定义")
	}

	if strings.HasPrefix(remainingPath, "re:") {
		// 处理正则类型：{re:正则表达式:变量名} 格式
		// 跳过"re:"后找第一个:（分割正则和变量名）
		colonAfterRe := strings.Index(remainingPath[3:], ":")
		if colonAfterRe == -1 {
			log.Fatal("路径: " + path + " 有问题，正则类型需格式 {re:正则:变量名}")
		}

		// 找闭合}的位置
		closingBrace := strings.Index(remainingPath[3+colonAfterRe:], "}")
		if closingBrace == -1 {
			log.Fatal("路径: " + path + " 有问题，正则类型缺少闭合}")
		}

		// 提取变量名（支持多个变量用,分隔）
		varNamesStr := remainingPath[3+colonAfterRe+1 : 3+colonAfterRe+closingBrace]
		varNamesStr = strings.Trim(varNamesStr, " ")
		if varNamesStr != "" {
			if strings.Contains(varNamesStr, ",") {
				// 分割多变量并去空格
				for _, v := range strings.Split(varNamesStr, ",") {
					trimmed := strings.Trim(v, " ")
					if trimmed != "" { // 过滤空字符串
						varlist = append(varlist, trimmed)
					}
				}
			} else {
				varlist = append(varlist, varNamesStr)
			}
		} else {
			log.Fatal("路径: " + path + " 有问题，正则类型{}内变量名不能为空")
		}

		// 拼接处理后的路径（正则表达式部分）
		regexPart := remainingPath[3 : 3+colonAfterRe]
		newpath += head + regexPart
		// 更新剩余待处理路径（跳过当前}）
		path = remainingPath[3+colonAfterRe+closingBrace+1:]
	} else {
		// 处理普通类型：{类型:变量名} 或 {变量名} 格式
		// 找闭合}的位置
		closingBrace := strings.Index(path, "}")
		if closingBrace == -1 {
			log.Fatal("路径: " + path + " 有问题，缺少闭合}")
		}
		if closingBrace == firstPrev+1 {
			log.Fatal("路径: " + path + " 有问题，{}之间不能为空")
		}

		// 提取{}内的内容（类型:变量名 或 变量名）
		contentInBraces := path[firstPrev+1 : closingBrace]
		// 检查:的数量（最多1个）
		colonCount := strings.Count(contentInBraces, ":")
		if colonCount > 1 {
			log.Fatal("路径: " + path + " 有问题，非正则的{}里面最多只能出现一个:")
		}

		// 解析普通类型的匹配规则和变量名
		re, opt := normal(contentInBraces)
		if opt == "" {
			log.Fatal("路径: " + path + " 有问题，非正则类型{}内变量名不能为空")
		}
		varlist = append(varlist, opt)

		// 拼接处理后的路径（匹配规则部分）
		newpath += head + re
		// 更新剩余待处理路径（跳过当前}）
		path = path[closingBrace+1:]
	}

	// 递归处理剩余路径（先trim避免空字符串递归）
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return newpath, varlist
	}
	return macheOne(trimmedPath, newpath, varlist)
}

// normal 解析非正则类型的{}内容，返回匹配正则和变量名
// 支持格式：变量名 或 类型:变量名（类型：int/word/all/string）
// func normal(path string) (string, string) {
// 	ts := strings.SplitN(path, ":", 2) // SplitN避免多个:分割错误
// 	switch len(ts) {
// 	case 1:
// 		// 无类型，默认word匹配，变量名为ts[0]
// 		opt := strings.Trim(ts[0], " ")
// 		return word, opt
// 	case 2:
// 		// 有类型：类型:变量名
// 		typ := strings.Trim(ts[0], " ")
// 		typ = strings.ToLower(typ)
// 		opt := strings.Trim(ts[1], " ")

// 		switch typ {
// 		case "int":
// 			return interger, opt
// 		case "word":
// 			return word, opt
// 		case "all":
// 			return all, opt
// 		case "string":
// 			return str, opt
// 		default:
// 			// 未知类型，默认使用word匹配
// 			return word, opt
// 		}
// 	default:
// 		panic("路径: " + path + " 有问题，非正则类型{}内格式错误")
// 	}
// }
