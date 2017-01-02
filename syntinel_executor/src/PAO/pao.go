package PAO

import "log"

var paoRun chan int = make(chan int)
var paoKill chan int = make(chan int)

var cancelMap map[int]*TestQueue = make(map[int]*TestQueue)

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
	test, ok := cancelMap[id]
	if !ok {
		// @TODO make the information passing jive.
		test = NewTestQueue(200, id)
		cancelMap[id] = test
	}
	test.Run()
}

func Kill(id int) {
	log.Println("Received kill request in PAO.")
	paoKill <- id
}

func kill(id int) {
	if test, ok := cancelMap[id]; ok {
		log.Println("Found killable process.")
		test.Kill()
	}
	log.Println("No killable process found.")
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
