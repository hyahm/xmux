package xmux

// bridge  数据二次封装

type Data struct {
	Header map[string]string
	Request interface{}
	
}

type Bridge  struct {
	Var map[string]interface{}  // 参数
	*Data   // 请求头和body
}