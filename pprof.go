package xmux

import (
	"net/http/pprof"
)

func Pprof() *RouteGroup {
	pp := NewRouteGroup().BindResponse(nil)

	pp.Get("/debug/pprof/{all:name}", pprof.Index)
	pp.Get("/debug/pprof/cmdline", pprof.Cmdline)
	pp.Get("/debug/pprof/profile", pprof.Profile)
	pp.Get("/debug/pprof/symbol", pprof.Symbol)
	pp.Get("/debug/pprof/trace", pprof.Trace)
	return pp
}

// Cmdline responds with the running program's
// command line, with arguments separated by NUL bytes.
// The package initialization registers it as /debug/pprof/cmdline.
// func cmdline(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	fmt.Fprint(w, strings.Join(os.Args, "\x00"))
// }

// func sleep(r *http.Request, d time.Duration) {
// 	select {
// 	case <-time.After(d):
// 	case <-r.Context().Done():
// 	}
// }

// func durationExceedsWriteTimeout(r *http.Request, seconds float64) bool {
// 	srv, ok := r.Context().Value(http.ServerContextKey).(*http.Server)
// 	return ok && srv.WriteTimeout != 0 && seconds >= srv.WriteTimeout.Seconds()
// }

// func serveError(w http.ResponseWriter, status int, txt string) {
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	w.Header().Set("X-Go-Pprof", "1")
// 	w.Header().Del("Content-Disposition")
// 	w.WriteHeader(status)
// 	fmt.Fprintln(w, txt)
// }

// // Profile responds with the pprof-formatted cpu profile.
// // Profiling lasts for duration specified in seconds GET parameter, or for 30 seconds if not specified.
// // The package initialization registers it as /debug/pprof/profile.
// func profile(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	sec, err := strconv.ParseInt(r.FormValue("seconds"), 10, 64)
// 	if sec <= 0 || err != nil {
// 		sec = 30
// 	}

// 	if durationExceedsWriteTimeout(r, float64(sec)) {
// 		serveError(w, http.StatusBadRequest, "profile duration exceeds server's WriteTimeout")
// 		return
// 	}

// 	// Set Content Type assuming StartCPUProfile will work,
// 	// because if it does it starts writing.
// 	w.Header().Set("Content-Type", "application/octet-stream")
// 	w.Header().Set("Content-Disposition", `attachment; filename="profile"`)
// 	if err := pprof.StartCPUProfile(w); err != nil {
// 		// StartCPUProfile failed, so no writes yet.
// 		serveError(w, http.StatusInternalServerError,
// 			fmt.Sprintf("Could not enable CPU profiling: %s", err))
// 		return
// 	}
// 	sleep(r, time.Duration(sec)*time.Second)
// 	pprof.StopCPUProfile()
// }

// // Trace responds with the execution trace in binary form.
// // Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified.
// // The package initialization registers it as /debug/pprof/trace.
// func tra(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	sec, err := strconv.ParseFloat(r.FormValue("seconds"), 64)
// 	if sec <= 0 || err != nil {
// 		sec = 1
// 	}

// 	if durationExceedsWriteTimeout(r, sec) {
// 		serveError(w, http.StatusBadRequest, "profile duration exceeds server's WriteTimeout")
// 		return
// 	}

// 	// Set Content Type assuming trace.Start will work,
// 	// because if it does it starts writing.
// 	w.Header().Set("Content-Type", "application/octet-stream")
// 	w.Header().Set("Content-Disposition", `attachment; filename="trace"`)
// 	if err := trace.Start(w); err != nil {
// 		// trace.Start failed, so no writes yet.
// 		serveError(w, http.StatusInternalServerError,
// 			fmt.Sprintf("Could not enable tracing: %s", err))
// 		return
// 	}
// 	sleep(r, time.Duration(sec*float64(time.Second)))
// 	trace.Stop()
// }

// // Symbol looks up the program counters listed in the request,
// // responding with a table mapping program counters to function names.
// // The package initialization registers it as /debug/pprof/symbol.
// func symbol(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

// 	// We have to read the whole POST body before
// 	// writing any output. Buffer the output here.
// 	var buf bytes.Buffer

// 	// We don't know how many symbols we have, but we
// 	// do have symbol information. Pprof only cares whether
// 	// this number is 0 (no symbols available) or > 0.
// 	fmt.Fprintf(&buf, "num_symbols: 1\n")

// 	var b *bufio.Reader
// 	if r.Method == "POST" {
// 		b = bufio.NewReader(r.Body)
// 	} else {
// 		b = bufio.NewReader(strings.NewReader(r.URL.RawQuery))
// 	}

