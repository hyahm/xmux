# xmux， go语言 路由(router)
应该是基于原生net.http包唯一一个带缓存， 使用简单并强大的路由， 内嵌接口文档，告别另外写文档的烦恼

### 已完成功能
- [x] xmux.NewGroupRoute(), 为了避免各种异常， 请使用自带的来创建路由
- [x] 支持路由分组
- [x] 支持全局请求头， 组请求头， 私有请求头
- [x] 支持自定义method， 多method
- [x] 支持正则匹配和参数获取
- [x] 完全匹配优先于正则匹配
- [x] 正则匹配支持（int(\d+), word(\w+), re, all(.*?)，不写默认 string([^\/])）建议使用string
- [x] 支持三大全局handle ,MethodnotFound(忘记写方法), MethodNotAllowed(method没定义), HandleNotFound(没有找到页面), Options请求）  
- [x] 强大的模块让你的代码模块化变得非常简单 
- [x] 中间件支持 
- [x] 内嵌接口文档
- [x] 数据绑定
- [x] 增加数据结构绑定， 适合模块间传递
- [x] 增加websocket， 可以学习，不建议使用, 如果其他的不好可以试试  
- [x] 集成pprof， router.AddGroup(xmux.Pprof())
- [x] 支持代理（参考于:  https://github.com/ouqiang/goproxy）

### 最简单的运行
```
package main

import (
	"net/http"

	"github.com/hyahm/xmux"
)

func main() {
	router := xmux.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>hello world!<h1>"))
	})
	router.Run()
}

```

打开 localhost:8080 就能看到 hello world!

### 添加了组的概念

> aritclegroup.go
```go
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

var Article *xmux.GroupRoute

func init() {
	Article = xmux.NewGroupRoute()
	Article.Get("/{int:id}", hello)

}
```
> main.go
```go
func main() {
	router := xmux.NewRouter()
	router.AddGroup(aritclegroup.Article)
}

```
### 更灵活的匹配
```go
func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
	router.Get("/{all:age}", Who)   // 这个可以匹配任何路由
}
```

记住， 是100%，  此路由优先匹配完全匹配规则， 匹配不到再寻找 正则匹配， 加快了寻址速度  
访问 /get -> 返回 show   
访问  /post   -> 返回 Who  

### 自动检测重复项,
```go
func main() {
	router := xmux.NewRouter()
	router.Get("/get",show) // 不同请求分别处理
	router.Get("/get",show) // 不同请求分别处理

}
写一大堆路由，  有没有重复的都不知道  
运行上面将会报错， 如下  
2019/11/29 21:51:11 pattern duplicate for /get

```
###  自动格式化url
将任意多余的斜杠去掉例如
/asdf/sadf//asdfsadf/asdfsdaf////as///, 转为-》 /asdf/sadf/asdfsadf/asdfsdaf/as
```go
func main() {
	router := xmux.NewRouter()
	router.Get("/get",show) // 不同请求分别处理
	router.Get("/get/",show) // 不同请求分别处理

}

所以运行上面将会报错，/get/被转为 /get 如下  
2019/11/29 21:51:11 pattern duplicate for /get

```


### 三大全局handle
```go
HandleOptions:        handleoptions(),   //这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理, 默认会返回ok， 也可以自定义
HandleNotFound: handleNotFound(),   // 默认返回404 ， 也可以自定义
HanleFavicon methodNotAllowed(),    // 默认请求 favicon

// 默认调用的方法如下， 没有找到路由
func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	})
}

```

###  Header ， 中间件 和 模块 

- 路由有3种路由头  
全局路由： 所有请求都会带上这个  
组路由： 所有组的路由都会带上这个， 还有带上全局的， 组的请求头覆盖全局的  
私有路由： 单一路由的请求头， 属于某组的话， 带上组路由头， 全局的话带上全局的  
优先级  
私有路由 > 组路由 > 全局路由  (如果存在优先级大的就覆盖优先级小的)

