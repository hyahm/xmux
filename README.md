# xmux， go语言 路由
之前一直使用的 mux， 但是xmux 已经无法满足自己优化代码的需求  

目前还有不足，本人是小白鼠，  估计还有问题， 请勿在生产环境使用

### 添加了组的概念
几十个路由写一个文件里面， 嗯， 还好， 但是多了呢， 眼睛是不是有点花， 并且某些名字重名了也不知道
 增加组了后， 可以分开写路由了，
###### 第一种
> 函数封装   

aritclegroup.article.go
```go


func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

func Article() *xmux.GroupRoute {
	article := xmux.NewGroupRoute()
	article.Pattern("/{int:id}").Get(hello)
	return article
}
```
> main.go
```go
func main() {
	router := xmux.NewRouter()
	router.AddGroup(aritclegroup.Article())
}
```

###### 第二种
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
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理
	router.Pattern("/{string:age}").Get(Who)
}
```

记住， 是100%，  此路由优先匹配完全匹配规则， 匹配不到再寻找 正则匹配， 加快了寻址速度  
访问 /get -> 返回 show   
访问  /post   -> 返回 Who  

### 如上所示， 多methods 分开使用
```go
func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理
	router.Pattern("/{string:age}").Get(Who)
}
```
/get get 请求 返回 show，    post 请求 返回 postme

### 自动检测重复项,
```go
func main() {
	router := xmux.NewRouter()
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理

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
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理
	router.Pattern("/get/").Get(show).Post(postme) // 不同请求分别处理

}

所以运行上面将会报错，/get/被转为 /get 如下  
2019/11/29 21:51:11 pattern duplicate for /get

```

### 看到上面的代码可能想到了什么，  上面如果忘了 设置请求怎么办， 原来不写会语法报错
```go
func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理

	router.AddGroup(aritclegroup.Article)

	router.Pattern("/{string:age}").Get(Who).SetHeader("Host", "two")
	router.Pattern("/home/id").SetHeader("Host", "two")
	log.Fatal(http.ListenAndServe(":8080", router))
}
```
上面的路由  /home/id  没写方法  

当get请求的时候， 页面返回了这个. 应该就明白了  
<h1>when you see this page, it means you forget set handle in /home/id</h1>  

### 三大全局handle
```go
Options:        options(),   //这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
NotFound:       notFound(),   // 404 返回
HandleNotFound: handleNotFound(),   // 这个就是上面提示的忘了写handle 的提示页面

// 默认调用的方法如下
func notFound() http.Handler {
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

可以选择自定义的， 只要new路由赋值即可
router := xmux.NewRouter()
router.Options = Options()  

```

### 嗯， 注意的人应该注意了，  上面有header  
路由有3种路由头  
全局路由： 所有请求都会带上这个  
组路由： 所有组的路由都会带上这个， 还有带上全局的， 组的请求头覆盖全局的  
私有路由： 单一路由的请求头， 属于某组的话， 带上组路由头， 全局的话带上全局的  
优先级  
私有路由 > 组路由 > 全局路由  
会覆盖， 后面会补充删除的， 不提供add头， 不然又复杂了，  


### 获取正则匹配的参数
```go
func Who(w http.ResponseWriter, r *http.Request) {
fmt.Println(xmux.Var[r.URL.Path]["name"])
fmt.Println(xmux.Var[r.URL.Path]["age"])
w.Write([]byte("yes is mine"))
return
}

```
有没有一种坑爹的冲动， 加上这个也是为了高并发使用，引入了路由表而这么做的  
多写个[r.URL.Path]就好了， 实在是有不得不放弃简化的理由  

### 路由表，  嗯， 就是这个东西了
只有几个路由是看不到效果的， 成千上万个路由， 虽然也不会有那么多， 优势越大， 其实就是路由缓存  
第一次匹配到某路由的handle后， 下次不用寻址了， 直接从路由表获取即可  
```go
type rt struct {
	Handle http.Handler
	Header map[string]string
}

routeTable     map[string]*rt   // 就是这个东西了， 保存了handle和请求头信息
```

### 匹配路由
如上面代码所示
/aaa/{name}          这个和下面一个一样， 省略类型， 默认是string
/aaa/{string:name}   这个和上面一样， string类型
/aaa/{int:name}       这个匹配int类型
```
xmux.Var[r.URL.Path]["name"]  // 获取方法
```
后面会增加自定义正则匹配

### 看看速度对比吧
里面有个bench_test.go 文件  
从mux里面来的  
本框架的压力测试数据  
```
canderdeMacBook-Air:xmux cander$ go test -bench=.
goos: darwin
goarch: amd64
pkg: xmux
BenchmarkMux-4                          21019719                52.3 ns/op
BenchmarkMuxAlternativeInRegexp-4       11333706               105 ns/op
BenchmarkManyPathVariables-4            10704848               106 ns/op
PASS
ok      xmux    4.993s
canderdeMacBook-Air:xmux cander$ 
```
mux 框架的, 他的框架更新了， 注释掉空函数  
```go
canderdeMacBook-Air:mux cander$ go test -bench=.
goos: darwin
goarch: amd64
pkg: mux
BenchmarkMux-4                            689672              1835 ns/op
BenchmarkMuxAlternativeInRegexp-4         461038              2579 ns/op
BenchmarkManyPathVariables-4              453232              2604 ns/op
PASS
ok      mux     7.402s
canderdeMacBook-Air:mux cander$ 

```
嗯， 不比不知道， 一比吓一跳，20倍以上的速度， 不知道是寻址的问题还是路由表的功劳  

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
	fmt.Println(xmux.Var[r.URL.Path]["name"])
	fmt.Println(xmux.Var[r.URL.Path]["age"])
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
	fmt.Println(xmux.Var[r.URL.Path]["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

func Article() *xmux.GroupRoute {
	article := xmux.NewGroupRoute()
	article.Pattern("/{int:id}").Get(hello)
	return article
}



```



