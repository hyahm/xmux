package xmux

import (
	"bytes"
	"encoding/json"

	jsonv2 "encoding/json/v2"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type MethodStrcut struct {
	Summary    string              `json:"summary,omitempty" yaml:"summary"`
	Parameters []Parameter         `json:"parameters,omitempty" yaml:"parameters"`
	Responses  map[string]Response `json:"responses,omitempty" yaml:"responses"`
	Produces   []string            `json:"produces,omitempty" yaml:"produces" required:""`
	Consumes   []string            `json:"consumes,omitempty" yaml:"consumes"`
}

type ParameterType string

const (
	Query  ParameterType = "query"
	Path   ParameterType = "path"
	Header ParameterType = "header"
	Form   ParameterType = "formData"
)

type Parameter struct {
	In          ParameterType     `json:"in,omitempty" yaml:"in"`
	Name        string            `json:"name,omitempty" yaml:"name"`
	Required    bool              `required:"in,omitempty" yaml:"required"`
	Type        string            `json:"type,omitempty" yaml:"type"`
	Schema      map[string]string `json:"schema,omitempty" yaml:"schema"`
	Minimum     int64             `json:"minimum,omitempty" yaml:"minimum"`
	Enum        []string          `json:"enum,omitempty" yaml:"enum"`
	Default     any               `json:"default,omitempty" yaml:"default"`
	Description string            `json:"description,omitempty" yaml:"description"`
}

type Schema struct {
	Type       string          `json:"type" yaml:"type"`
	Properties map[string]Type `json:"properties" yaml:"properties"` // key是字段名
	Ref        string          `json:"$ref" yaml:"$ref"`             // $ref: '#/definitions/User'
}

type Type struct {
	Type        string `json:"type" yaml:"type"`
	Description string `json:"description" yaml:"description"`
}

type Response struct {
	Description string            `json:"type,omitempty" yaml:"type"`
	Schema      map[string]string `json:"schema,omitempty" yaml:"schema"`
}

type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
}

type Swagger struct {
	Swagger             string                             `json:"swagger" yaml:"swagger"`
	Info                Info                               `json:"info" yaml:"info"`
	Host                string                             `json:"host" yaml:"host"`
	BasePath            string                             `json:"basePath,omitempty" yaml:"basePath"`
	Schemes             []string                           `json:"schemes,omitempty" yaml:"schemes"`
	Paths               map[string]map[string]MethodStrcut `json:"paths,omitempty" yaml:"paths"`
	Definitions         map[string]Definition              `json:"definitions,omitempty" yaml:"definitions"`
	Security            []map[string][]string              `json:"security,omitempty" yaml:"security"`
	SecurityDefinitions map[string]Type                    `json:"securityDefinitions,omitempty" yaml:"securityDefinitions"`
}

type Definition struct {
	Properties map[string]Type `json:"properties" yaml:"properties"`
	Required   []string        `json:"required" yaml:"required"`
}

type SwaggerUIOpts struct {
	// BasePath for the UI path, defaults to: /
	SpecURL string

	// The three components needed to embed swagger-ui
	SwaggerURL       string
	SwaggerPresetURL string
	SwaggerStylesURL string

	Favicon32 string
	Favicon16 string

	// Title for the documentation site, default to: API documentation
	Title string
}

func (r *Router) ShowSwagger(url, host string, schemes ...string) *RouteGroup {
	jsonPath := "/swagger.json"
	swagger := NewRouteGroup().BindResponse(nil).SetHeader("Access-Control-Allow-Origin", "*")
	swagger.SetHeader("Content-Type", "sec-ch-ua;sec-ch-ua-mobile;sec-ch-ua-platform")
	swagger.SetHeader("Access-Control-Allow-Methods", "*")
	swagger.Get(url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		opts := DefaultEnsure(jsonPath)
		tmpl := template.Must(template.New("swaggerui").Parse(swaggeruiTemplate))
		buf := bytes.NewBuffer(nil)
		_ = tmpl.Execute(buf, &opts)
		w.Write(buf.Bytes())
	})

	swagger.Get(jsonPath, JsonFile(jsonPath, url, host, r, schemes...))
	return swagger
}

