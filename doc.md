安装
要安装 Gin 软件包，需要先安装 Go 并设置 Go 工作区。

1. 下载并安装 gin：

```
go get -u github.com/hyahm/xmux
```
2. 使用 `xmux`：
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

运行你的项目
```go run main.go```

hello world!
