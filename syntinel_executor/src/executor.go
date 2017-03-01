package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syntinel_executor/DAO/database"
	"syntinel_executor/PAO"
	"syntinel_executor/controller"
	"syntinel_executor/service/statistics"

	"github.com/gorilla/mux"
)

func Port() string {
	if len(os.Args) > 1 {
		return fmt.Sprintf(":%v", os.Args[1])
	}
	return ":8080"
}

func main() {
	defer log.Println("Exiting")
	port := Port()
	log.Println("Initializing Database")
	database.InitDB()
	PAO.StartPAO()
	statistics.StartStatisticsServer()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := mux.NewRouter()
	// Makes, e.g., /test/1 and /test/1/ properly resolve to the same handler.
	router.StrictSlash(true)

	router.HandleFunc("/stats", controller.Statistics).Methods("GET")

	router.HandleFunc("/test/run", controller.RunTest).Methods("POST")
	router.HandleFunc("/test/run", controller.KillTest).Methods("DELETE")
	router.HandleFunc("/test/run", controller.QueryTest).Methods("GET")

	http.Handle("/", router)
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 5)
	// 		log.Println(fmt.Sprintf("Number of Goroutines: %v", runtime.NumGoroutine()))
	// 	}
	// }()
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalln(err)
	}
}
