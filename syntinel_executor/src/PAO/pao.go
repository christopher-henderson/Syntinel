package PAO

import "log"

var paoRun chan int = make(chan int)
var paoKill chan int = make(chan int)

var testQueueMap map[int]*TestQueue = make(map[int]*TestQueue)

func StartPAO() {
	go dispatch()
}

func Run(id int) {
	log.Println("Received run request in PAO.")
	select {
	case paoRun <- id:
	default:
	}
}

func run(id int) {
	log.Println("Running ", id)
	test, ok := testQueueMap[id]
	if !ok {
		// @TODO make the information passing jive.
		test = NewTestQueue(id)
		testQueueMap[id] = test
	}
	go test.Run()
}

func Kill(id int) {
	log.Println("Received kill request in PAO.")
	select {
	case paoKill <- id:
	default:
	}
}

func kill(id int) {
	if test, ok := testQueueMap[id]; ok {
		log.Println("Attempt to kill ", id)
		go test.Kill()
	} else {
		log.Println("No killable process found.")
	}
}

func dispatch() {
	log.Println("Starting PAO request loop.")
	for {
		select {
		case id := <-paoRun:
			log.Println("Received run request in dispatcher.")
			run(id)
		case id := <-paoKill:
			log.Println("Received kill request in dispatcher.")
			kill(id)
		}
	}
}