- 模块类似上面header         
不过优先级相反    
全局路由 > 组路由 > 私有路由 (如果存在优先级大的就覆盖优先级小的)
```go

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world home"))
	return
}

func mid() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("77777")
		return
	})
}

func hf(w http.ResponseWriter, r *http.Request)  bool {
	fmt.Println("44444444444444444444444444")
	r.Header.Set("name", "cander")
	
	return true
}

func hf1(w http.ResponseWriter, r *http.Request)  bool {
	fmt.Println("66666")
	fmt.Println(r.Header.Get("name"))
	return false
}

func TestHome(t *testing.T) {
	router := xmux.NewRouter()
	router.Get("/home/{test}",home).AddModule(hf).SetHeader("name", "cander").AddModule(hf1)
	var a string
	// client := http.Client{}
	r, err := http.NewRequest("GET", "/home/asdf", strings.NewReader(a))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	t.Log(w.Code)

	t.Log(w.Body.String())
}

```
模块会根据添加的顺序执行  
44444444444444444444444444  
66666  
cander  
页面将看到  
hello world home  

- 中间件最多只能有一个， 功能较多建议使用模块  
优先级与header 一样， 中间件如下， 这是个计算执行时间的例子  
```go
func GetExecTime(handle func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	handle(w, r)
	fmt.Printf("url: %s -- addr: %s -- method: %s -- exectime: %f\n", r.URL.Path, r.RemoteAddr, r.Method, time.Since(start).Seconds())
}

```

### 跨域处理
跨域主要是添加请求头的问题, 其余框架一般都是借助中间件来设置   
但是本路由借助上面请求头设置 大大简化跨域配置  

```
func main() {
	router := xmux.NewRouter()
	router.Slash = true
	router.SetHeader("Access-Control-Allow-Origin", "*")  // 主要的解决跨域, 因为是全局的请求头， 所以后面增加的路由全部支持跨域
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type,Access-Token,X-Token,Origin,smail,authorization")  // 新增加的请求头
	router.Get("/", index)
	router.Run()
}

```

### 适合在当前handle 的 module， 中间件， handle 传递值   
> 生命周期从定义开始， 到此handle执行完毕将被释放
```go
func filter(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")
	xmux.GetData(r).Set("name","xmux")
	r.Header.Set("bbb", "ccc")
	return false
}

func name(w http.ResponseWriter, r *http.Request) {
	var name string
	gd := xmux.GetData(r)
	if gd != nil {
		name= gd.Get("name").(string)
	}
	 
	w.Write([]byte("hello world " + name))
	return
}
router.Pattern("/aaa/{name}").Get(name).AddModule(filter).AddModule(login)
```
### 模块， 中间件， handle 执行顺序
- 模块 > 中间件 > handle  

### 获取正则匹配的参数
```go
func Who(w http.ResponseWriter, r *http.Request) {
fmt.Println(xmux.Var(r)["name"])
fmt.Println(xmux.Var(r)["age"])
w.Write([]byte("yes is mine"))
return
}

``` 
### 数据绑定（Bind(), 与Module 一起使用）
将数据结构绑定到此 Handle 里， 通过读取r.Body 来解析数据
因为解析的代码都是一样的， 绑定后可以共用同一份代码
```
func JsonToStruct(w http.ResponseWriter, r *http.Request) bool {
	// 任何报错信息， 直接return true， 就是此handle 直接执行完毕了， 不继续向后面走了
	resp := &response.Response{}
	if goconfig.ReadBool("debug", false) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write(resp.ErrorE(err))
			return true
		}
		err = json.Unmarshal(b, xmux.GetData(r).Data)
		if err != nil {
			w.Write(resp.ErrorE(err))
			return true
		}
	} else {
		err := json.NewDecoder(r.Body).Decode(xmux.GetData(r).Data)
		if err != nil {
			w.Write(resp.ErrorE(err))
			return true
		}

	}
	return false
}


Router.Post("/important/add", handle.AddName).Bind(&model.DataName{}).AddMidware(midware.JsonToStruct)
Router.Post("/important/add", handle.AddStd).Bind(&model.DataStd{}).AddMidware(midware.JsonToStruct)
Router.Post("/important/add", handle.AddFoo).Bind(&model.DataFoo{}).AddMidware(midware.JsonToStruct)


func AddName(w http.ResponseWriter, r *http.Request) {
	// 这里的data 就是解析后出来的数据
	data := xmux.GetData(r).Data.(*model.DataName)
	/// todo...
}
```

### 有model可以没有Handle(module 本身可以当handle使用）
确定没有midware， 是可以没有handle的  
```
func NoHandleModule(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("hello world"))
	// 这里必须返回true， 否则接口报错
	return true
}
user.Get("/no/handle", nil).AddModule(NoHandleModule)

resp: /no/handle
return:  hello world
```


