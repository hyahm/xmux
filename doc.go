package xmux

// 自动生成接口文档

// 获取每一个路由的url， 请求， header， body， response（后面3个需要手动添加）, 只支持json 格式, 数据结构绑定
type Bind struct {
	Title     string // 接口名
	Describe  string // 接口描述
	Url       string
	Method    []string          // 这里应该是个数组, 一般只有一个值
	Header    map[string]string // 请求头， 一般都是 Content-Type: application/x-www-form-urlencoded; charset=UTF-8  或 Content-Type: application/json; charset=UTF-8
	Requester interface{}       //  请求接口
	Responser interface{}       // 数据接口
	Request   []byte            // json 请求实例
	Response  []byte            // json 数据 实例
}

// 所有的接口数据
type Doc struct {
	*Bind
	Msg map[int]string // 这个是所有的错误码统计
}
