package routes

import (
	"net/http"
	"strings"
)

func (r *Router) handleOptions(w http.ResponseWriter, req *http.Request) bool {
	path := req.URL.Path
	var methods []string

	// Cek static route (O(1))
	if m, ok := r.Routes[path]; ok {
		for method := range m {
			methods = append(methods, method)
		}
	}

	// Cek dynamic route (O(N) traversal)
	if len(methods) == 0 {
		handlers, _, found := r.tree.search(path)
		if found && len(handlers) > 0 {
			for method := range handlers {
				methods = append(methods, method)
			}
		}
	}

	if len(methods) == 0 {
		return false
	}
	methods = append(methods, http.MethodOptions)

	// Set header CORS
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.WriteHeader(http.StatusNoContent)
	return true
}
