package routes

import (
	"net/http"
	"strings"

	"github.com/he-end/simproute/routes/routeutil"
)

type node struct {
	path       string
	isParam    bool
	isWildcard bool
	paramName  string
	handlers   map[string]http.Handler
	children   []*node
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

func (n *node) insert(path string, methods []string, handler http.Handler) {
	// Fast track root insertion
	if n.path == "" && len(n.children) == 0 {
		n.insertNode(path, methods, handler)
		return
	}

walk:
	for {
		// Find longest common prefix
		i := longestCommonPrefix(path, n.path)

		if i < len(n.path) {
			child := &node{
				path:       n.path[i:],
				isParam:    n.isParam,
				isWildcard: n.isWildcard,
				paramName:  n.paramName,
				handlers:   n.handlers,
				children:   n.children,
			}

			n.children = []*node{child}
			n.path = path[:i]
			n.isParam = false
			n.isWildcard = false
			n.paramName = ""
			n.handlers = nil
		}

		if i < len(path) {
			path = path[i:]

			for _, child := range n.children {
				if child.path[0] == path[0] {
					n = child
					continue walk
				}
			}

			// Path divergence at param indicator
			if len(path) > 0 && path[0] == ':' || path[0] == '*' || path[0] == '{' {
				for _, child := range n.children {
					if child.isParam || child.isWildcard {
						n = child
						continue walk
					}
				}
			}

			child := &node{}
			child.insertNode(path, methods, handler)
			n.children = append(n.children, child)
			return
		} else if i == len(path) {
			if n.handlers == nil {
				n.handlers = make(map[string]http.Handler)
			}
			for _, m := range methods {
				n.handlers[m] = handler
			}
			return
		}
	}
}

func (n *node) insertNode(path string, methods []string, handler http.Handler) {
	// Find if this path segment contains a parameter (slow path on registration but fast on lookup)
	for i := 0; i < len(path); i++ {
		if path[i] == ':' || path[i] == '*' || path[i] == '{' {
			// Save text prefix up to the param
			if i > 0 {
				n.path = path[:i]
				
				child := &node{}
				child.insertNode(path[i:], methods, handler)
				n.children = []*node{child}
				return
			}

			// We are at a param boundary
			isWildcard := path[i] == '*'
			isRegexTyle := path[i] == '{'

			// Find end of parameter
			end := i + 1
			for end < len(path) && path[end] != '/' {
				end++
			}

			paramName := ""
            if isRegexTyle {
                if closeIdx := strings.IndexByte(path, '}'); closeIdx != -1 {
                    paramName = path[1:closeIdx]
                    end = closeIdx + 1
                }
            } else {
			    paramName = path[1:end]
            }

			n.path = path[:end]
			n.isParam = !isWildcard
			n.isWildcard = isWildcard
			n.paramName = paramName

			if end < len(path) {
				child := &node{}
				child.insertNode(path[end:], methods, handler)
				n.children = []*node{child}
				return
			}
			
			if n.handlers == nil {
				n.handlers = make(map[string]http.Handler)
			}
			for _, m := range methods {
				n.handlers[m] = handler
			}
			return
		}
	}

	// No params in path
	n.path = path
	if n.handlers == nil {
		n.handlers = make(map[string]http.Handler)
	}
	for _, m := range methods {
		n.handlers[m] = handler
	}
}

func (n *node) search(path string) (map[string]http.Handler, routeutil.RouteParams, bool) {
	var params routeutil.RouteParams

walk:
	for {
		if path == n.path {
			if len(n.handlers) > 0 {
				return n.handlers, params, true
			}
		}

		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]

				// 1. Text match
match_children:
				for _, child := range n.children {
					if !child.isParam && !child.isWildcard && len(path) > 0 && child.path[0] == path[0] {
						n = child
						continue walk
					}
				}

				// 2. Param match
				for _, child := range n.children {
					if child.isParam {
						end := 0
						for end < len(path) && path[end] != '/' {
							end++
						}

						if params == nil {
							params = make(routeutil.RouteParams, 0, 1)
						}
						params = append(params, routeutil.Param{Key: child.paramName, Value: path[:end]})
						path = path[end:]

						if len(path) == 0 {
							if len(child.handlers) > 0 {
								return child.handlers, params, true
							}
						}

						n = child
						goto match_children
					}
				}

				// 3. Wildcard match
				for _, child := range n.children {
					if child.isWildcard {
						if params == nil {
							params = make(routeutil.RouteParams, 0, 1)
						}
						params = append(params, routeutil.Param{Key: child.paramName, Value: path})
						return child.handlers, params, true
					}
				}
			}
		}

		return nil, nil, false
	}
}
