// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/netzsinus/defluxio-software"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var configFile = flag.String("config", "defluxio-exporter.conf", "configuration file")
var mkConfigFile = flag.Bool("genconfig", false, "generate an example configuration file")
var meterID = flag.String("meter", "", "ID of the meter to query")
var startTimestamp = flag.Int64("start", 0,
	"data export start: first unix timestamp to export")
var endTimestamp = flag.Int64("end", 0,
	"data export end: last unix timestamp to export")
var exportFilename = flag.String("file", "defluxio-export.txt",
	"path to file to use for export")
var force = flag.Bool("force", false,
	"force export, overwriting existing files")
var cfg *defluxio.ExporterConfiguration
var dbclient *defluxio.DBClient

func init() {
	flag.Parse()
	if *mkConfigFile {
		log.Println("Creating default configuration in file " + *configFile)
		cfg := defluxio.MkDefaultExporterConfiguration()
		err := cfg.Save(*configFile)
		if err != nil {
			log.Fatal("Failed to create default configuration:", err.Error())
		}
		os.Exit(0)
	}
	if strings.EqualFold(*meterID, "") {
		log.Fatal("You must specify the meter ID (i.e. -meter=foometer)")
	}
	if *startTimestamp == 0 {
		log.Fatal("You must specify the start timestamp( i.e. -start=1405607436)")
	}
	if *endTimestamp == 0 {
		log.Fatal("You must specify the end timestamp( i.e. -end=1405607465)")
	}
	if *startTimestamp >= *endTimestamp {
		log.Fatal("start timestamp cannot be after end timestamp.")
	}
	if !*force {
		if _, err := os.Stat(*exportFilename); err == nil {
			log.Fatal("file ", *exportFilename, " exists - aborting.")
		}
	}
	var err error
	if cfg, err = defluxio.LoadExporterConfiguration(*configFile); err != nil {
		log.Fatal("Error loading configuration: ", err.Error())
	}
	if dbclient, err = defluxio.NewDBClient(&cfg.InfluxDB); err != nil {
		log.Fatal("Cannot initialize database client:", err.Error())
	}
}

func main() {
	log.Printf("Attempting to export from meter %s\n", *meterID)
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
	// TODO: Replace with real time.Unix foo from commandline
	timeReadings, terr := dbclient.GetFrequenciesBetween(*meterID,
		time.Unix(*startTimestamp, 0), time.Unix(*endTimestamp, 0))
	if terr != nil {
		log.Fatal("Failed to query database: ", terr.Error())
	}
	sort.Sort(defluxio.ByTimestamp(timeReadings))

	tsve, eerr := defluxio.NewTsvExporter(*exportFilename)
	if eerr != nil {
		log.Fatal("Cannot create exporter with file %s", *exportFilename)
	}
	if eerr = tsve.ExportDataset(timeReadings); eerr != nil {
		log.Fatal("Failed to export dataset: %s", eerr.Error())
	} else {
		log.Printf("Export finished successfully.")
	}
}
