package main

import (
	"gorilla/mux"
	"net/http"
	"syntinel_executor/handlers"
)

func main() {
	router = mux.NewRouter()
  router.HandleFunc("/test/{id:d\+}", handlers.RegisterTest).Methods("POST")
  router.HandleFunc("/test/{id:d\+}", handlers.DeleteTest).Methods("DELETE")
  router.HandleFunc("/test/{id:d\+}", handlers.UpdateTest).Methods("PATCH")
  router.HandleFunc("/docker", handlers.RegisterDocker).Methods("POST")
  router.HandleFunc("/docker", handlers.DeleteDocker).Methods("DELETE")
  router.HandleFunc("/docker", handlers.UpdateDocker).Methods("PATCH")
	router.HandleFunc("/test/run/{id:d\+}", handlers.RunTest).Methods("POST")
  router.HandleFunc("/test/run/{id:d\+}", handlers.Killest).Methods("DELETE")
  router.HandleFunc("/test/run/{id:d\+}", handlers.QueryTest).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
