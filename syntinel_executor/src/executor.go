package main

import (
	"log"
	"net/http"
	"syntinel_executor/controller"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// Makes, e.g., /test/1 and /test/1/ properly resolve to the same handler.
	router.StrictSlash(true)
	router.HandleFunc("/test/{id:\\d+}", controller.RegisterTest).Methods("POST")
	router.HandleFunc("/test/{id:\\d+}", controller.DeleteTest).Methods("DELETE")
	router.HandleFunc("/test/{id:\\d+}", controller.RegisterTest).Methods("PATCH")
	router.HandleFunc("/docker/{id:\\d+}", controller.RegisterDocker).Methods("POST")
	router.HandleFunc("/docker/{id:\\d+}", controller.DeleteDocker).Methods("DELETE")
	router.HandleFunc("/docker/{id:\\d+}", controller.UpdateDocker).Methods("PATCH")
	router.HandleFunc("/test/run/{id:\\d+}", controller.RunTest).Methods("POST")
	router.HandleFunc("/test/run/{id:\\d+}", controller.KillTest).Methods("DELETE")
	router.HandleFunc("/test/run/{id:\\d+}", controller.QueryTest).Methods("GET")
	http.Handle("/", router)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalln(err)
	}
	log.Println("Exiting")
}
