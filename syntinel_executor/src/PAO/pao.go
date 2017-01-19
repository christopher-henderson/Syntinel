package PAO

import "log"

var paoRun chan [2]int = make(chan [2]int)
var paoKill chan [2]int = make(chan [2]int)

var testQueueMap map[int]*TestRunner = make(map[int]*TestRunner)

func StartPAO() {
	go dispatch()
}

func Run(testID, testRunID int) {
	select {
	case paoRun <- [2]int{testID, testRunID}:
	default:
	}
}

func Kill(testID, testRunID int) {
	log.Println("Received kill request in PAO.")
	select {
	case paoKill <- [2]int{testID, testRunID}:
	default:
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

func run(testID, testRunID int) {
	testQueue, ok := testQueueMap[testID]
	if !ok {
		testQueue = NewTestRunner(testID)
		testQueueMap[testID] = testQueue
	}
	if testQueue.Query(testRunID) == NotFound {
		go testQueue.Run(testRunID)
	} else {
		log.Println("Receieved an attempt to replay the same test ID. ID = ", testRunID)
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
