# xmux， go语言 路由(router)
应该是基于原生net.http包唯一一个带缓存， 使用简单的路由,  目前只有作者本人的项目在使用， 没有发现bug， 有什么问题欢迎反馈

### 已完成功能
- [x] 支持路由分组
- [x] 支持全局请求头， 组请求头， 私有请求头
- [x] 支持自定义method， 多method
- [x] 支持正则匹配和参数获取
- [x] 完全匹配优先于正则匹配
- [x] 自动检查pattern
- [x] 支持修复pattern
- [x] 自定修复请求的url
- [x] 正则匹配支持（int(\d+), word(\w+), re, all(.*?)，不写默认 string([^\/])）建议使用string
- [x] 支持四大全局的handle（notFound, methodNotFound, handleNotFound, Options请求）  
- [x] 支持中间件  
- [x] 增加websocket， 可以学习，不建议使用


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
	Article = xmux.NewGroupRoute("article")
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
	router.Pattern("/{all:age}").Get(Who)   // 这个可以匹配任何路由
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

### 四大全局handle
```go
Options:        options(),   //这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
NotFound:       notFound(),   // 404 返回
HandleNotFound: handleNotFound(),   // 这个就是上面提示的忘了写handle 的提示页面
MethodNotAllowed http.Handler

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

func methodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
}


可以选择自定义的， 只要new路由赋值即可
router := xmux.NewRouter()
router.Options = Options()  
methodNotAllowed 和  handleNotFound的区别
当存在handle 但找不到 method 就返回 methodNotAllowed
不存在handle 就返回 handleNotFound的


```

###  header 和 中间件 func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, bool)  
路由有3种路由头  
全局路由： 所有请求都会带上这个  
组路由： 所有组的路由都会带上这个， 还有带上全局的， 组的请求头覆盖全局的  
私有路由： 单一路由的请求头， 属于某组的话， 带上组路由头， 全局的话带上全局的  
优先级  
私有路由 > 组路由 > 全局路由  

中间类似上面header, 最后一个返回值是表示是否直接返回， 如果直接返回， 后面的方法将不会执行, 例如下面的方法， 执行的话， 将不会打印66666  
```
go test -v -run TestHome github.com/hyahm/xmux/example
```       
不过优先级相反    
私有路由 < 组路由 < 全局路由    
```go
func mid() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("77777")
		return
	})
}

func hf(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request,  bool) {
	fmt.Println("44444444444444444444444444")
	r.Header.Set("name", "cander")
	
	return w, r, true
}

func hf1(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, bool) {
	fmt.Println("66666")
	fmt.Println(r.Header.Get("name"))
	return w, r, true
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
xmux.Var[r.URL.Path]["name"]  // 获取方法
```
后面会增加自定义正则匹配

### 压力测试
```
canderdeAir:xmux cander$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: xmux
BenchmarkMux-4                          22398705                53.5 ns/op             0 B/op          0 allocs/op
BenchmarkMuxAlternativeInRegexp-4       11393886               106 ns/op               0 B/op          0 allocs/op
BenchmarkManyPathVariables-4            10317334               110 ns/op               0 B/op          0 allocs/op
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



