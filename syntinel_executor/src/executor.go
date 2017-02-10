package main

import (
	"log"
	"net/http"
	"syntinel_executor/DAO/database"
	"syntinel_executor/PAO"
	"syntinel_executor/controller"

	"github.com/gorilla/mux"
)

func main() {
	defer log.Println("Exiting")
	log.Println("Initializing Database")
	database.InitDB()
	PAO.StartPAO()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := mux.NewRouter()
	// Makes, e.g., /test/1 and /test/1/ properly resolve to the same handler.
	router.StrictSlash(true)

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
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalln(err)
	}
}
