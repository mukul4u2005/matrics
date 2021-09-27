package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mukul4u2005/matrix-service/matrix"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PodCount struct {
	Namespace string `json:"Namespace"`
	PodCount  int    `json:"Pod Count"`
}

var (
	podCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cluster_pod_count",
		Help: "Total pods in cluster",
	})
	hitCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hit_counter",
		Help: "Endpoint hit count",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(podCount)
	prometheus.MustRegister(hitCounter)
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
		PodCount:  matrix.GetPodCount(namespace),
	}
	json.NewEncoder(w).Encode(podCount)
	//hitCounter.Inc()

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go home page")
	//hitCounter.Inc()

}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/podcount", podCountByNamespace)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}

func main() {
	var y float64 = float64(matrix.GetPodCount(""))
	podCount.Set(y)
	hitCounter.Inc()
	handleRequest()
}
