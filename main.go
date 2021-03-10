package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		&pr0Collector{
			statusCode: prometheus.NewDesc(
				"pr0gramm_http_status_code",
				"pr0gramm's status code",
				[]string{"path"}, nil,
			),
			APIDecodable: prometheus.NewDesc(
				"pr0gramm_api_json_decodable",
				"Returns 1 if JSON is decodable and 0 if not",
				[]string{"path"}, nil,
			),
		},
	)

	_ = http.ListenAndServe("0.0.0.0:4242", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
}

type pr0Collector struct {
	statusCode   *prometheus.Desc
	APIDecodable *prometheus.Desc
}

func (c *pr0Collector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.statusCode
	descs <- c.APIDecodable
}

func (c *pr0Collector) Collect(metrics chan<- prometheus.Metric) {
	resp, err := http.Get("https://pr0gramm.com")
	if err != nil {
		panic(err)
	}

	metrics <- prometheus.MustNewConstMetric(
		c.statusCode,
		prometheus.GaugeValue,
		float64(resp.StatusCode),
		"/",
	)

	resp, err = http.Get("https://pr0gramm.com/api/items")
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
		metrics <- prometheus.MustNewConstMetric(
			c.APIDecodable,
			prometheus.GaugeValue,
			0.0,
			"/api/items",
		)
	} else {
		metrics <- prometheus.MustNewConstMetric(
			c.APIDecodable,
			prometheus.GaugeValue,
			1.0,
			"/api/items",
		)
	}
}
