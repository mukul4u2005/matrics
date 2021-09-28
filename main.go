package main

/*
Main method for scrapper is handler() function, This function calling GetPodCountsByhttp() metrix.go
by using http call

Other method podCountByNamespace is calling GetPodCount() which is using rest api
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/mukul4u2005/matrix-service/metrix"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PodCount struct {
	Namespace string `json:"Namespace"`
	PodCount  int    `json:"Pod Count"`
}

var (
	requestsCounter uint64 = 0
	podCount               = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cluster_pod_count",
		Help: "Total pods in cluster",
	})
	hitCounter = prometheus.NewCounterFunc(prometheus.CounterOpts{
		Name: "hit_counter",
		Help: "Endpoint hit count",
	}, func() float64 {
		return float64(atomic.LoadUint64(&requestsCounter))
	})
	registry = prometheus.NewRegistry()
)

func init() {
	// Metrics have to be registered to be exposed:
	registry.Register(podCount)
	registry.Register(hitCounter)
}
func podCountByNamespace(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Finding pod count")
	var namespace string
	namespaceArr := r.URL.Query()["namespace"]
	if len(namespaceArr) == 0 {
		namespace = ""
	} else {
		namespace = namespaceArr[0]
	}
	podCount := PodCount{
		Namespace: namespace,
		PodCount:  metrix.GetPodCount(namespace),
	}
	json.NewEncoder(w).Encode(podCount)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&requestsCounter, 1)
	fmt.Fprintf(w, "Go home page")

}

func handler(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&requestsCounter, 1)
	var y float64 = float64(metrix.GetPodCountsByhttp())
	podCount.Set(y)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/metrics", handler)
	http.HandleFunc("/rest/metrics", podCountByNamespace)
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}

func main() {

	handleRequest()
}
