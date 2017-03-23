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

type TestSuccess struct {
	Successful bool `json:"successful"`
}

type TestFailure struct {
	Err        string `json:"error"`
	Successful bool   `json:"successful"`
}

func Finalize(ID int, testErr error) error {
	url := fmt.Sprintf("http://localhost/api/v1/testrun/%v", ID)
	var serializedResult []byte
	var err error
	if testErr == nil {
		if serializedResult, err = json.Marshal(TestSuccess{true}); err != nil {
			log.Println(err)
			return err
		}
	} else {
		if serializedResult, err = json.Marshal(TestFailure{testErr.Error(), false}); err != nil {
			log.Println(err)
			return err
		}
	}
	headers := map[string][]string{
		"Content-Type": []string{"application/json"},
	}
	_, _, r, err := http.DefaultClient.Patch(url, headers, bytes.NewBuffer(serializedResult))
	if err != nil {
		log.Println(err)
		return err
	}
	if r != nil {
		defer r.Close()
	}
	return nil
}
