package xmux

import (
	"bytes"
	"encoding/json"
	"html/template"
)

type JsonStr string

func (js JsonStr) Json() string {
	b := &bytes.Buffer{}
	err := json.Indent(b, []byte(js), "", "    ")
	if err != nil {
		return ""
	}
	return b.String()
}

// 获取每一个路由的url， 请求， header， body， response（后面3个需要手动添加）, 只支持json 格式, 数据结构绑定
type Document struct {
	Describe   string // 接口描述
	Url        string
	Method     string
	Header     map[string]string // 请求头
	Opt        []option
	Callbak    []option
	Request    string // json 请求实例
	Response   string // json 返回实例
	Supplement string // 最后的补充说明
	CodeField  string
	CodeMsg    map[string]string
}

type Doc struct {
	Api     []Document        // 单路由或组路由
	Title   string            // 内容主标题， 与label 一致
	Search  map[string]int    // 所有路由都是用的这个, 暂时不管
	Sidebar map[string]string // 所有路由的这个都是相同的, url 和 lebel
}

var sidebar map[string]string // 所有路由的这个都是相同的, url 和 lable

var apiDocument map[int]Doc // 应该是这个, 打开某地址， 返回对应的信息

var keys map[string]int // 侧边栏与id 对应

func NewDocs(r *Router) {
	apiDocument = make(map[int]Doc)
	sidebar = make(map[string]string) // url to lable
	keys = make(map[string]int)
	count := 1
	// 单路由
	r.route.AppendTo(&count)
	r.tpl.AppendTo(&count)

}

// 获取html 基础页面
func NewTemplate() *template.Template {
	var t *template.Template
	name := "doc"
	var html *template.Template
	if t == nil {
		t = template.New(name)
	}
	if name == t.Name() {
		html = t
	} else {
		html = t.New(name)
	}
	t, _ = html.Parse(tpl)

	return t
}

func NewHomeTemplate() *template.Template {
	var t *template.Template
	name := "doc"
	var html *template.Template
	if t == nil {
		t = template.New(name)
	}
	if name == t.Name() {
		html = t
	} else {
		html = t.New(name)
	}
	t, _ = html.Parse(tpl0)

	return t
}
