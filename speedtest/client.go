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
	"github.com/prometheus/common/log"
	"github.com/zpeters/speedtest/sthttp"
	"github.com/zpeters/speedtest/tests"
)

// Client defines the Speedtest client
type Client struct {
	Server         sthttp.Server
	Config         sthttp.Config
	AllServers     []sthttp.Server
	ClosestServers []sthttp.Server
}

// NewClient defines a new client for Speedtest
func NewClient(configURL string, serversURL string) (*Client, error) {
	log.Debugf("New Speedtest client")
	var testServer sthttp.Server

	log.Debugf("Retrieve Speedtest configuration from: %s", configURL)
	config, err := sthttp.GetConfig(configURL)
	if err != nil {
		return nil, err
	}
	// sthttp.CONFIG = config

	log.Debugf("Retrieve all servers")
	var allServers []sthttp.Server
	allServers, err = sthttp.GetServers(serversURL, "")
	if err != nil {
		return nil, err
	}

	log.Debugf("Retrieve closest servers")
	closestServers := sthttp.GetClosestServers(allServers, config.Lat, config.Lon)
	log.Debugf("Find test server")
	testServer = sthttp.GetFastestServer(closestServers)

	return &Client{
		Server:         testServer,
		Config:         config,
		AllServers:     allServers,
		ClosestServers: closestServers,
	}, nil
}

func (client *Client) NetworkMetrics() map[string]float64 {
	result := map[string]float64{}
	tester := tests.NewTester(tests.DefaultDLSizes, tests.DefaultULSizes, false, false)
	downloadMbps := tester.Download(client.Server)
	uploadMbps := tester.Upload(client.Server)
	ping := client.Server.Latency
	result["download"] = downloadMbps
	result["upload"] = uploadMbps
	result["ping"] = ping
	log.Debugf("Speedtest results: %s", result)
	return result
}
