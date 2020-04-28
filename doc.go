package xmux

import (
	"html/template"
)

// 获取每一个路由的url， 请求， header， body， response（后面3个需要手动添加）, 只支持json 格式, 数据结构绑定
type Document struct {
	Describe   string // 接口描述
	Url        string
	Method     string
	Header     map[string]string // 请求头
	Opt        []Opt
	Callbak    []Opt
	Request    string // json 请求实例
	Response   string // json 返回实例
	Supplement string // 最后的补充说明
	CodeField  string
	CodeMsg    map[string]string
}

// 所有的接口数据
type Doc struct {
	Api   []Document
	Title string // 这个是所有的错误码统计
}

func (doc *Doc) Add(line Document) *Doc {
	doc.Api = append(doc.Api, line)
	return doc
}

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
