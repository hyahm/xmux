package xmux

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type MethodStrcut struct {
	Summary    string              `json:"summary,omitempty" yaml:"summary"`
	Parameters []Parameter         `json:"parameters,omitempty" yaml:"parameters"`
	Responses  map[string]Response `json:"responses,omitempty" yaml:"responses"`
	Produces   []string            `json:"produces,omitempty" yaml:"produces"`
	Consumes   []string            `json:"consumes,omitempty" yaml:"consumes"`
}

type Parameter struct {
	In       string            `json:"in,omitempty" yaml:"in"`
	Name     string            `json:"name,omitempty" yaml:"name"`
	Required bool              `required:"in,omitempty" yaml:"required"`
	Type     string            `json:"type,omitempty" yaml:"type"`
	Schema   map[string]string `json:"schema,omitempty" yaml:"schema"`
	Minimum  int64             `json:"minimum,omitempty" yaml:"minimum"`
}

type Schema struct {
	Type       string                     `json:"type" yaml:"type"`
	Properties map[string]map[string]Type `json:"properties" yaml:"properties"`
}

type Type struct {
	Type    string `json:"type" yaml:"type"`
	Example string `json:"example" yaml:"example"`
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

func (r *Router) ShowSwagger(url, jsonPath, host string, schemes ...string) *GroupRoute {
	if jsonPath == "" {
		jsonPath = "/swagger.json"
	}
	swagger := NewGroupRoute().BindResponse(nil)
	swagger.Get(url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		opts := DefaultEnsure(jsonPath)
		tmpl := template.Must(template.New("swaggerui").Parse(swaggeruiTemplate))
		buf := bytes.NewBuffer(nil)
		_ = tmpl.Execute(buf, &opts)
		w.Write(buf.Bytes())
	})
	swagger.Get(jsonPath, func(w http.ResponseWriter, req *http.Request) {
		// 拿到路由
		ss := schemes
		if len(schemes) == 0 {
			ss = []string{"http"}
		}
		swagger := Swagger{
			Swagger: "2.0",
			Host:    host,
			Schemes: ss,
			Paths:   make(map[string]map[string]MethodStrcut),
		}
		for k, v := range r.route {
			if k == url || k == jsonPath {
				continue
			}
			path := make(map[string]MethodStrcut)
			for method := range v.methods {
				path[strings.ToLower(method)] = MethodStrcut{
					Summary: v.describe,
				}
			}
			swagger.Paths[k] = path
		}
		send, err := json.MarshalIndent(swagger, "", "\t")
		if err != nil {
			log.Println(err)
		}
		w.Write(send)
	})
	return swagger
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
