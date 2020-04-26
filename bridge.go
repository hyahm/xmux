package xmux

// bridge  数据二次封装

type Data struct {
	Header   map[string]string
	Request  interface{}
	Response interface{}
}

type Bridge struct {
	Var   map[string]string // 参数
	*Data                   // 请求头和body
}
