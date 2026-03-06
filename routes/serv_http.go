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

	if r.AutoCorelation {
		idCorelate := goruntime.GetCorelationID()
		w.Header().Set("X-Set-Corelation-ID", idCorelate.String())
		logger.NewLoggerOnRuntime(logger.RegisterRuntime{Key: "request_id", Value: idCorelate.String()})
		
		defer goruntime.ClearCorelationID()
		defer logger.DeferDeleteRuntimeValue()
	}
	defer func() {
		// panic recovered
		if r.RecoverOnPanic {
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

	// Remove old snapshot middleware comments since they are not used anymore.
	var handler http.Handler
	var routeParams routeutil.RouteParams

	// O(1) Static Route Match First
	if methodForPath, exist := r.Routes[path]; exist {
		if h, ok := methodForPath[method]; ok {
			handler = h
		} else {
			response.NewWithGlobalLogger().Fail(rec, "Method Not Allowed", "METHOD_NOT_ALLOWED", "The method is not allowed for the requested URL")
			return
		}
	} else {
		handlers, params, found := r.tree.search(path)
		if found {
			if h, ok := handlers[method]; ok {
				handler = h
				routeParams = params
			} else {
				response.NewWithGlobalLogger().Fail(rec, "Method Not Allowed", "METHOD_NOT_ALLOWED", "The method is not allowed for the requested URL")
				return
			}
		} else {
			response.NewWithGlobalLogger().Error(rec, "Not Found", "NOT_FOUND", "The requested resource was not found", http.StatusNotFound)
			return
		}
	}

	currentMws := r.Mws

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

func (r *Router) RoutesTreeSearchForTest(path string) (map[string]http.Handler, routeutil.RouteParams, bool) {
	return r.tree.search(path)
}
