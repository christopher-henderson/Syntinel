package PAO

import "log"

var paoRun chan [2]int = make(chan [2]int)
var paoKill chan [2]int = make(chan [2]int)

var testQueueMap map[int]*TestRunQueue = make(map[int]*TestRunQueue)

func StartPAO() {
	go dispatch()
}

func Run(testID, testRunID int) {
	log.Println("Received run request in PAO.")
	select {
	case paoRun <- [2]int{testID, testRunID}:
	default:
	}
}

func run(testID, testRunID int) {
	test, ok := testQueueMap[testID]
	if !ok {
		// @TODO make the information passing jive.
		test = NewTestRunQueue(testID)
		testQueueMap[testID] = test
	}
	go test.Run(testRunID)
}

func Kill(testID, testRunID int) {
	log.Println("Received kill request in PAO.")
	select {
	case paoKill <- [2]int{testID, testRunID}:
	default:
	}
}

func kill(testID, testRunID int) {
	if test, ok := testQueueMap[testID]; ok {
		log.Println("Attempt to kill ", testID)
		go test.Kill(testRunID)
	} else {
		log.Println("No killable process found.")
	}
}

func Query(testID, testRunID int) int {
	if testQueue, ok := testQueueMap[testID]; ok {
		return testQueue.Query(testRunID)
	}
	return NotFound
}

func dispatch() {
	log.Println("Starting PAO request loop.")
	for {
		select {
		case IDs := <-paoRun:
			log.Println("Received run request in dispatcher.")
			run(IDs[0], IDs[1])
		case IDs := <-paoKill:
			log.Println("Received kill request in dispatcher.")
			kill(IDs[0], IDs[1])
		}
	}
}
