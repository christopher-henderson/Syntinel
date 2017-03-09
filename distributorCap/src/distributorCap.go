package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"./LoadBalancer"
	"./controller"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting")
	router := mux.NewRouter()
	//requests to register must be in format {"hostName":"localhost", "port": "9093", "Scheme": "http"}
	router.HandleFunc("/register", controller.AddExecutor).Methods("POST")
	log.Println("somthing")
	router.HandleFunc("/schedule", controller.ScheduleTest).Methods("POST")
	router.HandleFunc("/cancel", controller.CancelTest).Methods("POST")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy := LoadBalancer.GetReverseProxy(0, false)
		proxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":9093", nil); err != nil {
		log.Println("huh")
		log.Fatalln(err)
	}

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL.Path)
			buf, _ := ioutil.ReadAll(r.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr2
			if r.URL.Path == "/test/run" {
				body, err := ioutil.ReadAll(rdr1)
				if err != nil {
					log.Println("error")
				}
				var t Scheduler.ExecutorRequestObj
				err = json.Unmarshal(body, &t)
				if err != nil {
					log.Println(t.ID)
					proxy := LoadBalancer.GetReverseProxy(0, false)
					proxy.ServeHTTP(w, r)
				} else {
					proxy := LoadBalancer.GetReverseProxy(t.ID, true)
					proxy.ServeHTTP(w, r)
				}
			} else {
				proxy := LoadBalancer.GetReverseProxy(0, false)
				proxy.ServeHTTP(w, r)
			}

		})

		http.HandleFunc("/cancelrunning", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("error")
			}
			log.Println(string(body))
			var t Scheduler.Job
			err = json.Unmarshal(body, &t)
			if err != nil {
				log.Println("error")
				io.WriteString(w, "There was a problem with your submission")
			} else {

			}
		})


		log.Fatal(http.ListenAndServe(":9093", nil))
	*/

}
