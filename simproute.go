package simproute

import "github.com/he-end/simproute/routes"

type Routes = routes.Router

func NewRouter() *Routes {
	return routes.New()
}
