package routes

import (
	"net/http"
	"time"

	"github.com/he-end/simproute/goruntime"
	logger "github.com/he-end/simproute/route_logger"
	"github.com/he-end/simproute/routes/response"
	"github.com/he-end/simproute/routes/routeutil"
	"go.uber.org/zap"
)

type responseRecorer struct {
	http.ResponseWriter
	status int
	size   int
}

func (rr *responseRecorer) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorer) Write(b []byte) (int, error) {
	if rr.status == 0 {
		rr.status = http.StatusOK
	}
	n, err := rr.ResponseWriter.Write(b)
	rr.size = n
	return n, err
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	rec := &responseRecorer{ResponseWriter: w, status: http.StatusOK}

	// if r.AutoCorelation {
	// 	r.Use(mwAutoCorelation())
	// 	defer goruntime.ClearCorelationID()
	// 	defer logger.DeferDeleteRuntimeValue()
	// }
	defer func() {
		// panic recovered
		if r.RecoverOnPanic {
			defer func() {
				if recvr := recover(); recvr != nil {
					fields := []zap.Field{zap.Any("error", recvr), zap.String("method", req.Method), zap.String("path", req.URL.Path)}
					if r.AutoCorelation {
						rID := goruntime.GetCorelationID()
						fields = append(fields, zap.String("request_id", rID.String()))
					}
					logger.Error("panic recovered",
						fields...,
					)
					// Use response handler to send a safe error response
					response.NewWithGlobalLogger().Error(rec, "Internal server error", response.ErrCodeInternalError, "An unexpected error occurred", http.StatusInternalServerError)
				}
			}()
			if rec.status == 0 {
				rec.status = http.StatusOK
			}
			dur := time.Since(start)
			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", rec.status),
				zap.String("ip", req.RemoteAddr),
				zap.Duration("duration", dur),
			}
			if r.AutoCorelation {
				rID := goruntime.GetCorelationID()
				fields = append(fields, zap.String("request_id", rID.String()))
			}

			// its call external package 'logger' for create auto log
			logger.GetLogger().Info("http_request", fields...)

		}
	}()

	// Lookup
	method := req.Method
	path := req.URL.Path

	// Handle preflight OPTIONS
	if method == http.MethodOptions {
		if r.handleOptions(w, req) {
			return
		}
	}

	// r.MU.RLock()
	// // snapshot middleware slice for this request
	// mws := append([]func(http.Handler) http.Handler(nil), r.Mws...)
	// // snapshot routes for this request
	// routes := make(map[string]map[string]http.Handler)
	// for p, methods := range r.Routes {
	// 	routes[p] = methods
	// }
	// // snapshot dynamic routes for this request
	// dynamicRoutes := make([]struct {
	// 	pattern routePattern
	// 	method  map[string]http.Handler
	// }, len(r.DynamicRoutes))

	// copy(dynamicRoutes, r.DynamicRoutes)
	r.MU.RLock()

	var handler http.Handler
	var routeParams routeutil.RouteParams

	// First, try exact match (static routes)
	if methodForPath, exist := r.Routes[req.URL.Path]; exist {
		handler = methodForPath[req.Method]
	}
	// else {
	// 	// 405 Method Not Allowed
	// 	response.NewWithGlobalLogger().Fail(rec, "Method Not Allowed", "METHOD_NOT_ALLOWED", "The method is not allowed for the requested URL")
	// 	return
	// }

	if handler == nil {
		// Try dynamic route matching
		found := false
		for _, dr := range r.DynamicRoutes {
			if matches := dr.pattern.regex.FindStringSubmatch(path); matches != nil {
				// Extract parameters
				if methodHandler, exists := dr.method[method]; exists {
					routeParams = make(routeutil.RouteParams)
					for i, paramName := range dr.pattern.paramNames {
						if i+1 < len(matches) {
							routeParams[paramName] = matches[i+1]
						}
					}

					// Check if method is allowed for this pattern
					handler = methodHandler
					found = true
					break
				} else {
					r.MU.RUnlock()
					// Pattern matches but method not allowed - 405
					response.NewWithGlobalLogger().Fail(rec, "Method Not Allowed", "METHOD_NOT_ALLOWED", "The method is not allowed for the requested URL")
					return
				}
			}
		}

		if !found {
			// 404 Not Found
			r.MU.RUnlock()
			response.NewWithGlobalLogger().Error(rec, "Not Found", "NOT_FOUND", "The requested resource was not found", http.StatusNotFound)
			return
		}
	}
	currentMws := r.Mws
	r.MU.RUnlock()

	// Inject route parameters into request context
	if routeParams != nil {
		ctx := routeutil.SetRouteParams(req.Context(), routeParams)
		req = req.WithContext(ctx)
	}

	// Wrap handler with middleware chain (outer-most last registered)
	for i := len(currentMws) - 1; i >= 0; i-- {
		handler = currentMws[i](handler)
	}

	handler.ServeHTTP(rec, req)
}
