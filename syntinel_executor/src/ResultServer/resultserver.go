package ResultServer

import (
	"bufio"
	"fmt"
	"log"
	"time"

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
	origin := "http://192.168.1.2:8000"
	url := fmt.Sprintf("ws://192.168.1.2:8000/testRun/console/%v-executor", ID)
	log.Println(origin)
	log.Println(url)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for stdout.Scan() {
			// log.Println(fmt.Sprintf("STDOUT ID %v: %v", ID, string(stdout.Bytes())))
			if _, err := ws.Write(stdout.Bytes()); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
