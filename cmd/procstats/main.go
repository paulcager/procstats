package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paulcager/procstats"
)

// Default Request Handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello %s!</h1>", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", defaultHandler)
	procstats.Start()
	go busy(4 * time.Second)
	go busy(7 * time.Second)
	for {
		time.Sleep(900 * time.Millisecond)
		fmt.Println("M", procstats.LastMinute()[58:])
		fmt.Println("H", procstats.LastHour()[59:])
	}
	http.ListenAndServe(":8080", nil)
}

func busy(delay time.Duration) {
	var x uint
	time.Sleep(delay)
	for {
		x = (x * 11) + 7
	}
}
