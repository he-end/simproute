package routes

import (
	"net/http"
	"regexp"
	"strings"
	"sync"

	logger "github.com/he-end/simproute/route_logger"
	"github.com/he-end/simproute/routes/response"
)

type HandlerFunc http.HandlerFunc

type routePattern struct {
	pattern    string
	regex      *regexp.Regexp
	paramNames []string
}

type Router struct {
	MU sync.RWMutex

	Routes map[string]map[string]http.Handler

	DynamicRoutes []struct {
		pattern routePattern
		method  map[string]http.Handler
	}

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
		DynamicRoutes: make([]struct {
			pattern routePattern
			method  map[string]http.Handler
		}, 0, 4),

		Mws:            make([]func(http.Handler) http.Handler, 0, 4),
		AutoCorelation: true,
		RecoverOnPanic: true,
	}

}

// compilePattern converts a route pattern to a regex and extracts parameter names
// Supports both :param and {param} syntax
// Example: "/users/:id/posts/{postId}" -> regex with paramNames ["id", "postId"]
func compilePattern(pattern string) routePattern {
	regexPattern := pattern
	paramNames := make([]string, 0)

	syntaxParam := `([^/]+)`
	// colon param handle :param
	colonRegex := regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9]*)`)
	colonMatch := colonRegex.FindAllStringSubmatch(pattern, -1)
	for _, match := range colonMatch {
		param := match[1]
		// if param != "" {
		// 	paramNames = append(paramNames, param)
		// }
		paramNames = append(paramNames, param)
		regexPattern = strings.ReplaceAll(regexPattern, match[0], syntaxParam)
	}

	// brace param handle {param}
	braceRegex := regexp.MustCompile(`\{([a-zA-Z][a-zA-Z0-9]*)\}`)
	barceMatch := braceRegex.FindAllStringSubmatch(pattern, -1)
	for _, match := range barceMatch {
		param := match[1]
		// if param != "" {
		// 		paramNames = append(paramNames, param)
		// }
		paramNames = append(paramNames, param)
		regexPattern = strings.ReplaceAll(regexPattern, match[0], syntaxParam)
	}

	regexPattern = regexp.QuoteMeta(regexPattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\(\[\^/\]\+\)`, syntaxParam)
	regexPattern = "^" + regexPattern + "$"

	compileRegex := regexp.MustCompile(regexPattern)

	return routePattern{
		pattern:    pattern,
		regex:      compileRegex,
		paramNames: paramNames,
	}
}

func isDynamixRoute(pattern string) bool {
	return strings.Contains(pattern, ":") || strings.Contains(pattern, "{")
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
	// set MU
	r.MU.Lock()
	defer r.MU.Unlock()

	if isDynamixRoute(path) {
		// check if exists pattern
		pattern := compilePattern(path)
		// pattern check
		for i, dnr := range r.DynamicRoutes {
			if dnr.pattern.pattern == pattern.pattern {
				for _, m := range method {
					method := strings.ToUpper(strings.TrimSpace(m))
					if method == "" {
						continue
					}
					r.DynamicRoutes[i].method[method] = http.HandlerFunc(handler)
				}
				return
			}
		}

		// if code here, thats mean pattern not yet to register
		// add pattern
		methodMaps := make(map[string]http.Handler)
		for _, m := range method {
			method := strings.ToUpper(strings.TrimSpace(m))
			if method == "" {
				continue
			}
			methodMaps[method] = http.HandlerFunc(handler)
		}

		r.DynamicRoutes = append(r.DynamicRoutes, struct {
			pattern routePattern
			method  map[string]http.Handler
		}{
			pattern: pattern,
			method:  methodMaps,
		})

	} else {
		if r.Routes[path] == nil {
			r.Routes[path] = map[string]http.Handler{}
		}
		for _, m := range method {
			method := strings.ToUpper(strings.TrimSpace(m))
			if method == "" {
				continue
			}
			r.Routes[path][method] = http.HandlerFunc(handler)
		}

	}

}

func (r *Router) Get(path string, handler HandlerFunc) {
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
