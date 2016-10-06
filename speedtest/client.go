// Copyright (C) 2016 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package speedtest

import (
	"time"

	"github.com/prometheus/common/log"
	"github.com/zpeters/speedtest/sthttp"
	"github.com/zpeters/speedtest/tests"
)

// Client defines the Speedtest client
type Client struct {
	Server          sthttp.Server
	SpeedtestClient *sthttp.Client
	AllServers      []sthttp.Server
	ClosestServers  []sthttp.Server
}

// NewClient defines a new client for Speedtest
func NewClient(configURL string, serversURL string) (*Client, error) {
	log.Debugf("New Speedtest client %s %s", configURL, serversURL)
	configTimeout, _ := time.ParseDuration("15s")
	latencyTimeout, _ := time.ParseDuration("15s")
	downloadTimeout, _ := time.ParseDuration("15s")
	stClient := sthttp.NewClient(
		&sthttp.SpeedtestConfig{
			ConfigURL:       configURL,
			ServersURL:      serversURL,
			AlgoType:        "max",
			NumClosest:      3,
			NumLatencyTests: 5,
			Interface:       "",
			Blacklist:       "",
		},
		&sthttp.HTTPConfig{
			ConfigTimeout:   configTimeout,
			LatencyTimeout:  latencyTimeout,
			DownloadTimeout: downloadTimeout,
		},
		true,
		"|")

	log.Debugf("Retrieve all servers")
	var allServers []sthttp.Server
	allServers, err := stClient.GetServers()
	if err != nil {
		return nil, err
	}

	log.Debugf("Retrieve closest servers")
	closestServers := stClient.GetClosestServers(allServers)
	log.Debugf("Find test server")
	testServer := stClient.GetFastestServer(closestServers)

	return &Client{
		Server:          testServer,
		SpeedtestClient: stClient,
		AllServers:      allServers,
		ClosestServers:  closestServers,
	}, nil
}

func (client *Client) NetworkMetrics() map[string]float64 {
	result := map[string]float64{}
	tester := tests.NewTester(client.SpeedtestClient, tests.DefaultDLSizes, tests.DefaultULSizes, false, false)
	downloadMbps := tester.Download(client.Server)
	uploadMbps := tester.Upload(client.Server)
	ping := client.Server.Latency
	result["download"] = downloadMbps
	result["upload"] = uploadMbps
	result["ping"] = ping
	log.Debugf("Speedtest results: %s", result)
	return result
}
