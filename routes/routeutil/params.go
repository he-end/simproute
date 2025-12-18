package routeutil

import (
	"context"
)

// Context key for route parameters
type routeParamsKey struct{}

// RouteParams represents extracted route parameters
type RouteParams map[string]string

// GetRouteParams extracts route parameters from request context
// Usage: params := routeutil.GetRouteParams(r.Context())
func GetRouteParams(ctx context.Context) RouteParams {
	if params, ok := ctx.Value(routeParamsKey{}).(RouteParams); ok {
		return params
	}
	return make(RouteParams)
}

// Get extracts a specific parameter value from context
// Usage: userID := routeutil.GetRouteParams(r.Context()).Get("id")
func (rp RouteParams) Get(key string) string {
	return rp[key]
}

// SetRouteParams sets route parameters in context (used internally by router)
func SetRouteParams(ctx context.Context, params RouteParams) context.Context {
	return context.WithValue(ctx, routeParamsKey{}, params)
}
