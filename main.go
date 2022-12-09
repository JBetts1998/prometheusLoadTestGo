package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricNumber = 20
const clusterNumber = 1

var Metric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "metric",
		Help: "A fake metric for testing Prometheus.",
	},
	[]string{"cluster", "node", "metric"},
)

func recordMetrics(numNodes []int) {
	go func() {
		for {
			rand.Seed(time.Now().UnixNano())
			for i := 0; i < clusterNumber; i++ {
				clusterName := fmt.Sprintf("cluster-%d", i)
				nodeNumber := numNodes[i]
				for j := 0; j < nodeNumber; j++ {
					nodeName := fmt.Sprintf("node-%d", j)
					for k := 0; k < 20; k++ {
						if k == 0 {
							metricName := "healthy"
							Metric.WithLabelValues(clusterName, nodeName, metricName).Set(math.Round(rand.Float64()))
						} else {
							metricName := fmt.Sprintf("metric-%d", k)
							Metric.WithLabelValues(clusterName, nodeName, metricName).Set(rand.Float64() * 100)
						}

					}
				}
			}
			time.Sleep(time.Second * 2)
		}
	}()
}

func main() {
	var numNodes []int
	for i := 0; i < clusterNumber; i++ {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(10)
		numNodes = append(numNodes, num)
	}
	recordMetrics(numNodes)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