func JsonFile(jsonPath, url, host string, router *Router, schemes ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 拿到路由
		ss := schemes
		if len(schemes) == 0 {
			ss = []string{"http"}
		}
		swagger := Swagger{
			Swagger: "2.0",
			Host:    host,
			Info: Info{
				Title:       router.SwaggerTitle,
				Version:     router.SwaggerVersion,
				Description: router.SwaggerDescription,
			},
			Schemes: ss,
			Paths:   make(map[string]map[string]MethodStrcut),
		}
		// 合并匹配请求
		for k, mr := range router.urlRoute {
			if k == url || k == jsonPath {
				continue
			}

			ms := MethodStrcut{
				Summary:  mr.summary,
				Produces: []string{"application/json"},
				Responses: map[string]Response{"200": {
					Description: "",
				}},
			}
			path := make(map[string]MethodStrcut)
			for _, method := range mr.methods {
				path[strings.ToLower(method)] = ms
			}
			if acc, ok := mr.header["Content-Type"]; ok {
				ms.Produces = strings.Split(acc, ";")
			} else {
				ms.Produces = []string{"application/json"}
			}
			ms.Parameters = mr.query
			swagger.Paths[k] = path
		}
		// 合并正则请求
		for url, mr := range router.urlTpl {
			path := make(map[string]MethodStrcut)
			for _, method := range mr.methods {
				ms := MethodStrcut{
					Summary:  mr.summary,
					Produces: []string{"application/json"},
					Responses: map[string]Response{"200": {
						Description: "",
					}},
				}

				if acc, ok := mr.header["Content-Type"]; ok {
					ms.Produces = strings.Split(acc, ";")
				} else {
					ms.Produces = []string{"application/json"}
				}
				ms.Parameters = mr.query

				// 正则请求还需要合并path, 自定义正则必须每个组都以^开头 $结尾， 不然无法自动生成
				for _, name := range mr.params {
					// url 进行填充
					// 将url的^ 替换成 { 将url的$ 替换成  }

					url = url[1 : len(url)-1]
					start := strings.Index(url, "(")
					end := strings.Index(url, ")")
					t := "string"
					if url[start+1:end] == interger {
						t = "interger"
					}
					// 将里面的值替换
					url = url[:start] + "{" + name + "}" + url[end+1:]
					p := Parameter{
						In:       Path,
						Name:     name,
						Required: true,
						Type:     t,
					}
					ms.Parameters = append(ms.Parameters, p)
				}

				path[strings.ToLower(method)] = ms
				swagger.Paths[url] = path
			}

		}
		var send []byte
		var err error
		if enableJsonV2 {
			send, err = jsonv2.Marshal(swagger, jsonv2.DefaultOptionsV2())
		} else {
			send, err = json.MarshalIndent(swagger, "", "  ")
		}

		if err != nil {
			log.Println(err)
		}
		w.Write(send)
	}
}

func DefaultEnsure(jsonPath string) *SwaggerUIOpts {
	return &SwaggerUIOpts{
		SpecURL:          jsonPath,
		SwaggerURL:       swaggerLatest,
		SwaggerPresetURL: swaggerPresetLatest,
		SwaggerStylesURL: swaggerStylesLatest,
		Favicon16:        swaggerFavicon16Latest,
		Favicon32:        swaggerFavicon32Latest,
		Title:            "API documentation",
	}
}

// SwaggerUI creates a middleware to serve a documentation site for a swagger spec.
// This allows for altering the spec before starting the http listener.
// func SwaggerUI(opts SwaggerUIOpts, next http.Handler) http.Handler {
// 	opts.EnsureDefaults()

// 	pth := path.Join(opts.BasePath, opts.Path)
// 	tmpl := template.Must(template.New("swaggerui").Parse(swaggeruiTemplate))

// 	buf := bytes.NewBuffer(nil)
// 	_ = tmpl.Execute(buf, &opts)
// 	b := buf.Bytes()

// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == pth {
// 			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
// 			rw.WriteHeader(http.StatusOK)

// 			_, _ = rw.Write(b)
// 			return
// 		}

// 		if next == nil {
// 			rw.Header().Set("Content-Type", "text/plain")
// 			rw.WriteHeader(http.StatusNotFound)
// 			_, _ = rw.Write([]byte(fmt.Sprintf("%q not found", pth)))
// 			return
// 		}
// 		next.ServeHTTP(rw, r)
// 	})
// }

const (
	swaggerLatest          = "https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"
	swaggerPresetLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"
	swaggerStylesLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui.css"
	swaggerFavicon32Latest = "https://unpkg.com/swagger-ui-dist/favicon-32x32.png"
	swaggerFavicon16Latest = "https://unpkg.com/swagger-ui-dist/favicon-16x16.png"
	swaggeruiTemplate      = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
		<title>{{ .Title }}</title>

    <link rel="stylesheet" type="text/css" href="{{ .SwaggerStylesURL }}" >
    <link rel="icon" type="image/png" href="{{ .Favicon32 }}" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{ .Favicon16 }}" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="{{ .SwaggerURL }}"> </script>
    <script src="{{ .SwaggerPresetURL }}"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: '{{ .SpecURL }}',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region

      window.ui = ui
    }
  </script>
  </body>
</html>
`
)
