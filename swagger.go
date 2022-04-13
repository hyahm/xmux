package xmux

import (
	"bytes"
	"html/template"
	"net/http"
)

type MethodStrcut struct {
	Summary    string              `json:"summary" yaml:"summary"`
	Parameters []Parameter         `json:"parameters" yaml:"parameters"`
	Responses  map[string]Response `json:"responses" yaml:"responses"`
}

type Parameter struct {
	Un       string `json:"in" yaml:"in"`
	Name     string `json:"name" yaml:"name"`
	Required bool   `required:"in" yaml:"required"`
	Type     string `json:"type" yaml:"type"`
}

type Response struct {
	Description string            `json:"type" yaml:"type"`
	Schema      map[string]string `json:"schema" yaml:"schema"`
}

type Path map[string]MethodStrcut

// swagger: "2.0"
// host: api.example.com
// basePath: /v1
// schemes:
//   - https
// paths:
//   /users/{userId}:
//     get:
//       summary: Returns a user by ID.
//       parameters:
//         - in: path
//           name: userId
//           required: true
//           type: integer
//       responses:
//         200:
//           description: OK
//           schema:
//             $ref: '#/definitions/User'
//   /users:
//     post:
//       summary: Creates a new user.
//       parameters:
//         - in: body
//           name: user
//           schema:
//             $ref: '#/definitions/User'
//       responses:
//         200:
//           description: OK

type Swagger struct {
	Swagger  string          `json:"swagger" yaml:"swagger"`
	Host     string          `json:"host" yaml:"host"`
	BasePath string          `json:"basePath" yaml:"basePath"`
	Schemes  []string        `json:"schemes" yaml:"schemes"`
	Paths    map[string]Path `json:"paths" yaml:"paths"`
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

func (r *Router) ShowSwagger(url, host string, schemes ...string) {
	r.Get(url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		opts := DefaultEnsure()
		tmpl := template.Must(template.New("swaggerui").Parse(swaggeruiTemplate))

		buf := bytes.NewBuffer(nil)
		_ = tmpl.Execute(buf, &opts)
		GetInstance(r).Response = ""
		w.Write(buf.Bytes())
	})
	r.Get("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		// 拿到路由
		// swagger := Swagger{
		// 	Swagger: "2.0",
		// 	Host:    host,
		// 	Schemes: schemes,
		// }
		// for k, v := range r.route {
		// 	path := make(map[string]MethodStrcut)
		// 	for method, rt := range v {
		// 		path[method] = MethodStrcut{
		// 			Summary: "",
		// 		}
		// 	}
		// 	swagger.Paths[k] = path
		// }
		w.Write([]byte(``))
	})
}

func DefaultEnsure() *SwaggerUIOpts {
	return &SwaggerUIOpts{
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
        url: '/swagger.json',
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
