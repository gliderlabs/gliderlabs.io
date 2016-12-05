package web

import (
	"net/http"
	"os"
	"text/template"

	"github.com/facebookgo/httpdown"

	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("web", &Component{},
		com.Option("listen_addr", "0.0.0.0:"+os.Getenv("PORT"), "Address and port to listen on"),
		com.Option("static_dir", "ui/static/", "Directory to serve static files from"),
		com.Option("static_path", "/static", "URL path to serve static files at"),
		com.Option("cookie_secret", "", "Random string to use for session cookies"))
}

// Handler extension point for matching and handling HTTP requests
type Handler interface {
	MatchHTTP(r *http.Request) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Handlers accessor for web.Handler extension point
func Handlers() []Handler {
	var handlers []Handler
	for _, com := range com.Enabled(new(Handler), nil) {
		handlers = append(handlers, com.(Handler))
	}
	return handlers
}

type TemplateFuncProvider interface {
	WebTemplateFuncMap(r *http.Request) template.FuncMap
}

func TemplateFuncMap(r *http.Request) template.FuncMap {
	funcMap := template.FuncMap{}
	for _, com := range com.Enabled(new(TemplateFuncProvider), nil) {
		for k, v := range com.(TemplateFuncProvider).WebTemplateFuncMap(r) {
			funcMap[k] = v
		}
	}
	return funcMap
}

// Web component
type Component struct {
	http httpdown.Server
}
