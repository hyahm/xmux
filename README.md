# xmux， go语言 路由(router)
应该是基于原生net.http包唯一一个带缓存， 使用简单并强大的路由， 内嵌接口文档，告别另外写文档的烦恼

### 已完成功能
- [x] 支持路由分组
- [x] 支持全局请求头， 组请求头， 私有请求头
- [x] 支持自定义method， 多method
- [x] 支持正则匹配和参数获取
- [x] 完全匹配优先于正则匹配
- [x] 正则匹配支持（int(\d+), word(\w+), re, all(.*?)，不写默认 string([^\/])）建议使用string
- [x] 支持四大全局handle ,MethodnotFound(忘记写方法), MethodNotAllowed(method没定义), HandleNotFound(没有找到页面), Options请求）  
- [x] 强大的中间件让你的代码模块化变得非常简单  
- [x] 增加全局上下文， 方便中间件传递值
- [x] 内嵌接口文档
- [x] 支持收尾操作
- [x] 增加数据结构绑定， 适合中间件传递
- [x] 增加websocket， 可以学习，不建议使用
- [x] 集成pprof， router.AddGroup(xmux.Pprof())


### 添加了组的概念

> aritclegroup.article.go
```go
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

var Article *xmux.GroupRoute

func init() {
	Article = xmux.NewGroupRoute()
	Article.Pattern("/{int:id}").Get(hello)

}
```
> main.go
```
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
	router.Pattern("/{all:age}").Get(Who)   // 这个可以匹配任何路由
}
```

记住， 是100%，  此路由优先匹配完全匹配规则， 匹配不到再寻找 正则匹配， 加快了寻址速度  
访问 /get -> 返回 show   
访问  /post   -> 返回 Who  

### 自动检测重复项,
```go
func main() {
	router := xmux.NewRouter()
	router.Pattern("/get").Get(show) // 不同请求分别处理
	router.Pattern("/get").Get(show) // 不同请求分别处理

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
	router.Pattern("/get").Get(show) // 不同请求分别处理
	router.Pattern("/get/").Get(show) // 不同请求分别处理

}

所以运行上面将会报错，/get/被转为 /get 如下  
2019/11/29 21:51:11 pattern duplicate for /get

```

### 看到上面的代码可能想到了什么，  上面如果忘了 设置请求怎么办， 原来不写会语法报错
```go
func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 
	router.Pattern("/get").Get(show) // 不同请求分别处理

	router.AddGroup(aritclegroup.Article)

	router.Pattern("/{string:age}").Get(Who).SetHeader("Host", "two")
	router.Pattern("/home/id").SetHeader("Host", "two")
	log.Fatal(http.ListenAndServe(":8080", router))
}
```
上面的路由  /home/id  没写方法  

当get请求的时候， 页面返回了这个. 应该就明白了  
<h1>when you see this page, it means you forget set handle in /home/id</h1>  

### 四大全局handle
```go
Options:        options(),   //这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
HandleNotFound: handleNotFound(),   // 404 返回
MethodNotFound: methodNotFound(),   // 这个就是上面提示的忘了写handle 的提示页面
MethodNotAllowed methodNotAllowed(),

// 默认调用的方法如下
func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	})
}

func options() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>when you see this page, it means you forget set handle in " + r.URL.Path + "<h1>"))
		return
	})
}

func methodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
}

```

###  header 和 中间件 func(w http.ResponseWriter, r *http.Request)  bool  
路由有3种路由头  
全局路由： 所有请求都会带上这个  
组路由： 所有组的路由都会带上这个， 还有带上全局的， 组的请求头覆盖全局的  
私有路由： 单一路由的请求头， 属于某组的话， 带上组路由头， 全局的话带上全局的  
优先级  
私有路由 > 组路由 > 全局路由  (如果存在优先级大的就覆盖优先级小的)

中间类似上面header   
最后一个返回值是表示是否直接返回， 如果直接返回，后面的中间件和方法将不会执行   
例如下面的方法， 执行的话， 将不会打印66666    
但是如果 hf 返回false   
那么将打印  
44444444444444444444444444  
66666  
cander  
页面将看到  
hello world home  
```
go test -v -run TestHome github.com/hyahm/xmux/example
```       
不过优先级相反    
私有路由 < 组路由 < 全局路由    (如果存在优先级大的就覆盖优先级小的)
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
	router.Pattern("/home/{test}").Get(home).AddMidware(hf).SetHeader("name", "cander").AddMidware(hf1)
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

### 全局 数据  GetData(r) , 对基础路由的数据补充， 用法后面详细补充

```go
func filter(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")
	r.Header.Set("bbb", "ccc")
	return false
}

func name(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world name"))
	return
}
router.Pattern("/aaa/{name}").Get(name).AddMidware(filter).AddMidware(login)
```


### 获取正则匹配的参数
```go
func Who(w http.ResponseWriter, r *http.Request) {
fmt.Println(xmux.Var[r.URL.Path]["name"])
fmt.Println(xmux.Var[r.URL.Path]["age"])
w.Write([]byte("yes is mine"))
return
}

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
如上面代码所示  
  
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

### 收尾操作
用来处理一些无关紧要的收尾处理，比如日志和消息通知，  处理的时候， 与客户端的连接已经断开了， 只有单路由(Route)才支持
```go
这是一个单路由实例（end 是xmux.GetData(r).End  的值）
func end(end interface{}) {
	fmt.Println("-----------------------")
	fmt.Println(end)  // fmt.Println(xmux.GetDate(r).End) 和这个是一个东西
	fmt.Println("end function ")

}

