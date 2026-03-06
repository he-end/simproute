package routeutil

import (
	"context"
)

// Context key for route parameters
type routeParamsKey struct{}

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// RouteParams represents extracted route parameters
type RouteParams []Param

// GetRouteParams extracts route parameters from request context
// Usage: params := routeutil.GetRouteParams(r.Context())
func GetRouteParams(ctx context.Context) RouteParams {
	if params, ok := ctx.Value(routeParamsKey{}).(RouteParams); ok {
		return params
	}
	return nil
}

// Get extracts a specific parameter value from context
// Usage: userID := routeutil.GetRouteParams(r.Context()).Get("id")
func (rp RouteParams) Get(key string) string {
	for i := range rp {
		if rp[i].Key == key {
			return rp[i].Value
		}
	}
	return ""
}

// SetRouteParams sets route parameters in context (used internally by router)
func SetRouteParams(ctx context.Context, params RouteParams) context.Context {
	return context.WithValue(ctx, routeParamsKey{}, params)
}
