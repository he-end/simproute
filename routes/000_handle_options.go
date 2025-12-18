package routes

import (
	"net/http"
	"strings"
)

func (r *Router) handleOptions(w http.ResponseWriter, req *http.Request) bool {
	path := req.URL.Path
	var methods []string

	// Cek static route
	if m, ok := r.Routes[path]; ok {
		for method := range m {
			methods = append(methods, method)
		}
	}

	// Cek dynamic route juga
	if len(methods) == 0 {
		for _, dr := range r.DynamicRoutes {
			if dr.pattern.regex.MatchString(path) {
				for method := range dr.method {
					methods = append(methods, method)
				}
				break
			}
		}
	}

	// Kalau gak ada route terdaftar, lewatin (biar 404 normal)
	if len(methods) == 0 {
		return false
	}

	// Tambahkan OPTIONS ke daftar method
	methods = append(methods, http.MethodOptions)

	// Set header CORS
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.WriteHeader(http.StatusNoContent)
	return true
}
