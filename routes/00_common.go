package routes

import (
	"net/http"
	"strings"
	"sync"

	logger "github.com/he-end/simproute/route_logger"
	"github.com/he-end/simproute/routes/response"
)

type HandlerFunc http.HandlerFunc



type Router struct {
	MU sync.RWMutex

	Routes map[string]map[string]http.Handler
	tree   *node

	// prefix for gouping
	Prefix string

	Mws []func(http.Handler) http.Handler

	AutoCorelation bool

	RecoverOnPanic bool
}

// # return of
//
//	Autocorelation = default(true)
//	RecoverOnPanic = default(true)
func New() *Router {
	logger.InitLogger("dev", "debug")
	if resp := response.NewWithGlobalLogger(); resp != nil {
		resp.Dev = true
	}

	return &Router{
		Routes: make(map[string]map[string]http.Handler),
		tree: &node{
			handlers: make(map[string]http.Handler),
		},

		Mws:            make([]func(http.Handler) http.Handler, 0, 4),
		AutoCorelation: true,
		RecoverOnPanic: true,
	}

}

func (r *Router) Handle(method []string, path string, handler HandlerFunc) {
	// fixing path if abnormal
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	// apply prefix if set (use from grouping)
	if r.Prefix != "" {
		if path == "/" {
			path = r.Prefix
		} else {
			path = r.Prefix + path
		}
	}

	r.MU.Lock()
	defer r.MU.Unlock()

	var normalizedMethods []string
	for _, m := range method {
		methodStr := strings.ToUpper(strings.TrimSpace(m))
		if methodStr != "" {
			normalizedMethods = append(normalizedMethods, methodStr)
		}
	}

	if isDynamicRoute(path) {
		r.tree.insert(path, normalizedMethods, http.HandlerFunc(handler))
	} else {
		if r.Routes[path] == nil {
			r.Routes[path] = make(map[string]http.Handler)
		}
		for _, m := range normalizedMethods {
			r.Routes[path][m] = http.HandlerFunc(handler)
		}
	}
}

func isDynamicRoute(path string) bool {
	return strings.Contains(path, ":") || strings.Contains(path, "{") || strings.Contains(path, "*")
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.Handle([]string{"GET"}, path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.Handle([]string{"POST"}, path, handler)
}
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.Handle([]string{"PATCH"}, path, handler)
}
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.Handle([]string{"PUT"}, path, handler)
}
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.Handle([]string{"DELETE"}, path, handler)
}
