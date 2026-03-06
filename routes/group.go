package routes

import (
	"net/http"
	"strings"
)

func (r *Router) Group(prefix string, fn func(gr *Router)) {
	if prefix == "" || prefix[0] != '/' {
		prefix = "/"
	}

	group := &Router{
		Routes:         make(map[string]map[string]http.Handler),
		tree:           r.tree,
		Prefix:         joinPrefix(r.Prefix, prefix),
		AutoCorelation: r.AutoCorelation,
		RecoverOnPanic: r.RecoverOnPanic,
	}
	fn(group)

	r.MU.Lock()
	defer r.MU.Unlock()

	for p, method := range group.Routes {
		if r.Routes[p] == nil {
			r.Routes[p] = make(map[string]http.Handler)
		}
		for m, h := range method {
			method := strings.ToUpper(strings.TrimSpace(m))
			if method == "" {
				continue
			}
			r.Routes[p][method] = h
		}
	}
}

func joinPrefix(a, b string) string {
	if a == "" {
		return b
	}
	if b == "/" {
		return a
	}
	return a + b
}
