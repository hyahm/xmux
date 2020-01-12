package xmux

import (
	"testing"
)

type matchStruct struct {
	title      string
	path       string
	expect     string
	expectlist []string
}

func TestMatch(t *testing.T) {
	tests := []matchStruct{
		{
			title:  "完全匹配",
			path:   "/aaaa",
			expect: "/aaaa",
		},
		{
			title:  "多路径完全匹配",
			path:   "/aaaa/bbb",
			expect: "/aaaa/bbb",
		},
		{
			title:      "默认正则",
			path:       "/{name}/bbb",
			expect:     "^/([^/ ])/bbb$",
			expectlist: []string{"name"},
		},
		{
			title:      "带类型正则",
			path:       "/{int:name}/bbb",
			expect:     "^/(\\d+)/bbb$",
			expectlist: []string{"name"},
		},
		{
			title:      "re正则",
			path:       "/{re:(.?+):name}/bbb",
			expect:     "^/(.?+)/bbb$",
			expectlist: []string{"name"},
		},
		{
			title:      "path正则1",
			path:       "/{re:(.?+):name}/bbb",
			expect:     "^/(.?+)/bbb$",
			expectlist: []string{"name"},
		},
		{
			title:      "re多参数正则",
			path:       "/{re:([a-z])444([0-9]):word,num}/bbb",
			expect:     "^/([a-z])444([0-9])/bbb$",
			expectlist: []string{"word", "num"},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			testMatch(t, test)
		})
	}
}

func testMatch(t *testing.T, m matchStruct) {
	path, varlist := match(m.path)
	if path != m.expect {
		t.Fatalf("title: %s", m.title)
	}
	for i, v := range varlist {
		if v != m.expectlist[i] {
			t.Fatalf("title: %s", m.title)
		}
	}

}

func TestCount(t *testing.T) {
	c := parenthesesCount("(asdf))asdf(asdfasdfasdf)AsdF(AsDF)SA(DF)ASDF(SD(F)SADF(AS)DF(ASD(F)SADF()ASDF", 0)
	if c != 8 {
		t.Fatal("error count")
	}
}

type tPath struct {
	url   string
	reUrl string
}

func TestPath(t *testing.T) {
	tests := []tPath{
		{
			url:   "/aaaa/bbb/{name}",
			reUrl: "/aaaa/bbb/ccc.asdf.png/",
		},
	}
	for _, test := range tests {
		path, vl := match(test.url)
		t.Log(vl)
		if !matchUrlTest(test.reUrl, path) {
			t.Fatal("not match")
		}
	}

}
