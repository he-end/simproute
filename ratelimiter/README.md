# Rate Limit as Middleware

Developer can use this Rate Limiter function as Middleware.

## 1. how to use

It's very easy, you just call `AddRateLimit()`to create a Middleware Function, which is then applied to`*Route{}`

Or, here we provide a `function` to create a Global RateLimit. It is slightly different from using `AddRateLimit()`,

like `GlobalRateLimit(NewRateLimit())`

## 2. example

### 2.1 use global middleware

```go
package main

import (
	"time"

	"github.com/he-end/simproute/ratelimiter"
)

func main() {
	rtr := NewRouter()

	// 1. create moiddlewar function for Rate Limit
	globalRateLimti := ratelimiter.GlobalRateLimit(ratelimiter.NewRateLimiter(10, 5, 5*time.Second))
	// 2. register to Router
	rtr.Use(globalRateLimti)
}

```

### 2.2 use with explicit path

```go
package main

import (
	"time"

	"github.com/he-end/simproute/ratelimiter"
)

func main() {
	rtr := NewRouter()

	// 1. create moiddlewar function for Rate Limit
	userRateLimit := ratelimiter.AddRateLimit("/api/users", 10, 5, 1*time.Second)
	// 2. register to Router
	rtr.Use(userRateLimit)
}
```
