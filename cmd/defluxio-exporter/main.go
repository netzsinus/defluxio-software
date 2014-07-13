// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gonium/defluxio"
	"log"
)

var configFile = flag.String("config", "defluxio-exporter.conf", "configuration file")
var cfg *defluxio.ExporterConfigurationData
var dbclient *defluxio.DBClient

func init() {
	flag.Parse()
	var err error
	cfg, err = defluxio.LoadExporterConfiguration(*configFile)
	if err != nil {
		log.Fatal("Error loading configuration: ", err.Error())
	}
	dbclient, err = defluxio.NewDBClient(&cfg.InfluxDB)
	if err != nil {
		log.Fatal("Cannot initialize database client:", err.Error)
	}
}

func main() {
	log.Println("Exporter starting up.")
}
