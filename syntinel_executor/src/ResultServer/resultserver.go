package ResultServer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/http"
	"golang.org/x/net/websocket"
)

const maxBackoff = 20

// origin := "http://localhost:8000/"
// url := fmt.Sprintf("ws://localhost:8000/testRun/console/%v-executor", ID)
// ws, err := websocket.Dial(url, "", origin)
// if err != nil {
//     log.Fatal(err)
// }
// if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
//     log.Fatal(err)
// }
// var msg = make([]byte, 512)
// var n int
// if n, err = ws.Read(msg); err != nil {
//     log.Fatal(err)
// }
// fmt.Printf("Received: %s.\n", msg[:n])

func Stream(ID int, stdout *bufio.Scanner) {
	origin := "http://localhost"
	url := fmt.Sprintf("ws://localhost/testRun/console/%v-executor", ID)
	log.Println(origin)
	log.Println(url)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for stdout.Scan() {
			log.Println(fmt.Sprintf("STDOUT ID %v: %v", ID, string(stdout.Bytes())))
			if _, err := ws.Write(stdout.Bytes()); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

type TestFinalization struct {
	Err        string `json:"error"`
	Successful bool   `json:"successful"`
}

func Finalize(ID int, testErr error) error {
	log.Println("Finalizing")
	log.Println("Ping")
	url := fmt.Sprintf("http://localhost/api/v1/testrun/%v", ID)
	var errorMessage string
	if testErr != nil {
		errorMessage = testErr.Error()
	}
	log.Println("Ping")
	obj, _ := json.Marshal(TestFinalization{errorMessage, testErr == nil})
	headers := map[string][]string{
		"Content-Type": []string{"application/json"},
	}
	log.Println("Ping")
	status, _, r, err := http.DefaultClient.Patch(url, headers, bytes.NewBuffer(obj))
	log.Println("Ping")
	if err != nil {
		log.Println("Ping")
		log.Println("error!")
		log.Println(err)
		return err
	}
	log.Println("Ping")
	if r != nil {
		defer r.Close()
		log.Println("Ping")
	}
	log.Println("Ping")
	log.Printf("Patch result: %v", status)
	return nil
}
