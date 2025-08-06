package xmux

import (
	"log"
	"regexp"
	"strings"
	"testing"
)

type routerTest struct {
	title       string
	method      string
	pattern     string
	url         string
	keys        []string
	values      []string
	isSlash     bool
	shouldMatch bool
}

// 测试路径是否能匹配
func TestPattern(t *testing.T) {
	tests := []routerTest{
		{
			title:       "完全匹配",
			pattern:     "/aaaa/bbbb",
			url:         "/aaaa/bbbb",
			keys:        nil,
			values:      nil,
			isSlash:     false,
			method:      "GET",
			shouldMatch: true,
		},
		{
			title:       "字符串正则匹配",
			pattern:     "/aaaa/{bbbb}",
			url:         "/aaaa/bbbb",
			method:      "GET",
			isSlash:     false,
			keys:        []string{"bbbb"},
			values:      []string{"bbbb"},
			shouldMatch: true,
		},
		{
			title:       "数字也是字符串",
			pattern:     "/aaaa/{bbbb}",
			url:         "/aaaa/12334",
			method:      "POST",
			isSlash:     false,
			keys:        []string{"bbbb"},
			values:      []string{"12334"},
			shouldMatch: true,
		},
		{
			title:       "带参数的字符串匹配",
			pattern:     "/aaaa/{string:bbbb}",
			url:         "/aaaa/hioj",
			method:      "GET",
			isSlash:     false,
			keys:        []string{"bbbb"},
			values:      []string{"hioj"},
			shouldMatch: true,
		},
		{
			title:       "int类型正则匹配",
			pattern:     "/aaaa/{int:bbbb}",
			method:      "GET",
			isSlash:     false,
			url:         "/aaaa/joijoa324",
			shouldMatch: false,
		},
		{
			title:       "int类型正则匹配",
			pattern:     "/aaaa/{int:bbbb}",
			url:         "/aaaa/334",
			keys:        []string{"bbbb"},
			method:      "GET",
			isSlash:     false,
			values:      []string{"65555555555555"},
			shouldMatch: true,
		},
		{
			title:       "多路径截断",
			pattern:     "/aaaa////{int:bbbb}",
			url:         "/aaaa/334",
			method:      "POST",
			isSlash:     false,
			keys:        []string{"bbbb"},
			values:      []string{"334"},
			shouldMatch: true,
		},
		{
			title:       "多路径截断",
			pattern:     "/aaaa///{string:bbbb}///",
			url:         "/aaaa/sdf/",
			method:      "POST",
			isSlash:     true,
			keys:        []string{"bbbb"},
			values:      []string{"sdf"},
			shouldMatch: true,
		},
		{
			title:       "多路径截断",
			pattern:     "/aaaa///{string:bbbb}///",
			url:         "/aaaa/sdf/",
			method:      "POST",
			isSlash:     true,
			keys:        []string{"bbbb"},
			values:      []string{"sdf"},
			shouldMatch: true,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			testPattern(t, test)
		})
	}
}

func testPattern(t *testing.T, test routerTest) {

	t.Log("before: ", test.pattern)
	test.pattern = PrettySlash(test.pattern)
	t.Log("after: ", test.pattern)
	// r.makeRoute(test.pattern)
	patternMatched := false

	url, _, ok := makeRoute(test.pattern)
	if !ok {
		patternMatched = url == test.url
	} else {
		re := regexp.MustCompile(url)
		patternMatched = re.MatchString(test.url)
	}

	if patternMatched != test.shouldMatch {
		t.Errorf("not matched")
	}

}

// 正则的url 换成 正常的url
func testFormat(path string, newpath string) string {

	start := strings.Index(path, "{")
	end := strings.Index(path, "}")

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
			} else {
				log.Fatal("invalid uri ," + path)
			}
		} else if len(ts) == 2 {
			if ts[0] == "int" {
				prefix += "(\\d+)"
			} else {
				prefix += "(\\w+)"
			}
		} else if len(ts) == 1 {
			prefix += "(\\w+)"
		} else {
			log.Fatal("invalid uri ," + path)
		}
		newpath += prefix
		if end+1 == len(path) {
			// last url
			newpath += "$"
			return newpath
		} else {
			return testFormat(path[end+1:], newpath)
		}
	}
	return ""

}
