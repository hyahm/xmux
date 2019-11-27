package xmux

import (
	"fmt"
	"strings"
)

//var Var map[string]interface{}
var Var map[string]string

// match url , is true is a regx, false is fullurl
func match(path string, newpath string, varlist []string) (string, []string, bool) {

	// /article/24/content

	start := strings.Index(path, "{")
	end := strings.Index(path, "}")
	if start == -1 && end == -1 {
		//非正则的
		return path, varlist, false
	} else if start >= 0 && end > 0 && end > start {
		//正则匹配的
		re := strings.Trim(path[start+1:end], " ")
		if re == "" {
			panic("invaild uri " + path)
		} else {
			prefix := path[:start]
			fmt.Println(prefix)
			//判断:
			ts := strings.Split(re, ":")
			if len(ts) == 2 {
				if ts[0] == "int" {
					prefix += "(\\d+)"
				} else {
					prefix += "(\\w+)"
				}
				varlist = append(varlist, strings.Trim(ts[1], " "))
			} else if len(ts) == 1 {
				prefix += "(\\w+)"
				varlist = append(varlist, strings.Trim(ts[0], " "))
			} else {
				panic("invaild uri " + path)
			}
			newpath += prefix
			if end+1 == len(path) {
				// last url
				newpath += "$"
				return newpath, varlist, true
			} else {
				return match(path[end+1:], newpath, varlist)
			}
		}
	} else {
		panic("invaild uri " + path)
	}
}
