# xmux
基于原生net.http 极简高灵活性 专注前后端分离项目的路由   
功能自己做主  


[视频教程](https://www.bilibili.com/video/BV1Ji4y1D7o3/)

简体中文 | [English](./README.md) | [简体中文](./README_zh.md) 
### 环境条件
如果想使用json/v2    `go get github.com/hyahm/xmux@jsonv2`   

go >= 1.23

### 导航
- [安装](#install)
- [快速开始](#start)
- [http3](#http3)
- [Using GET, POST, PUT, PATCH, DELETE and OPTIONS](#method)
- [路由组](#group)
- [前缀](#prefix)
- [自动检测重复项](#check)
- [上下文传值](#variable)
- [自动忽略斜杠](#slash)
- [模块](#module)
- [钩子函数](#hook)
- [设置请求头](#header)
- [数据绑定](#bind)
- [url正则匹配](#regex)
- [websocket](#websocket)
- [权限控制模块](#permission)
- [缓存模块](#cache)
- [pprof组](#pprof)
- [swagger组](#swagger)
- [连接的实例](#instance)
- [文件浏览](#browse)
- [限流](#limit)


### 安装<a id="install"></a>  
```
go get github.com/hyahm/xmux
```

### 快速开始<a id="start"></a>  

```go
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

打开 http://localhost:8080 就能看到 hello world!

### http3 <a id="http3"></a>  

```go
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
	xmux.GenerateCertificate("cert.pem", "key.pem", "localhost")
	err := router.RunQuic("cert.pem", "key.pem")
}
```
client.go
```go
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	client := http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 仅测试用
			},
		},
	}

	resp, err := client.Get("https://localhost:8080/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

```
运行 `go run client.go` 就能看到 hello world!

### 请求方式<a id="method"></a>
```go
package main

import (
	"net/http"

	"github.com/hyahm/xmux"
)

func main() {
	router := xmux.NewRouter()
	// 只是例子不建议下面的写法， 而是使用   router.Reqeust("/",nil, "POST", "GET")
	router.Get("/",nil)  // get请求
	router.Post("/",nil)  // post请求
	router.Request("/getpost",nil, "POST", "GET")  // 同时支持get，post请求
	router.Any("/any",nil)  // 支持除了options 之外的所有请求
	router.Run()
}



2019/11/29 21:51:11 Found that the / has multiple request methods. Please use Request method to merge the processing
```

### <a id="group">路由组</a>

> aritclegroup.go
```go
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(xmux.Var(r)["id"])
	w.Write([]byte("hello world!!!!"))
	return
}

var Article *xmux.RouteGroup

func init() {
	Article = xmux.NewRouteGroup()
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

### 前缀 <a id="prefix"></a>
// 注意： DelPrefix和 Prefix 只有路由组才会生效, 直接 再 router 里面无法生效
```go

func main() {
	router := xmux.NewRouter().Prefix("test")
	router.Get("/bbb", c)   // /test/bbb
	router.Get("/ccc", c).DelPrefix("test")   // /test/bbb
	g := xmux.NewRouteGroup()
	g.Get("/aaa", noCache).DelModule(setKey) // /test/bbb
	g.Get("/no/cache1", noCache1).DelModule(setKey).DelPrefix("test") // /no/cache1
	router.AddGroup(g)
	router.Run()
}


```

### 自动检测重复项 <a id="check"></a>
```go
func main() {
	router := xmux.NewRouter()
	router.Get("/get",show) // 不同请求分别处理
	router.Get("/get",nil) // 不同请求分别处理
	router.Run()
}
写一大堆路由，  有没有重复的都不知道  
运行上面将会报错， 如下  
2019/11/29 21:51:11 GET pattern duplicate for /get

```



###  忽略url斜杠 <a id="slash"></a>

将任意多余的斜杠去掉例如
/asdf/sadf//asdfsadf/asdfsdaf////as///, 转为-》 /asdf/sadf/asdfsadf/asdfsdaf/as

```go
func main() {
	router := xmux.NewRouter()
	router.IgnoreSlash = true
	router.Get("/get",show) // 不同请求分别处理
	router.Get("/get/",show) // 不同请求分别处理
	router.Run()
}

如果 router.IgnoreSlash = false
那么运行上面将会报错，/get/被转为 /get 如下  
2019/11/29 21:51:11 pattern duplicate for /get

```


### 三大全局handle
```go
HandleOptions:        handleoptions(),   //这个是全局的options 请求处理， 前端预请求免除每次都要写个预请求的处理, 默认会返回ok， 也可以自定义
HandleNotFound: 	  handleNotFound(),   // 默认返回404 ， 也可以自定义
HanleFavicon：        methodNotAllowed(),    // 默认请求 favicon

// 默认调用的方法如下， 没有找到路由
func handleNotFound(w http.ResponseWriter, r *http.Request)  {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// 注意这一行， 这是为了将状态码传递到exit中打印状态码
	GetInstance(r).StatusCode = http.StatusNotFound
	w.WriteHeader(http.StatusNotFound)
}

```

###  模块 （代替其他框架的中间件功能，并且更灵活更简单）<a id="module"></a>
**核心理念： 任何逻辑都可以是一个模块， 最后组合而成的就是一个完整接口功能**

- 模块类 优先级   
  全局路由 > 组路由 > 私有路由 (如果存在优先级大先执行，。 
  如果不想用可以在不想使用的路由点或路由组 DelModule 来单独删除)

> 模块返回值的含义： true： 直接返回给客户端，不做后面的处理， false就是继续向下执行


```go
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world home"))
	return
}


func hf(w http.ResponseWriter, r *http.Request)  bool {
	fmt.Println("44444444444444444444444444")
	return true
}

func hf1(w http.ResponseWriter, r *http.Request)  bool {
	fmt.Println("66666")
	return false
}

func main() {
	router := xmux.NewRouter().AddModule(hf).SetHeader("name", "cander")
	router.Get("/home/{test}",home).AddModule(hf1)  // 此处会先执行 hf -> hf1 -> home
	router.Get("/test/{test}",home).DelModule(hf)  // 此处直接执行 home
	router.Run()
}

```


###  上下文传值<a id="variable"></a>

- 其中，这几个是内置的，
  -   xmux.GetInstance(r).GetConnectId() ：             连接的id（任何地方都可以使用）
  -   xmux.GetInstance(r).GetCurrFuncName()：       	它的值永远是处理函数的函数名（从模块开始才有值） 
  -   xmux.GetInstance(r).CacheKey:       				缓存的 key
  -   xmux.GetInstance(r).Body:               			获取的body内容  （[]byte）
  -   xmux.GetInstance(r).GetPageKeys()：               跟页面权限有关（从模块开始才有值）
  -   xmux.GetInstance(r).StatusCode：            		接口返回的状态码（有些情况要修改，比如页面跳转，任何地方都可以使用）
  -   xmux.GetInstance(r).Data:            		        数据绑定
-  自定义的值是从模块开始才能赋值

```
# 设置值
xmux.GetInstance(r).Set("key", "value")

# 获取值
xmux.GetInstance(r).Get("key")
```

```go
func hf1(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("login mw")
	xmux.GetInstance(r).Set("name","xmux")
	r.Header.Set("bbb", "ccc")
	return false
}

func hf(w http.ResponseWriter, r *http.Request) bool {
	name := xmux.GetInstance(r).Get("name").(string)
    fmt.Println(name)
	w.Write([]byte("hello world " + name))
	return false
}

func home(w http.ResponseWriter, r *http.Request)  {
    fmt.Println(name)
}

func main() {
	router := xmux.NewRouter().AddModule(hf).SetHeader("name", "cander")
	router.Exit = exit
	router.Enter = enter
	router.Get("/home/{test}",home).AddModule(hf1)  // 此处会先执行 hf -> hf1 -> home
	router.Get("/test/{test}",home).DelModule(hf)  // 此处直接执行 home
	router.Run()
}


```

> curl http://localhost:8080/aaa/name

```
login mw
name
```


### 钩子<a id="hook"></a>

- NotFoundRequiredField                                             : 必要字段验证失败的处理勾子
- UnmarshalError                                                   : 内置解析解析错误的勾子
- Exit (start time.Time, w http.ResponseWriter, r *http.Request)   :    // 匹配到的路由才会进来 
- Enter( w http.ResponseWriter, r *http.Request) bool              :    // 匹配到的路由才会进来 
- HandleAll(w http.ResponseWriter,r *http.Request)   bool            :     为了性能考虑新增   所有请求都能在这里获取到数据， 用来替代之前的 enter 和  exit 的请求记录
```go

func exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	//  主要为了打印执行的时间  注意此时间没有计算寻址时间， 本路径寻址使用了缓存，可以无视寻址时间
	 // 匹配到的路由才会进来 
	fmt.Println(time.Since(start).Seconds(), r.URL.Path)
}

// 与module一样的效果， return true 就是直接返回， return false 就是继续 但是不支持 xmux.GetInstence(r)传参  
// 主要用来过滤请求和调试
func enter( w http.ResponseWriter, r *http.Request) bool {
 // 匹配到的路由才会进来 
	fmt.Println(time.Since(start).Seconds(), r.URL.Path)
}

func HandleAll( w http.ResponseWriter, r *http.Request) bool {
	// 任何请求都会进入到这里，比如过滤ip， 域名
 // 匹配到的路由才会进来 
	fmt.Println(time.Since(start).Seconds(), r.URL.Path)
}
```

### 设置请求头 <a id="header"></a>
跨域主要是添加请求头的问题, 其余框架一般都是借助中间件来设置   
但是本路由借助上面请求头设置 大大简化跨域配置  

优先级  
私有路由 > 组路由 > 全局路由  (如果存在优先级大的就覆盖优先级小的)

```go
// 跨域处理的例子， 设置下面的请求头后， 所有路由都将挂载上请求头， 
// 如果某些路由有单独请求头， 可以单独设置
func main() {
	router := xmux.NewRouter()
	router.IgnoreSlash = true
	router.SetHeader("Access-Control-Allow-Origin", "*")  // 主要的解决跨域, 因为是全局的请求头， 所以后面增加的路由全部支持跨域
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type,Access-Token,X-Token,Origin,smail,authorization")  // 新增加的请求头
	router.Get("/", index)
	router.Run()
}
```



### 数据绑定（绑定解析后的数据从模块之前开始生效）<a id="bind"></a>

- BindJson:       绑定的是一个json
- BindXml：     绑定是一个xml
- BindForm：  绑定的是一个form, file文件无法获取， 需要额外通过 r.FormFile() 获取
- Bind：           自定义处理绑定（通过模块来进行处理）



验证字段的必须存在

```go
type User struct {
	Username string `json:"username,required" form:"username,required"`
}
```

router.PrintRequestStr：  是否打印接受请求体内容

xmux.MaxPrintLength： 打印的form表单的大小如果超过指定的大小就不打印（默认2k）

```go
func JsonToStruct(w http.ResponseWriter, r *http.Request) bool {
	// 任何报错信息， 直接return true， 就是此handle 直接执行完毕了， 不继续向后面走了
	if goconfig.ReadBool("debug", false) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return true
		}
		err = json.Unmarshal(b, xmux.GetInstance(r).Data)
		if err != nil {
			return true
		}
	} else {
		err := json.NewDecoder(r.Body).Decode(xmux.GetInstance(r).Data)
		if err != nil {
			return true
		}

	}
	return false
}

type DataName struct{}
type DataStd struct{}
type DataFoo struct{}

func AddName(w http.ResponseWriter, r *http.Request) {
	df := xmux.GetInstance(r).Data.(*DataName)
	fmt.Printf("%#v", df)
}

func AddStd(w http.ResponseWriter, r *http.Request) {
	df := xmux.GetInstance(r).Data.(*DataStd)
	fmt.Printf("%#v", df)
}

func AddFoo(w http.ResponseWriter, r *http.Request) {
	df := xmux.GetInstance(r).Data.(*DataFoo)
	fmt.Printf("%#v", df)
}

func main() {
	router := xmux.NewRouter()
	router.Post("/important/name", AddName).Bind(&DataName{}).AddModule(JsonToStruct)
	router.Post("/important/std", AddStd).Bind(&DataStd{}).AddModule(JsonToStruct)
	router.Post("/important/foo", AddFoo).Bind(&DataFoo{}).AddModule(JsonToStruct)
	// 也可以直接使用内置的
	router.Post("/important/foo/by/json", AddFoo).BindJson(&DataFoo{}) // 如果是json格式的可以直接 BindJson 与上面是类似的效果
	router.Run()
}

```

- 绑定返回值

```
    data := &Response{
		Code: 200,
	}
	router := xmux.NewRouter().BindResponse(data)
```

 通过.BindResponse(nil) 来设置取消使用全局绑定  



### 匹配路由 <a id="regex"></a>  
支持以下5种   
 word   只匹配数字和字母下划线（默认）    
 string  匹配所有不含/的字符    
 int  匹配整数    
 all： 匹配所有的包括/   
 re： 自定义正则   
/aaa/{name}          这个和下面一个一样， 省略类型， 默认是word  
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


### websocket <a id="websocket"></a>
下面是一个完整的例子
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hyahm/xmux"
)

type client struct {
	msg string
	c   *xmux.BaseWs
}

var msgchan chan client
var wsmu sync.RWMutex
var ps map[*xmux.BaseWs]byte

func sendMsg() {
	for {
		c := <-msgchan
		for p := range ps {
			if c.c == p {
				// 不发给自己
				continue
			}
			fmt.Println(c.msg)
			// 发送的msg的长度不能超过 1<<31, 否则掉内容， 建议分包
			p.SendMessage([]byte(c.msg), ps[p])
		}
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	p, err := xmux.UpgradeWebSocket(w, r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	p.SendMessage([]byte("hello"), xmux.TypeMsg)
	wsmu.Lock()
	ps[p] = xmux.TypeMsg
	wsmu.Unlock()
	tt := time.NewTicker(time.Second * 2)
	go func() {
		for {
			<-tt.C
			if err := p.SendMessage([]byte(time.Now().String()), xmux.TypeMsg); err != nil {
				break
			}
		}
	}()
	for {
		if p.Conn == nil {
			return
		}
		// 封包
		msgType, msg, err := p.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			// 连接断开
			wsmu.Lock()
			delete(ps, p)
			wsmu.Unlock()
			break
		}
		ps[p] = msgType
		c := client{
			msg: msg,
			c:   p,
		}
		msgchan <- c
	}
}

func main() {
	router := xmux.NewRouter()
	wsmu = sync.RWMutex{}
	msgchan = make(chan client, 100)
	ps = make(map[*xmux.BaseWs]byte)
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.Get("/{int:uid}", ws)

	go sendMsg()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}

```

```html
<!DOCTYPE html>
<html>

<head>
    <title>go websocket</title>
    <meta charset="utf-8" />
</head>

<body>
    <script type="text/javascript">
        var wsUri = "ws://localhost:8080/3";
        var output;
        var connect = false;
   
        function init() {
            output = document.getElementById("output");
            testWebSocket();
        }

        function testWebSocket() {
            websocket = new WebSocket(wsUri, WebSocket.binaryType);
            websocket.onopen = function(evt) {
                onOpen(evt)
            };
            websocket.onclose = function(evt) {
                onClose(evt)
            };
            websocket.onmessage = function(evt) {
                onMessage(evt)
            };
            websocket.onerror = function(evt) {
                onError(evt)
            };
        }

        function onOpen(evt) {
            writeToScreen("CONNECTED");
            connect = true
                // doSend("WebSocket rocks");
        }

        function onClose(evt) {
            connect = false
            writeToScreen("DISCONNECTED");
        }

        function onMessage(evt) {

            msg = String.fromCharCode(evt.data)
            console.log(msg)
            writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data + '</span>');
            // websocket.close();
        }

        function onError(evt) {
            writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
        }

        function doSend(message) {
            if (!connect) {
                console.log("connect error")
                return
            }
            writeToScreen("SENT: " + message);
            websocket.send(message);
        }

        function writeToScreen(message) {
            var pre = document.createElement("p");
            pre.style.wordWrap = "break-word";

            pre.innerHTML = message;
            output.appendChild(pre);
        }

        window.addEventListener("load", init, false);

        function sendBtnClick() {
            var msg = document.getElementById("input").value;
            doSend(msg);
            document.getElementById("input").value = '';
        }

        function closeBtnClick() {
            websocket.close();
        }
    </script>
    <h2>WebSocket Test</h2>
    <input type="text" id="input"></input>
    <button onclick="sendBtnClick()">send</button>
    <button onclick="closeBtnClick()">close</button>
    <div id="output"></div>

</body>

</html>
```

### 集成swagger <a id="swagger"></a>

> 适用的函数全部以Swagger开头， 数据结构与[swagger](https://swagger.io/docs/specification/2-0/basic-structure/  )文档的一样，请swagger数据结构  
    
```
package main

import "github.com/hyahm/xmux"

func main() {
	router := xmux.NewRouter()
	router.Get("/", nil)
	router.AddGroup(router.ShowSwagger("/docs", "localhost:8080"))
	router.Run()
}

```

### 获取当前的连接数
```go
xmux.GetConnents()
```

### 优雅的停止
```
xmux.StopService()
```


### 内置路由缓存
```go

xmux.NewRouter(cache ...uint64) // cache 是一个内置lru 路径缓存， 不写默认缓存10000， 请根据情况自己修改
```

###  权限控制 <a id="permission"></a>
- 页面权限
  思路来自前端框架路由组件 meta 的 roles  
  通过给定数组来判断

  > 以github前端star最多的vue后端项目为例子    https://github.com/PanJiaChen/vue-element-admin

  

  ```
  src/router/index.js 里面的页面权限路由
  
  
  {
      path: '/permission',
      component: Layout,
      redirect: '/permission/page',
      alwaysShow: true, // will always show the root menu
      name: 'Permission',
      meta: {
        title: 'Permission',
        icon: 'lock',
        roles: ['admin', 'editor'] // you can set roles in root nav
      },
      children: [
        {
          path: 'page',
          component: () => import('@/views/permission/page'),
          name: 'PagePermission',
          meta: {
            title: 'Page Permission',
            roles: ['admin'] // or you can only set roles in sub nav
          }
        },
        {
          path: 'directive',
          component: () => import('@/views/permission/directive'),
          name: 'DirectivePermission',
          meta: {
            title: 'Directive Permission'
            // if do not set roles, means: this page does not required permission
          }
        },
        {
          path: 'role',
          component: () => import('@/views/permission/role'),
          name: 'RolePermission',
          meta: {
            title: 'Role Permission',
            roles: ['admin']
          }
        }
      ]
    },
  ```

  > xmux 对应的写法

  ```
  
  func AddName(w http.ResponseWriter, r *http.Request) {
  	fmt.Printf("%v", "AddName")
  }
  
  func AddStd(w http.ResponseWriter, r *http.Request) {
  	fmt.Printf("%v", "AddStd")
  }
  
  func AddFoo(w http.ResponseWriter, r *http.Request) {
  	fmt.Printf("%v", "AddFoo")
  }
  
  func role(w http.ResponseWriter, r *http.Request) {
  	fmt.Printf("%v", "role")
  }
  
  func DefaultPermissionTemplate(w http.ResponseWriter, r *http.Request) (post bool) {
  
  	// 拿到对应uri的权限， 也就是AddPageKeys和DelPageKeys所设置的
  	pages := xmux.GetInstance(r).GetPageKeys()
  	// 如果长度为0的话，说明任何人都可以访问
  	if len(pages) == 0 {
  		return false
  	}
  
  	// 拿到用户对应的 role，判断是都在
  	roles := []string{"admin"} //从数据库中获取或redis获取用户的权限
  	for _, role := range roles {
  		if _, ok := pages[role]; ok {
  			// 这里匹配的是存在这个权限， 那么久继续往后面的走
  			return false
  		}
  	}
  	// 没有权限
  	w.Write([]byte("no permission"))
  	return true
  }
  
  func 
  
  func main() {
  	router := xmux.NewRouter()
	// 添加验证模块，直接用模版就可以，也可以自己写
	router.AddModule(xmux.DefaultPermissionTemplate)
	// AddPageKeys() 里面的字符串权限与js 路由里面的验证是一样的， 这里是全局添加， 所有下面所有路由都有
  	router.AddPageKeys("admin", "editor")
  	router.Post("/permission", AddName)
  	router.Post("/permission/page", AddStd).DelPageKeys("editor")
  	router.Post("/permission/directive", AddFoo)
  	// 也可以直接使用内置的
  	router.Post("/permission/role", role).DelPageKeys("editor")
  	router.Run()
  }
  
  ```

  

- 更加细致的增删改查权限但不限于 增删改查
  想过最简单的是根据 handle 的函数名 来判断， 

  可以参考xmux的权限模板  xmux.DefaultPermissionTemplate


### 缓存 <a id="cache"></a>

- 初始化缓存  cache.InitResponseCache()
- 需要设置缓存的 key 的模块（核心的模块， 如果没设置的话， 就不用缓存）
  - 为了设置 CacheKey的值 xmux.GetInstance(r).Set(xmux.CacheKey, fmt.Sprintf("%s_%v", r.URL.Path, uid))
- 需要挂载缓存模块 

> 没有绑定返回数据的例子
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hyahm/gocache"
	"github.com/hyahm/xmux"
)

func c(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comming c")
	now := time.Now().String()
	xmux.GetInstance(r).Response.(*Response).Data = now
}

func noCache(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update c")
	xmux.NeedUpdate("/aaa")
}

func noCache1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comming noCache1")
	now := time.Now().String()
	xmux.GetInstance(r).Response.(*Response).Data = now
}

func setKey(w http.ResponseWriter, r *http.Request) bool {
	xmux.GetInstance(r).Set(xmux.CacheKey, r.URL.Path)
	fmt.Print(r.URL.Path + " is cached")
	return false
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func main() {
	r := &Response{
		Code: 0,
	}
	cth := gocache.NewCache[string, []byte](100, gocache.LFU)
	xmux.InitResponseCache(cth)
	// 如果没有绑定返回， 那么使用 xmux.DefaultCacheTemplateCacheWithoutResponse 代替 xmux.DefaultCacheTemplateCacheWithResponse
	router := xmux.NewRouter().AddModule(setKey, xmux.DefaultCacheTemplateCacheWithResponse)
	router.BindResponse(r)
	router.Get("/aaa", c)
	router.Get("/update/aaa", noCache).DelModule(setKey)
	router.Get("/no/cache1", noCache1).DelModule(setKey)
	router.Run()
}


```


### 客户端文件下载（官方内置方法 mp4文件为例）

```go
func PlayVideo(w http.ResponseWriter, r *http.Request) {
	filename := xmux.Var(r)["filename"]
	f, err := os.Open(<mp4file>)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("X-Download-Options", "noopen")
	http.ServeContent(w, r, filename, time.Now(), f)

```

### 客户端文件上传(官方内置方法)
```go
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// 官方默认上传文件的大小最大是32M， 可以通过方法设置新的大小
	r.ParseMultipartForm(100 << 20)   // 100M
	// 读取文件
	file, header, err := r.FormFile("file")
	if err != nil {
		return
	}
	f, err := os.OpenFile(<storefilepath>, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, err := io.Copy(f, file)
	if err != nil {
		return
	}
}
```

### 连接的实例 <a id="instance"></a>  
```
xmux.GetInstance(r).Body // 请求过来的数据， 只有绑定值了才有这个数据
xmux.GetInstance(r).CacheKey  // 绑定缓存的key
xmux.GetInstance(r).Data   // 数据绑定解析出来的值
xmux.GetInstance(r).Response  // 绑定返回值
xmux.GetInstance(r).StatusCode   // status_code
xmux.GetInstance(r).Get()  // 上下文传值用的
xmux.GetInstance(r).Set()  // 上下文传值用的
xmux.GetInstance(r).GetConnectId()  // 获取当前的链接id
xmux.GetInstance(r).GetFuncName()  // 增删改查权限
xmux.GetInstance(r).GetPageKeys()  // 页面权限
```


# 目录列表 <a id="browse"></a>  
```go
package xmux

import (
	"log"
	"github.com/hyahm/xmux"
)


func main() {
	router := xmux.NewRouter()
	// 第一个参数是 url前缀路径，  第二个参数是本地目录，  第三个是 是否显示列出文件列表， 第四个是是都下载
	router.AddGroup(FileBrowse("/static", "D:\\ProgramData", true, false))
	log.Fatal(router.Run())
}
```

# 限流 <a id="limit"></a>
```go 
// 固定窗口的模板，其他的 滑动窗口计数器，漏桶算法，令牌桶算法  在同一个文件下注释，因为需要定义全部变量，减少没西药的内存消耗也是需要自己做少量修改
router.AddModule(xmux.LimitFixedWindowCounterTemplate)
```

### 性能分析

```
func main() {
	router := xmux.NewRouter()
	router.Post("/", nil)
	// 也可以直接使用内置的
	router.AddGroup(xmux.Pprof())
	router.Run()
}
```



> open http://localhost:8080/debug/pprof  can see pprof page



### 查看某handel详细的中间件模块等信息

```
// 查看某个指定路由的详细信息
router.DebugAssignRoute("/user/info")  
// 查看某个正则火匹配路由的固定uri来获取某路由的详细
router.DebugIncludeTpl("")
// 显示全部的， 不建议使用，
router.DebugRoute()
router.DebugTpl()
```

> out

```
2022/01/22 17:16:11 url: /user/info, method: GET, header: map[], module: xmux.funcOrder{"github.com/hyahm/xmux.DefaultPermissionTemplate"}, midware: "" , pages: map[string]struct {}{}
```



### 流程图总结 （没有匹配路由不会进入下面的图）
![xmux流程图](xmux.jpg)

