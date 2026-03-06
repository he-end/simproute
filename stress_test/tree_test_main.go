package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/he-end/simproute/routes"
// )

// func main() {
// 	r := routes.New()
// 	handler := func(w http.ResponseWriter, req *http.Request) {
// 		w.WriteHeader(200)
// 	}

// 	r.GET("/users/{id}/name/{name}/aksndaskdnask", handler)

// 	path := "/users/100/name/you/aksndaskdnask"
// 	req, _ := http.NewRequest("GET", path, nil)
// 	handlers, params, found := r.RoutesTreeSearchForTest(req.URL.Path)
// 	fmt.Printf("Tree Search Direct - found: %v, params: %v, methods: %d\n", found, params, len(handlers))
	
// 	// Test via HTTP Server
// 	go http.ListenAndServe(":8081", r)
// 	time.Sleep(200*time.Millisecond)

// 	res, err := http.Get("http://localhost:8081" + path)
// 	if err != nil {
// 		fmt.Printf("HTTP Error: %v\n", err)
// 	} else {
// 		fmt.Printf("HTTP Status: %d\n", res.StatusCode)
// 	}
// }
