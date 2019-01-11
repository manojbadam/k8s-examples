package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		delay, err := strconv.Atoi(os.Getenv("DELAY_SECONDS"))
		if err != nil {
			delay = 0
		}
		time.Sleep(time.Duration(delay) * time.Second)
		fmt.Printf("Request recieved %s at %v\n", r.URL.Path, time.Now())
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	gracePeriodDealySeconds, err := strconv.Atoi(os.Getenv("GRACE_PERIOD_DELAY_SECONDS"))
	if err == nil && gracePeriodDealySeconds > 0 {
		fmt.Printf("Sleeping for %s seconds \n", os.Getenv("GRACE_PERIOD_DELAY_SECONDS"))
		time.Sleep(time.Duration(gracePeriodDealySeconds) * time.Second)
	}
	fmt.Printf("Server is ready to handle requests at http://localhost:8080/ \n")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		fmt.Printf("Could not listen on 8080: %v\n", err)
	}
}
