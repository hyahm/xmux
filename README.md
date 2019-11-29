# xmux
之前一直使用的 mux， 但是xmux 已经无法满足自己优化代码的需求  

作者在蒙圈研究go test 测试，  估计还有问题， 请勿在生产环境使用

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
	router.AddGroup(aritclegroup.Article())
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

记住， 是100%，  此框架优先匹配完全匹配规则， 匹配不到再寻找 正则匹配， 加快了寻址速度
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

### 自动检测重复项


### 看到上面的代码可能想到了什么，  上面如果忘了 设置请求怎么办， 原来不写会语法报错
```go
func main() {
	router := xmux.NewRouter()
	router.Options = Options()                    // 这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
	router.Pattern("/get").Get(show).Post(postme) // 不同请求分别处理

	router.AddGroup(aritclegroup.Article)

	router.Pattern("/{string:age}").Get(Who).SetHeader("Host", "two")
	router.Pattern("/home/id").SetHeader("Host", "two")
	log.Fatal(http.ListenAndServe(":8080", router))
}
```
上面的路由  /home/id  没写方法

当get请求的


### 添加了请求方式分别处理， 同时支持post, get 等， 不同的处理函数
exmaple下面的例子
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

### 为了加速匹配， 增加路由表概念， 一担添加进去无法修改， 不会过期， 重启会清空