### 自动修复请求的url
例如： 请求的url 是这个样子的
http://www.hyahm.com/mmm///af/af,  默认是请求不到的
但是设置后
```go
router := xmux.NewRouter()
router.Slash = true  
```
是可以直接访问 http://www.hyahm.com/mmm/af/af 这个地址的请求

### 匹配路由
支持以下5种
 word   只匹配数字和字母（默认）    
 string  匹配所有不含/的字符    
 int  匹配整数    
 all： 匹配所有的包括/   
 re： 自定义正则   
/aaa/{name}          这个和下面一个一样， 省略类型， 默认是string  
/aaa/{string:name}   这个和上面一样， string类型  
/aaa/{int:name}       这个匹配int类型   
/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd:name,age}  这个是一段里面匹配了2个参数 name, age，   
大括号表示是个匹配规则，里面2个冒号分割了3部分 起头的  
第一个: re表示用到自定义正则，只有re才会有2个冒号分割,   
第二个： 正则表达式， 里面不能出现: 需要提取的参数用()括起来，   
第三个: 参数名， 前面有多少对()， 后面就需要匹配多少个参数， 用逗号分割   
例如： /aaa/adfaasf16sd  
这个是匹配的， name: aa   age: 16  
```
xmux.Var(r)["name"] 
```
后面会增加自定义正则匹配


### 编写接口文档， 
> 使用接口文档,  第一个参数是组路由名， 第二个参数是挂载的路由uri
== 组路由里面的静态文件 默认挂在 /-/css/xxx.css 和  /-/js/xxx.js 下 ==
== 动态路由   /-/api/{int}.html  ==

```go
// 所有的文档相关的方法都以Api开头， 文档只支持单路由的单请求方式， 多请求方式会乱, 调用的时候只会显示到当前位置以上的路由
router := xmux.NewRouter()
api := router.ShowApi("/doc")
router.ShowApi(api).
ApiCreateGroup("test", "api test", "apitest").  //增加了侧边栏 所有组路由或单路由必须加上这个才会显示, 第一个参数是组key, 第二个是组的标题， 第三个是侧边栏url显示的文字 ， 或者添加到某个组上 ApiAddGroup(key), 组路由添加的key 会被子路由继承， 如果不想显示可以ApiAddGroup 挂载到其他路由或者 ApiExitGroup， 移除此组
ApiDescribe("这是home接口的测试").  // 接口的简述
ApiReqHeader("content-type", "application/json"). // 接口请求头
ApiReqStruct(&Home{}).    // 接口请求参数， 由struct tag 提供（可以是结构体，也可以是结构体指针）
ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).   // 接口请求示例
ApiResStruct(Call{}).     // 接口返回参数， 由struct tag 提供 （可以是结构体，也可以是结构体指针）
ApiResponseTemplate(`{"code": 0, "msg": ""}`).  // 接口返回示例
ApiSupplement("这个是接口的说明补充， 没补充就不填"). // 接口补充
ApiCodeField("133").    // 错误码字段
ApiCodeMsg("1", "56").ApiCodeMsg("3", "akhsdklfhl")   // 错误码说明， 多次调用添加多次
```
>  接口请求参数tag 示例
```
type Home struct {
	Addr   string `json:"addr" type:"string" need:"是" default:"深圳" information:"家庭住址"`
	People int    `json:"people" type:"int" need:"是" default:"1" information:"有多少个人"`
}
```
>  接口接收参数tag 示例, 比请求示例少了 default
```
type Call struct {
	Code int    `json:"code" type:"int" need:"是" information:"错误返回码"`
	Msg  string `json:"msg" type:"string" need:"否" information:"错误信息"`
}
```
### 压力测试
```
PS F:\xmux> go test -bench .
goos: windows
goarch: amd64
pkg: github.com/hyahm/xmux
BenchmarkOneRoute-8                      1000000              1007 ns/op              32 B/op          1 allocs/op
Benchmark404Many-8                         61387             21462 ns/op           13044 B/op        170 allocs/op
BenchmarkMux-8                           1270545               932 ns/op              32 B/op          1 allocs/op
BenchmarkMuxAlternativeInRegexp-8        1285495               939 ns/op              32 B/op          1 allocs/op
BenchmarkManyPathVariables-8               75098             16003 ns/op           10432 B/op        130 allocs/op
```


###  代理使用
```
func main() {
	# proxy 也是一个路由
	proxy := NewProxy()
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

```
