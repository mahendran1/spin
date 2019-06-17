package main

import (
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
  	"github.com/prometheus/client_golang/prometheus/promhttp"
  	"log"
  	"time"
  	"math/rand
)

var (
  counter = prometheus.NewCounter(
     prometheus.CounterOpts{
        Namespace: "golang",
        Name:      "my_counter",
        Help:      "This is my counter",
     })

  gauge = prometheus.NewGauge(
     prometheus.GaugeOpts{
        Namespace: "golang",
        Name:      "my_gauge",
        Help:      "This is my gauge",
     })

  histogram = prometheus.NewHistogram(
     prometheus.HistogramOpts{
        Namespace: "golang",
        Name:      "my_histogram",
        Help:      "This is my histogram",
     })

  summary = prometheus.NewSummary(
     prometheus.SummaryOpts{
        Namespace: "golang",
        Name:      "my_summary",
        Help:      "This is my summary",
     })
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %+v\n", r);
	bs, err := ioutil.ReadFile("/content/index.html")

	if err != nil {
		fmt.Printf("Couldn't read index.html: %v", err);
		os.Exit(1)
	}

	io.WriteString(w, string(bs[:]))
}

func main() {
	http.HandleFunc("/", index)
	port := ":8000"
	fmt.Printf("Starting to service on port %s\n", port);
	http.ListenAndServe(port, nil)
	rand.Seed(time.Now().Unix())

  http.Handle("/metrics", promhttp.Handler())

  prometheus.MustRegister(counter)
  prometheus.MustRegister(gauge)
  prometheus.MustRegister(histogram)
  prometheus.MustRegister(summary)

  go func() {
     for {
        counter.Add(rand.Float64() * 5)
        gauge.Add(rand.Float64()*15 - 5)
        histogram.Observe(rand.Float64() * 10)
        summary.Observe(rand.Float64() * 10)

        time.Sleep(time.Second)
     }
  }()

  log.Fatal(http.ListenAndServe(":8080", nil))
}
}