func all(w http.ResponseWriter, r *http.Request) {
	xmux.GetData(r).End = "13333"
	w.Write([]byte("hello world all"))
	return
}
router.Pattern("/bbb/ccc").Get(all).End(end)

```
上面的路由请求后 , 客户端收到 hello world all, 最后会打印 
```
-----------------------
13333
end function 
```
### 同一个路由各组件中通讯 （当前连接断开后，下面的数据会被清空）
```vim
xmux.GetData(r).Data  // 这里对应的是Bind 方法绑定的数据
xmux.GetData(r).Set(k string, v interface{})
xmux.GetData(r).Get(k string) (v interface{})
xmux.GetData(r).Del(k string)
```

### 编写接口文档， 
> 使用接口文档,  第一个参数是组路由名， 第二个参数是挂载的路由uri
== 组路由里面的静态文件 默认挂在 /-/css/xxx.css 和  /-/js/xxx.js 下 ==
== 动态路由   /-/api/{int}.html  ==

```go
// 所有的文档相关的方法都以Api开头， 文档只支持单路由的单请求方式， 多请求方式会乱, 调用的时候只会显示到当前位置以上的路由
router := xmux.NewRouter()
api := xmux.ShowApi("/doc", router)
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
canderdeAir:xmux cander$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: xmux
BenchmarkMux-6                           2913423               397 ns/op             128 B/op          4 allocs/op
BenchmarkMuxAlternativeInRegexp-6        1489166               799 ns/op             224 B/op          8 allocs/op
BenchmarkManyPathVariables-6              802186              1404 ns/op             841 B/op         10 allocs/op
PASS
ok      xmux    5.674s
```


### exmaple下面的例子

example.go
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"xmux"
	"xmux/example/aritclegroup"
)

func show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("show me!!!!"))
	return
}

func postme(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post me!!!!"))
	return
}

func Who(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var(r)["name"])
	fmt.Println(xmux.Var(r)["age"])
	w.Write([]byte("yes is mine"))
	return
}

// 默认已经是这样的了，  如果有其他的请自定义
func Options() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理

	router.AddGroup(aritclegroup.Article())

	router.Pattern("/people/{string:name}/{int:age}").Get(Who).SetHeader("Host", "two")
	log.Fatal(http.ListenAndServe(":8080", router))
}


```
articlegroup/route.go
```go
package aritclegroup

import (
	"fmt"
	"net/http"
	"xmux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var(r)["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

func Article() *xmux.GroupRoute {
	article := xmux.NewGroupRoute()
	article.Pattern("/{int:id}").Get(hello)
	return article
}



```

###### new   

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hyahm/xmux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var(r)["test"])
	fmt.Println(xmux.GetData(r).Data)
	w.Write([]byte("hello world home"))
	return
}

func name(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	fmt.Println(xmux.Var(r)["name"])

	w.Write([]byte("hello world name"))
	return
}

func me(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["me"])
	w.Write([]byte("hello world me"))
	return
}

func all(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello world all"))
	return
}

func login(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("not found data"))
		return true
	}
	err = json.Unmarshal(b, xmux.GetData(r).Data)

	fmt.Println(xmux.GetData(r).Data)
	return false
}

func filter(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("-----------------------")
	fmt.Println("login filter")
	r.Header.Set("bbb", "ccc")

	return false
}

type aaa struct {
	Age int
}

type bbb struct {
	Name string
}

type Home struct {
	Addr   string `json:"addr" type:"string" need:"是" default:"深圳" information:"家庭住址"`
	People int    `json:"people" type:"int" need:"是" default:"1" information:"有多少个人"`
}

type Call struct {
	Code int    `json:"code" type:"int" information:"错误返回码"`
	Msg  string `json:"msg" type:"string" information:"错误信息"`
}

func main() {

	router := xmux.NewRouter()
	router.IgnoreIco = true
	// fmt.Println(router.Slash)
	router.AddMidware(filter)
	router.Pattern("/home").Post(home).ApiDescribe("这是home接口的测试").ApiCreateGroup("home","home page", "home").
		ApiReqHeader("content-type": "application/json").
		ApiReqStruct(&Home{}).
		ApiRequestTemplate(`{"addr": "shenzhen", "people": 5}`).
		ApiResStruct(Call{}).
		ApiResponseTemplate(`{"code": 0, "msg": ""}`).
		ApiSupplement("这个是接口的说明补充， 没补充就不填").Bind(&Home{}).AddMidware(login).Get(home)
	router.Pattern("/aaa/{name}").Post(name).DelMidware(filter).Get(name)
	router.Pattern("/aaa/bbbb/{path:me}").Post(me)
	router.Pattern("/bbb/ccc/{int:oid}/{string:all}").Get(all)

	router.ShowApi("/doc") // 开启文档， 一般都是写在路由的最后, 后面的api不会显示
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatal(err)
	}

}


```
运行上面的代码， 打开localhost:9000/doc, 会看到下面的api

![api.png](http://download.hyahm.com/api.png)




