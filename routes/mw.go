package routes

import (
	"net/http"

	"github.com/he-end/simproute/goruntime"
	logger "github.com/he-end/simproute/route_logger"
)

func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.MU.Lock()
	r.Mws = append(r.Mws, mw)
	r.MU.Unlock()
}

func mwAutoCorelation() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idCorelate := goruntime.GetCorelationID()
			w.Header().Set("X-Set-Corelation-ID", idCorelate.String())
			logger.NewLoggerOnRuntime(logger.RegisterRuntime{Key: "request_id", Value: idCorelate.String()})
			next.ServeHTTP(w, r)
		})
	}
}
