package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func main() {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	for i := 1; i < 5; i++ {
		time.Sleep(400 * time.Millisecond)
		//r := 0
		//for r < 1 {
		//	r = rand.Intn(50)
		//	}
		fmt.Println(i)
		body := strings.NewReader(`{"TestID":` + fmt.Sprint(i) + `, "Interval": ` + fmt.Sprint(i) + `}`)
		req, err := http.NewRequest("POST", "http://localhost:9093/schedule", body)
		if err != nil {
			// handle err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
	}
}
