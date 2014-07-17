// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gonium/defluxio"
	"log"
	"sort"
	"strings"
	"time"
)

var configFile = flag.String("config", "defluxio-exporter.conf", "configuration file")
var meterID = flag.String("meter", "", "ID of the meter to query")
var cfg *defluxio.ExporterConfigurationData
var dbclient *defluxio.DBClient

func init() {
	flag.Parse()
	if strings.EqualFold(*meterID, "") {
		log.Fatal("You must specify the meter ID (i.e. -meter=foometer)")
	}
	var err error
	cfg, err = defluxio.LoadExporterConfiguration(*configFile)
	if err != nil {
		log.Fatal("Error loading configuration: ", err.Error())
	}
	dbclient, err = defluxio.NewDBClient(&cfg.InfluxDB)
	if err != nil {
		log.Fatal("Cannot initialize database client:", err.Error())
	}
}

func main() {
	fmt.Printf("Attempting to export from meter %s\n", *meterID)
	//result, err := dbclient.GetLastFrequency(*meterID)
	//if err != nil {
	//	log.Fatal("Failed to query database: ", err.Error())
	//}
	//fmt.Printf("On %v, the frequency was recorded as %f\n",
	//	result.Reading.Timestamp, result.Reading.Value)
	//meterReadings, err := dbclient.GetLastFrequencies(*meterID, 10)
	//if err != nil {
	//	log.Fatal("Failed to query database: ", err.Error())
	//}
	//for _, element := range meterReadings {
	//	fmt.Printf("%v: %f\n", element.Reading.Timestamp,
	//		element.Reading.Value)
	//}

	// Hack for testing
	timeReadings, terr := dbclient.GetFrequenciesBetween(*meterID,
		time.Unix(1405525188, 0), time.Unix(1405525439, 0))
	if terr != nil {
		log.Fatal("Failed to query database: ", terr.Error())
	}
	sort.Sort(defluxio.ByTimestamp(timeReadings))
	fmt.Printf("timestamp\treading\n")
	for _, element := range timeReadings {
		fmt.Printf("%d\t%f\n", element.Reading.Timestamp.Unix(),
			element.Reading.Value)
	}
	//TODO: Write exporter. The exporter should use a TSV format with
	// unix timestamp \t value. Header: "timestamp\treading"
}
