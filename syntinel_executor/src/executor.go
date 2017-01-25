package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"syntinel_executor/DAO/database"
	"syntinel_executor/PAO"
	"syntinel_executor/controller"
	"time"

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
	router.HandleFunc("/test/{id:\\d+}", controller.RegisterTest).Methods("POST")
	router.HandleFunc("/test/{id:\\d+}", controller.DeleteTest).Methods("DELETE")
	router.HandleFunc("/test/{id:\\d+}", controller.RegisterTest).Methods("PATCH")

	router.HandleFunc("/test/run", controller.RunTest).Methods("POST")
	router.HandleFunc("/test/run", controller.KillTest).Methods("DELETE")
	router.HandleFunc("/test/run", controller.QueryTest).Methods("GET")

	router.HandleFunc("/script/{id:\\d+}", controller.RegisterScript).Methods("POST")
	router.HandleFunc("/script/{id:\\d+}", controller.DeleteScript).Methods("DELETE")
	router.HandleFunc("/script/{id:\\d+}", controller.UpdateScript).Methods("PATCH")

	router.HandleFunc("/docker/{id:\\d+}", controller.RegisterDocker).Methods("POST")
	router.HandleFunc("/docker/{id:\\d+}", controller.DeleteDocker).Methods("DELETE")
	router.HandleFunc("/docker/{id:\\d+}", controller.UpdateDocker).Methods("PATCH")

	http.Handle("/", router)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			log.Println(fmt.Sprintf("Number of Goroutines: %v", runtime.NumGoroutine()))
		}
	}()
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalln(err)
	}
}
