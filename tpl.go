package xmux

// html 接口文档的 模板

var tpl0 = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="stylesheet" href="/-/css/style.css"/>
    <link rel="stylesheet" href="/-/css/left.css">
    <link rel="stylesheet" href="/-/css/font.css">
</head>
<body>
<div class="body-content">
    <div class="body-head">api接口文档</div>
        <div class="left-side-menu" >
            <div class="left-search"><i class=" iconfont icon-cc-search"></i><input type="search" value="" class="input" placeholder="关键字"/></div>
            <div class="lsm-expand-btn">
                <div class="lsm-mini-btn">
                    <label>
                        <input type="checkbox" checked="checked">
                        <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                            <circle cx="50" cy="50" r="30" />
                            <path class="line--1" d="M0 40h62c18 0 18-20-17 5L31 55" />
                            <path class="line--2" d="M0 50h80" />
                            <path class="line--3" d="M0 60h62c18 0 18 20-17-5L31 45" />
                        </svg>
                    </label>

                </div>
            </div>
            <div class="lsm-container">
                <div class="lsm-scroll" >
                    <div class="lsm-sidebar">
                        <ul>
                            <li class="lsm-sidebar-item">
                                {{ range $k,$v := .Sidebar }}
                                <a href="{{ $k }}"><i class="iconfont lsm-sidebar-icon icon_1"></i><span>{{ $v }}</span><i class="iconfont lsm-sidebar-more"></i></a>
                                {{ end }}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

    </div>
    <div class="body-right">
        <h2 class="text-center">{{ .Title }}</h2>
       
    </div>
</div>
</body>

<script type="text/javascript" src="/-/js/jquery.js"></script>
<script type="text/javascript" src="/-/js/slimscroll.js"></script>
<script type="text/javascript" src="/-/js/left.js"></script>
<script type="text/javascript" src="/-/js/click.js"></script>

</html>
`

var tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="stylesheet" href="/-/css/style.css"/>
    <link rel="stylesheet" href="/-/css/left.css">
    <link rel="stylesheet" href="/-/css/font.css">
</head>
<body>
<div class="body-content">
    <div class="body-head">api接口文档</div>
        <div class="left-side-menu" >
            <div class="left-search"><i class=" iconfont icon-cc-search"></i><input type="search" value="" class="input" placeholder="关键字"/></div>
            <div class="lsm-expand-btn">
                <div class="lsm-mini-btn">
                    <label>
                        <input type="checkbox" checked="checked">
                        <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                            <circle cx="50" cy="50" r="30" />
                            <path class="line--1" d="M0 40h62c18 0 18-20-17 5L31 55" />
                            <path class="line--2" d="M0 50h80" />
                            <path class="line--3" d="M0 60h62c18 0 18 20-17-5L31 45" />
                        </svg>
                    </label>

                </div>
            </div>
            <div class="lsm-container">
                <div class="lsm-scroll" >
                    <div class="lsm-sidebar">
                        <ul>
                            <li class="lsm-sidebar-item">
                                {{ range $k,$v := .Sidebar }}
                                <a href="{{ $k }}"><i class="iconfont lsm-sidebar-icon icon_1"></i><span>{{ $v }}</span><i class="iconfont lsm-sidebar-more"></i></a>
                                {{ end }}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

    </div>
    <div class="body-right">
        <h2 class="text-center">{{ .Title }}</h2>
     
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
                        <pre class="dl-expl">{{ .Request }} </pre>
                    </div>
                {{ end }}

                {{ if .Callbak }}
                    <h3>返回参数说明</h3>
                    <div class="dl-table dl-table1">
                        <span>参数名</span>
                        <span>类型</span>
                        <span>必定返回</span>
                        <span>说明</span>
                    </div>
                    {{ range .Callbak }}
                        <div class="dl-table dl-table1 dl-table-msg  dl-table-msg1">
                            <span>{{ .Name }}</span>
                            <span>{{ .Typ }}</span>
                            <span>{{ .Need }}</span>
                            <span>{{ .Information }}</span>
                        </div>
                    {{ end }}
                {{ end }}

                {{ if .Response }}
                    <h3>返回示例</h3>
                    <div class="dl-ex-box">
                        <pre class="dl-expl">{{ .Response }}
                        </pre>
                    </div>
                {{ end }}
                {{ $field := .CodeField }}
                {{ if .CodeMsg }}
                    <h3>错误码</h3>
                    <div class="dl-table dl-table1">
                        <span>字段名</span>
                        <span>错误码</span>
                        <span>说明</span>
                    </div>
                    {{ range $k, $v := .CodeMsg }}
                        <div class="dl-table dl-table1 dl-table-msg  dl-table-msg1">
                            <span>{{ $field }}</span>
                            <span>{{ $k }}</span>
                            <span>{{ $v }}</span>
                        </div>
                    {{ end }}
                {{ end }}

                {{ if .Supplement }}
                    <h3>备注</h3>
                    <div class="dl-bz">{{ .Supplement }}</div>
                {{ end }}
            </div>
        </div>
        {{ end }}
    </div>
</div>
</body>

<script type="text/javascript" src="/-/js/jquery.js"></script>
<script type="text/javascript" src="/-/js/slimscroll.js"></script>
<script type="text/javascript" src="/-/js/left.js"></script>
<script type="text/javascript" src="/-/js/click.js"></script>

</html>
`

// <script type="text/javascript" src="/-/js/left.js"></script>
// <script type="text/javascript" src="/-/js/slimscroll.js"></script>
// <script type="text/javascript" src="/-/js/click.js"></script>
// <script type="text/javascript" src="http://download.hyahm.com/js/jquery.min.js"></script>
// <script type="text/javascript" src="http://download.hyahm.com/js/jquery.slimscroll.min.js"></script>
// <script type="text/javascript" src="http://download.hyahm.com/js/left-side-menu.js"></script>
// <script type="text/javascript" src="http://download.hyahm.com/js/click.js"></script>
