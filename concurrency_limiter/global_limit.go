package concurrencylimiter

import (
	"net/http"
)

func (cl *ConcurrenctLimit) MwCCLimit() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case cl.limitReq <- struct{}{}:
				defer func() { <-cl.limitReq }()
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
		})
	}
}
