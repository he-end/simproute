package concurrencylimiter

import (
	"net/http"

	"github.com/he-end/simproute/routes"
)

func (cl *ConcurrenctLimit) PerHandlerMwCCLimit(capacity int64, next http.HandlerFunc) routes.HandlerFunc {
	limitReq := make(chan struct{}, capacity)
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case limitReq <- struct{}{}:
			defer func() { <-limitReq }()
			next.ServeHTTP(w, r)
			//
		default:
			if cl.defaultResponseError == nil {
				cl.responser.Error(w, "please try again later", "SERVER_BUSY", "", http.StatusServiceUnavailable)
				return
			}
			cl.responser.Error(
				w,
				cl.defaultResponseError.Message,
				cl.defaultResponseError.ErrCode,
				cl.defaultResponseError.Details,
				cl.defaultResponseError.StatusCode,
			)
			return
		}
	}
}
