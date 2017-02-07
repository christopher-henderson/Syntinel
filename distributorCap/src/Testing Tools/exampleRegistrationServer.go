package main

import (
	"io"
	"net/http"
    "io/ioutil"
    "fmt"
    "bytes"
    
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world from server 3!")
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":9093",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = hello
    
    url := "http://localhost:9090/register"
    fmt.Println("URL:>", url)
    var jsonStr = []byte(`{"HostName":"localhost", "Port": "9093"}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}