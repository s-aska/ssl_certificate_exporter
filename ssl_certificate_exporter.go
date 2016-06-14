package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"fmt"
)

const namespace = "ssl_certificate"

type exporter struct {
	expires *prometheus.GaugeVec
}

type config struct {
	Domains []string `json:"domains"`
}

var configUrl string
var domains []string
var m = new(sync.Mutex)

func newExporter() *exporter {
	return &exporter{
		expires: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "expires",
				Help:      "expires",
			},
			[]string{
				"domain",
			},
		),
	}
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	e.expires.Describe(ch)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	m.Lock()
	defer m.Unlock()
	for _, domain := range domains {
		s := check(domain)
		if s != math.NaN() {
			e.expires.WithLabelValues(domain).Set(s)
		}
	}

	e.expires.Collect(ch)
}

func check(domain string) float64 {
	config := tls.Config{}

	conn, err := tls.Dial("tcp", domain+":443", &config)
	if err != nil {
		log.Fatal("domain:" + domain + " error:" + err.Error())
		return math.NaN()
	}

	state := conn.ConnectionState()
	certs := state.PeerCertificates

	defer conn.Close()

	duration := certs[0].NotAfter.Unix() - time.Now().Unix()

	return float64(duration)
}

func load() {
	if configUrl == "" {
		return
	}

	resp, err := http.Get(configUrl)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode > 201 {
		log.Printf("Failure loading url:%v code:%v error:%v", configUrl, resp.StatusCode, resp.Status)
		return
	}

	var config config
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&config)

	m.Lock()
	defer m.Unlock()

	domains = config.Domains

	log.Printf("Successful loading domains:%v", strings.Join(domains, ","))
}

func main() {
	exporter := newExporter()

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/reload", reload)

	port := ":" + os.Getenv("PORT")

	// Sample
	// https://gist.githubusercontent.com/s-aska/03c41cf0d3f8b369cf0ae80d02a26c02/raw/3c742b80c4c1c7e79fb6705cda19808efb8048eb/config.json
	configUrl = os.Getenv("CONFIG_URL")
	if configUrl == "" {
		log.Fatal("Missing ENV CONFIG_URL")
	}

	load()

	if len(domains) == 0 {
		log.Fatal("Missing domains for config")
	}

	log.Print("Listening 127.0.0.1", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func reload(w http.ResponseWriter, r *http.Request) {
	load()
	fmt.Fprintf(w, "Reloading configuration file... domains:%v", strings.Join(domains, ","))
}