// 	for {
// 		word, err := b.ReadSlice('+')
// 		if err == nil {
// 			word = word[0 : len(word)-1] // trim +
// 		}
// 		pc, _ := strconv.ParseUint(string(word), 0, 64)
// 		if pc != 0 {
// 			f := runtime.FuncForPC(uintptr(pc))
// 			if f != nil {
// 				fmt.Fprintf(&buf, "%#x %s\n", pc, f.Name())
// 			}
// 		}

// 		// Wait until here to check for err; the last
// 		// symbol will have an err because it doesn't end in +.
// 		if err != nil {
// 			if err != io.EOF {
// 				fmt.Fprintf(&buf, "reading request: %v\n", err)
// 			}
// 			break
// 		}
// 	}

// 	w.Write(buf.Bytes())
// }

// // Handler returns an HTTP handler that serves the named profile.
// // func Handler(name string) http.Handler {
// // 	return handler(name)
// // }

// // type handler string

// func debug(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	name := Var(r)["name"]
// 	p := pprof.Lookup(name)
// 	if p == nil {
// 		serveError(w, http.StatusNotFound, "Unknown profile")
// 		return
// 	}
// 	gc, _ := strconv.Atoi(r.FormValue("gc"))
// 	if name == "heap" && gc > 0 {
// 		runtime.GC()
// 	}
// 	debug, _ := strconv.Atoi(r.FormValue("debug"))
// 	if debug != 0 {
// 		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	} else {
// 		w.Header().Set("Content-Type", "application/octet-stream")
// 		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
// 	}
// 	p.WriteTo(w, debug)
// }

// var profileDescriptions = map[string]string{
// 	"allocs":       "A sampling of all past memory allocations",
// 	"block":        "Stack traces that led to blocking on synchronization primitives",
// 	"cmdline":      "The command line invocation of the current program",
// 	"goroutine":    "Stack traces of all current goroutines",
// 	"heap":         "A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.",
// 	"mutex":        "Stack traces of holders of contended mutexes",
// 	"profile":      "CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.",
// 	"threadcreate": "Stack traces that led to the creation of new OS threads",
// 	"trace":        "A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.",
// }

// // Index responds with the pprof-formatted profile named by the request.
// // For example, "/debug/pprof/heap" serves the "heap" profile.
// // Index responds to a request for "/debug/pprof/" with an HTML page
// // listing the available profiles.
// func index(w http.ResponseWriter, r *http.Request) {
// 	// if strings.HasPrefix(r.URL.Path, "/debug/pprof/") {
// 	// 	name := Var(r)["name"]
// 	// 	if name != "" {
// 	// 		handler(name).ServeHTTP(w, r)
// 	// 		return
// 	// 	}
// 	// }

// 	type profile struct {
// 		Name  string
// 		Href  string
// 		Desc  string
// 		Count int
// 	}
// 	var profiles []profile
// 	for _, p := range pprof.Profiles() {
// 		profiles = append(profiles, profile{
// 			Name:  p.Name(),
// 			Href:  p.Name() + "?debug=1",
// 			Desc:  profileDescriptions[p.Name()],
// 			Count: p.Count(),
// 		})
// 	}

// 	// Adding other profiles exposed from within this package
// 	for _, p := range []string{"cmdline", "profile", "trace"} {
// 		profiles = append(profiles, profile{
// 			Name: p,
// 			Href: p,
// 			Desc: profileDescriptions[p],
// 		})
// 	}

// 	sort.Slice(profiles, func(i, j int) bool {
// 		return profiles[i].Name < profiles[j].Name
// 	})

// 	if err := indexTmpl.Execute(w, profiles); err != nil {
// 		log.Print(err)
// 	}
// }

// var indexTmpl = template.Must(template.New("index").Parse(`<html>
// <head>
// <title>/debug/pprof/</title>
// <style>
// .profile-name{
// 	display:inline-block;
// 	width:6rem;
// }
// </style>
// </head>
// <body>
// /debug/pprof/<br>
// <br>
// Types of profiles available:
// <table>
// <thead><td>Count</td><td>Profile</td></thead>
// {{range .}}
// 	<tr>
// 	<td>{{.Count}}</td><td><a href={{.Href}}>{{.Name}}</a></td>
// 	</tr>
// {{end}}
// </table>
// <a href="goroutine?debug=2">full goroutine stack dump</a>
// <br/>
// <p>
// Profile Descriptions:
// <ul>
// {{range .}}
// <li><div class=profile-name>{{.Name}}:</div> {{.Desc}}</li>
// {{end}}
// </ul>
// </p>
// </body>
// </html>
// `))
