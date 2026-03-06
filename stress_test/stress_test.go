package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	// "time"

	"github.com/he-end/simproute/routes"
	"github.com/he-end/simproute/routes/response"
)

func main() {
	// serv()
	// scenarios()
	// Wait for server to start
	// time.Sleep(500 * time.Millisecond)
}
func serv() {
	routes := routes.New()

	responser := response.NewWithGlobalLogger()
	responser.Dev = true
	for i := 1; i <= 10000; i++ {
		routes.GET(fmt.Sprintf("/users/{id}/:name/%v", i), func(w http.ResponseWriter, r *http.Request) {
			// time.Sleep(1 * time.Microsecond)
			responser.Success(w, "ok", nil)
		})
	}

	server := &http.Server{Addr: ":8080", Handler: routes}
	log.Println("server run in :8080")
	server.ListenAndServe()
}

func scenarios() {
	url := "http://localhost:8080/users/12/myname/10"

	job := make(chan string, 1000)

	result := make(chan int, 100000)
	// Gunakan tipe data int64 untuk atomic
	var resErr int64 = 0
	var resSucc int64 = 0
	// var resBusy int64 = 0

	worker := func(jobs <-chan string, result chan<- int) {
		client := http.Client{}
		for j := range jobs {
			req, err := http.NewRequest("GET", j, nil)
			if err != nil {
				resErr++
				result <- 1
				req.Context().Done()
				continue
			}

			res, err := client.Do(req)

			if res.StatusCode == 200 {
				atomic.AddInt64(&resSucc, 1)
			} else {
				atomic.AddInt64(&resErr, 1)
			}
			res.Body.Close()
			result <- 1
		}
	}

	for w := 1; w <= 10000; w++ {
		go worker(job, result)
	}

	for j := 1; j <= 100000; j++ {
		job <- url
	}

	// Saat print, gunakan atomic.LoadInt64
	for range result {
		fmt.Printf("success: %d\n", atomic.LoadInt64(&resSucc))
		fmt.Printf("error: %d\n", atomic.LoadInt64(&resErr))
	}
}
