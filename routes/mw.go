package routes

import (
	"net/http"
)

func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.MU.Lock()
	r.Mws = append(r.Mws, mw)
	r.MU.Unlock()
}


