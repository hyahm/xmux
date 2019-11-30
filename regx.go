package xmux

import (
	"log"
	"strings"
)

//var Var map[string]interface{}

// match url , is true is a regx, false is fullurl
func match(path string, newpath string, varlist []string) (string, []string, bool) {

	start := strings.Index(path, "{")
	if start == -1 {
		return path, varlist, false
	}
	end := -1

	sc := strings.Index(path[start:], "/")

	if sc != -1 {
		pp := path[:start+sc]
		end = strings.LastIndex(pp, "}")
	} else {
		end = strings.LastIndex(path, "}")
	}
	//找到最后一个}
	//strings.LastIndex()

	if start == -1  {
		//非正则的

	} else if start >= 0 && end > 0 && end > start {

		//正则匹配的
		re := strings.Trim(path[start+1:end], " ")
		if re == "" {
			log.Fatal("invalid uri " + path)
		} else {
			prefix := path[:start]
			//判断:
			ts := strings.Split(re, ":")
			if len(ts) == 3 {
				//正则 匹配
				// /asdf/{re:([a-z]{1,3})([0-9]{1,2}):ch,num}
				if ts[0] == "re" {
					// 检测参数是否匹配, 同时禁止匹配()
					pfc := strings.Count(ts[1], "(")
					sfc := strings.Count(ts[1], ")")
					if pfc != sfc {
						log.Fatal("can not include ( or ) ," + path)
					}
					//查看后面参数是否匹配
					vl := strings.Split(ts[2], ",")
					if len(vl) != sfc {
						log.Fatal("variable not matched , " + path)
					}
					prefix += ts[1]
					varlist = append(varlist, vl...)
				} else {
					log.Fatal("invalid uri ," + path)
				}
			} else 	if len(ts) == 2 {
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
				log.Fatal("invalid uri ," + path)
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
		log.Fatal("invalid uri ," + path)
	}
	return "", nil,false
}


// 将多个连续斜杠合成一个， 去掉末尾的斜杠，
// 例如   /asdf/sadf//asdfsadf/asdfsdaf////as///, 转为-》 /asdf/sadf/asdfsadf/asdfsdaf/as
func slash(s string) string {

	sl := strings.Split(s, "/")
	n := make([]string, 0)
	for _, v := range sl {
		if v != "" {
			n = append(n, v)
		}
	}
	return "/" + strings.Join(n, "/")
}