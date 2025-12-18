# 1. Example

## 1.1 Simple Use

```go
  
package main

import (
	"fmt"
	"net/http"

	logger "github.com/he-end/simproute/route_logger"
	"github.com/he-end/simproute/routes"
	"go.uber.org/zap"
)

func main() {
        // register new Routes -> return *Router
	r := routes.New()

        // you can set more options
        //
        // AutoCorelation bool
        // RecoverOnPanic bool
        // as default is 'true', but you can change as 'false'
        r.AutoCorelation = false
        r.RecoverOnPanic = false

        // [*] using single route / or simple Route without Prefix Grouping
  
        // ==>> GET /api
        r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"oke":true}`))
	}) 
        // ==>> GET /api/user1/1
        r.GET("/api/{name_user}/:id", PostHandler)  

        // [*] using grouping params/path
  
  
        r.Group("/api/v2/users", gr func(*Router){
                // ==>> GET /api/v2/users
                gr.Get("/", GetUserHandler)
                // ==>> PUT /api/v2/users/products/12
                gr.PUT("/products/:id", PostPaymentHandler)
        })
  
        // we also create logger in the package, as default logger will run as 
        // logger.InitLogger("dev", "debug")
        // in the engine, but you can change the option with "prod"/"production",
        // and level as "debug", "info", "warn", "error", "fatal", "panic" for minimal to getting Logs.
        //
        // Note the order, InitLogger() is executed simultaneously with Router() registration.
	logger.InitLogger("dev", "debug")

        // add the Http Server With Handler or final result *Router{} where you register using New()
	newServ := &http.Server{Addr: ":8080", Handler: r}
	fmt.Println("==============================")
	fmt.Println(`server run in : localhost:8080`)
	fmt.Println("==============================")
	if err := newServ.ListenAndServe(); err != nil || err != http.ErrServerClosed {
		logger.Error("undefine error", zap.Error(err))
	}
}

```

# 2. Logger

## 2.1 The Printer logger "DEV" Mode

structure of logger is  `TEXT`+`JSON` format.

### 2.1.1 Logger No Resource Request

```json
2025-12-19T04:07:30.143+0700	INFO	routes/serv_http.go:62	http_request	{"method": "GET", "path": "/as", "status": 404, "ip": "[::1]:38758", "duration": "76.879µs", "request_id": "627a2619-bf20-48f2-a8ad-bba0a4fa86b8"}
```

### 2.1.2 Logger Request Success

```json
2025-12-19T04:08:23.117+0700	INFO	routes/serv_http.go:62	http_request	{"method": "GET", "path": "/api", "status": 200, "ip": "[::1]:58864", "duration": "122.278µs", "request_id": "2c900f34-cbbb-49a4-8f29-153f2fb4a0a5"}
```

### 2.1.3 Logger Error on Panic Recovered

```json
2025-12-19T04:13:54.156+0700	ERROR	routes/serv_http.go:74	panic recovered	{"error": "runtime error: invalid memory address or nil pointer dereference", "method": "GET", "path": "/api", "request_id": "2b00362d-3950-4424-a232-75ef7d3b9c29", "request_id": "2b00362d-3950-4424-a232-75ef7d3b9c29"}
github.com/he-end/simproute/routes.(*Router).ServeHTTP.func2
	/home/hend/development/github/simproute/routes/serv_http.go:74
runtime.gopanic
	/usr/lib/go-1.22/src/runtime/panic.go:770
runtime.panicmem
	/usr/lib/go-1.22/src/runtime/panic.go:261
runtime.sigpanic
	/usr/lib/go-1.22/src/runtime/signal_unix.go:881
main.main.func1
	/home/hend/development/github/simproute/main.go:16
net/http.HandlerFunc.ServeHTTP
	/usr/lib/go-1.22/src/net/http/server.go:2171
github.com/he-end/simproute/routes.(*Router).ServeHTTP.mwAutoCorelation.func3.1
	/home/hend/development/github/simproute/routes/mw.go:22
net/http.HandlerFunc.ServeHTTP
	/usr/lib/go-1.22/src/net/http/server.go:2171
github.com/he-end/simproute/routes.(*Router).ServeHTTP
	/home/hend/development/github/simproute/routes/serv_http.go:166
net/http.serverHandler.ServeHTTP
	/usr/lib/go-1.22/src/net/http/server.go:3142
net/http.(*conn).serve
	/usr/lib/go-1.22/src/net/http/server.go:2044
2025-12-19T04:13:54.156+0700	INFO	routes/serv_http.go:62	http_request	{"method": "GET", "path": "/api", "status": 500, "ip": "[::1]:50070", "duration": "275.525µs", "request_id": "2b00362d-3950-4424-a232-75ef7d3b9c29"}
```


## 2.2 The Printer Logger on 'Prod' Mode

when use `prod` envoronment, log will be saved on `logs/` directory. but you still can see the log as `json` literal.

### 2.2.1 Logger Not Found

```json
{
  "level":"info","timestamp":"2025-12-19T04:20:56.054+0700",
  "caller":"routes/serv_http.go:62",
  "msg":"http_request",
  "method":"GET",
  "path":"/apssi",
  "status":404,
  "ip":"[::1]:40510",
  "duration":0.000068929,
  "request_id":"8fbd60a8-287a-4c39-badd-73895dea3bf7"
}

```

### 2.2.2 Logger Error

```json
{
  "level":"error",
  "timestamp":"2025-12-19T04:22:40.267+0700",
  "caller":"routes/serv_http.go:74",
  "msg":"panic recovered",
  "error":"runtime error: invalid memory address or nil pointer dereference","method":"GET","path":"/api","request_id":"c5576ff2-bfa5-4c18-b2a0-14fa55cfdea8","request_id":"c5576ff2-bfa5-4c18-b2a0-14fa55cfdea8","stacktrace":"github.com/he-end/simproute/routes.(*Router).ServeHTTP.func2\n\t/home/hend/development/github/simproute/routes/serv_http.go:74\nruntime.gopanic\n\t/usr/lib/go-1.22/src/runtime/panic.go:770\nruntime.panicmem\n\t/usr/lib/go-1.22/src/runtime/panic.go:261\nruntime.sigpanic\n\t/usr/lib/go-1.22/src/runtime/signal_unix.go:881\nmain.main.func1\n\t/home/hend/development/github/simproute/main.go:16\nnet/http.HandlerFunc.ServeHTTP\n\t/usr/lib/go-1.22/src/net/http/server.go:2171\ngithub.com/he-end/simproute/routes.(*Router).ServeHTTP.mwAutoCorelation.func3.1\n\t/home/hend/development/github/simproute/routes/mw.go:22\nnet/http.HandlerFunc.ServeHTTP\n\t/usr/lib/go-1.22/src/net/http/server.go:2171\ngithub.com/he-end/simproute/routes.(*Router).ServeHTTP.mwAutoCorelation.func3.1\n\t/home/hend/development/github/simproute/routes/mw.go:22\nnet/http.HandlerFunc.ServeHTTP\n\t/usr/lib/go-1.22/src/net/http/server.go:2171\ngithub.com/he-end/simproute/routes.(*Router).ServeHTTP\n\t/home/hend/development/github/simproute/routes/serv_http.go:166\nnet/http.serverHandler.ServeHTTP\n\t/usr/lib/go-1.22/src/net/http/server.go:3142\nnet/http.(*conn).serve\n\t/usr/lib/go-1.22/src/net/http/server.go:2044"
}
```

# 3 OPTIONS Method

in the package method OPTIONS is automatic include.

any `path` can `called with Method OPTIONS`, the response code is `201 No Content`.
