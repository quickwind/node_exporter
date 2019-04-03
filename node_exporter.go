// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "net/http/pprof"
	"sort"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/prometheus/node_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		metricsOutput = kingpin.Flag(
			"metrics.output",
			"Output file for the metrics.",
		).Default("./metrics.out").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting node_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	nc, err := collector.NewNodeCollector()
	if err != nil {
		log.Fatalf("couldn't create collector: %s", err)
	}

	log.Infof("Enabled collectors:")
	collectors := []string{}
	for n := range nc.Collectors {
		collectors = append(collectors, n)
	}
	sort.Strings(collectors)
	for _, n := range collectors {
		log.Infof(" - %s", n)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("node_exporter"))
	if err := r.Register(nc); err != nil {
		log.Fatalf("couldn't register node collector: %s", err)
	}

	log.Infoln("Collecting metrics...")
	if err := prometheus.WriteToTextfile(*metricsOutput, r); err != nil {
		log.Fatalf("collection failed: %s", err)
	}
	log.Infoln("Done")
}
