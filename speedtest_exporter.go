// Copyright (C) 2016, 2017 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/dchest/uniuri"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	prom_version "github.com/prometheus/common/version"

	"github.com/nlamirault/speedtest_exporter/speedtest"
	"github.com/nlamirault/speedtest_exporter/version"
)

const (
	namespace = "speedtest"
)

var (
	ping = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ping"),
		"Latency (ms)",
		nil, nil,
	)
	download = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download"),
		"Download bandwidth (Mbps).",
		nil, nil,
	)
	upload = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload"),
		"Upload bandwidth (Mbps).",
		nil, nil,
	)
)

// Exporter collects Speedtest stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	Client *speedtest.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter(config string, server string, interval time.Duration) (*Exporter, error) {
	log.Infof("Setup Speedtest client with interval %s", interval)
	client, err := speedtest.NewClient(config, server)
	if err != nil {
		return nil, fmt.Errorf("Can't create the Speedtest client: %s", err)
	}

	log.Debugln("Init exporter")
	return &Exporter{
		Client: client,
	}, nil
}

// Describe describes all the metrics ever exported by the Speedtest exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- ping
	ch <- download
	ch <- upload
}

// Collect fetches the stats from configured Speedtest location and delivers them
// as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Infof("Speedtest exporter starting")
	if e.Client == nil {
		log.Errorf("Speedtest client not configured.")
		return
	}

	metrics := e.Client.NetworkMetrics()
	ch <- prometheus.MustNewConstMetric(ping, prometheus.GaugeValue, metrics["ping"])
	ch <- prometheus.MustNewConstMetric(download, prometheus.GaugeValue, metrics["download"])
	ch <- prometheus.MustNewConstMetric(upload, prometheus.GaugeValue, metrics["upload"])
	log.Infof("Speedtest exporter finished")
}

func init() {
	prometheus.MustRegister(prom_version.NewCollector("speedtest_exporter"))
}

func main() {
	var (
		showVersion   = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("web.listen-address", ":9112", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		configURL     = flag.String("speedtest.config-url", "http://c.speedtest.net/speedtest-config.php?x="+uniuri.New(), "Speedtest configuration URL")
		serverURL     = flag.String("speedtest.server-url", "http://c.speedtest.net/speedtest-servers-static.php?x="+uniuri.New(), "Speedtest server URL")
		//interval      = flag.Int("interval", 60*time.Second, "Interval for metrics.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Speedtest Prometheus exporter. v%s\n", version.Version)
		os.Exit(0)
	}

	log.Infoln("Starting speedtest exporter", prom_version.Info())
	log.Infoln("Build context", prom_version.BuildContext())

	interval := 60 * time.Second
	exporter, err := NewExporter(*configURL, *serverURL, interval)
	if err != nil {
		log.Errorf("Can't create exporter : %s", err)
		os.Exit(1)
	}
	log.Infoln("Register exporter")
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Speedtest Exporter</title></head>
             <body>
             <h1>Speedtest Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
