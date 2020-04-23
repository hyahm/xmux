package xmux

// 自动生成接口文档

// 获取每一个路由的url， 请求， header， body， response（后面3个需要手动添加）

type Doc struct {
	Title string  // 接口名
	Describe   // 接口描述
	Url string
	Method []string // 这里应该是个数组, 一般只有一个值
	Header map[string]string  // 请求头， 一般都是 Content-Type: application/x-www-form-urlencoded; charset=UTF-8  或 Content-Type: application/json; charset=UTF-8 
	Body []byte     //  
	Response []byte
}


type Docs struct {
	*Doc 
	Msg map[int]string  // 这个是所有的错误码统计
}