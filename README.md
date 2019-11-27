# xmux
之前一直使用的 mux， 但是xmux 已经无法满足自己优化代码的需求

### 初始阶段， 为了满足自己代码的高封装， 暂时不支持正则匹配路由

### 添加了组的概念

### 添加了请求方式分别处理， 同时支持post, get 等， 不同的处理函数
exmaple下面的例子
example.go
```go
package main

import (
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

// 默认已经是这样的了，  如果有其他的请自定义
func Options() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func main() {
	router := xmux.NewRouter()
	router.Options = Options()                       // 这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理
	router.HandleFunc("/get").Get(show).Post(postme) // 不同请求分别处理
	router.AddGroup(aritclegroup.Article())

	log.Fatal(http.ListenAndServe(":8080", router))
}


```
articlegroup/route.go
```go
package aritclegroup

import (
	"net/http"
	"xmux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!!!!"))
	return
}

func Article() *xmux.GroupRoute {
	article := xmux.NewGroupRoute("/article")
	article.HandleFunc("name").Get(hello)
	return article
}

```
### 因为没有正则， 全部采用map匹配路由， 速度肯定是快速的(后面会增加)
