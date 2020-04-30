package circuitbreaker

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func Example_circuitBreaker() {
	getGoogle := func(hc *http.Client, circuitBreaker Interface) {
		circuitBreaker.Do(
			func() error {
				resp, err := hc.Get("https://google.com")
				if err != nil {
					return err
				}
				resp.Write(os.Stdout)
				return nil
			},
			func(err error) {
				fmt.Println("google is not available. Error:", err)
			},
			func() {
				fmt.Println("google is not available. Circuit breaker is opened")
			},
		)
	}

	circuitBreaker := NewSyncCircuitBreaker(1, time.Second)
	client := &http.Client{Timeout: time.Nanosecond}

	// this call get will timeout
	getGoogle(client, circuitBreaker)
	// this call will be denied by circuit breaker because of it's opened
	getGoogle(client, circuitBreaker)

}
