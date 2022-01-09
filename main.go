package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var requests_count = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests",
	},
	[]string{"path"},
)


func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		log.Printf(path)
		next.ServeHTTP(w, r)
		requests_count.WithLabelValues(path).Inc()
	})
}


func welcomeHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/index.html")
}


func goodbyeHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/goodbye.html")
}


func init(){
	prometheus.Register(requests_count)
}


func main(){
	// Serve static files
	staticRouter := mux.NewRouter()
	staticRouter.Use(middleware)

	staticRouter.HandleFunc("/welcome", welcomeHandler)
	staticRouter.HandleFunc("/goodbye", goodbyeHandler)	

	go func(){
		log.Printf("Serving static files on port :8080")
		err1 := http.ListenAndServe(":8080", staticRouter)
		log.Fatal(err1)
	}()

	metricsRouter := mux.NewRouter()
	metricsRouter.Path("/metrics").Handler(promhttp.Handler())

	log.Printf("Exposing metrics on port :8081")
	err := http.ListenAndServe(":8081", metricsRouter)
	log.Fatal(err)
}
