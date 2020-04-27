package xmux

// html 接口文档的 模板

var tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="stylesheet" href="http://download.hyahm.com/css/style.css"/>
    <link rel="stylesheet" href="http://download.hyahm.com/css/left.css">
    <link rel="stylesheet" type="text/css" href="http://download.hyahm.com/font/iconfont.css">
</head>
<body>
<div class="body-content">
        <h4>一共 {{ .Api | len }} 个路由</h4>
        {{ range .Api }}
        <div class="right-light">
            <div class="right-dl">
                <span class="right-get">{{ .Method }}</span>
                <span class="right-url">{{ .Url }}</span>
            </div>
            <div class="dl-box dl-none">
          
                {{ if .Describe }}
                    <h3>简述</h3>
                    <div class="dl-bz">{{ .Describe }}</div>
                {{ end }}
    
                {{ if .Header }}
                    <h3>请求头</h3>
                    {{ range $k, $v := .Header }}
                    <div class="dl-bz">{{ $k }} : {{ $v }}</div>
                    {{ end }}
                {{ end }}
         
                {{ if .Opt }}
                    <h3>参数</h3>
                    <div class="dl-table">
                        <span>参数名</span>
                        <span>类型</span>
                        <span>必选</span>
                        <span>默认值</span>
                        <span>说明</span>
                    </div>

                    {{ range .Opt }}
                        <div class="dl-table dl-table-msg">
                            <span>{{ .Name }}</span>
                            <span>{{ .Typ }}</span>
                            <span>{{ .Need }}</span>
                            <span>{{ .Default }}</span>
                            <span>{{ .Information }}</span>
                        </div>
                    {{ end }}
                {{ end }}

                {{ if .Request }}
                    <h3>请求示例</h3>
                    <div class="dl-ex-box">
                        <p hidden id="req">{{ .Request }} </p>
                        <pre class="dl-expl" id="json_req">
                        </pre>
                    </div>
                {{ end }}

                {{ if .Callbak }}
                    <h3>返回参数说明</h3>
                    <div class="dl-table dl-table1">
                        <span>参数名</span>
                        <span>类型</span>
                        <span>说明</span>
                    </div>
                    {{ range .Callbak }}
                        <div class="dl-table dl-table1 dl-table-msg  dl-table-msg1">
                            <span>{{ .Name }}</span>
                            <span>{{ .Typ }}</span>
                            <span>{{ .Information }}</span>
                        </div>
                    {{ end }}
                {{ end }}

                {{ if .Response }}
                    <h3>返回示例</h3>
                    <div class="dl-ex-box">
                        <p hidden id="res">{{ .Response }} </p>
                        <pre class="dl-expl" id="json_res">
                        </pre>
                    </div>
                {{ end }}

                {{ if .Supplement }}
                    <h3>备注</h3>
                    <div class="dl-bz">{{ .Supplement }}</div>
                {{ end }}
            </div>
        </div>
        {{ end }}

</div>
</body>
<script type="text/javascript" src="http://download.hyahm.com/js/jquery.min.js"></script>
<script type="text/javascript" src="http://download.hyahm.com/js/jquery.slimscroll.min.js"></script>
<script type="text/javascript" src="http://download.hyahm.com/js/left-side-menu.js"></script>
<script type="text/javascript" src="http://download.hyahm.com/js/click.js"></script>
</html>
`
