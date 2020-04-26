package xmux

// html 接口文档的 模板

var tpl = `
<!doctype html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.Title}}</title>
</head>
<body>
{{range .Api}}

<h3>{{.Title}}</h3>
<h5>描述</h5>
<code>{{.Describe}}</code>
<h5>url</h5>
<code>{{.Url}}</code>
<h5>请求方式</h5>
<code>{{.Method}}</code>
<h5>请求头</h5>
<code>
{{range $key, $value := .Header }}
    {{$key}} : {{$value}} <br/>
{{end}}
</code>
<h5>请求参数</h5>

{{end}}
</body>
</html>`
