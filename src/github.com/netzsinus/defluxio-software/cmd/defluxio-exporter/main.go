// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/netzsinus/defluxio-software"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	SecsPerDay int64 = 24 * 60 * 60
)

var configFile = flag.String("config", "defluxio-exporter.conf", "configuration file")
var mkConfigFile = flag.Bool("genconfig", false, "generate an example configuration file")
var meterID = flag.String("meter", "", "ID of the meter to query")
var startTimestamp = flag.Int64("start", 0,
	"data export start: first unix timestamp to export")
var endTimestamp = flag.Int64("end", 0,
	"data export end: last unix timestamp to export")
var exportDirectory = flag.String("dir", ".", "path to use for export")
var force = flag.Bool("force", false, "force export, overwriting existing files")
var verbose = flag.Bool("verbose", false, "verbose logging")
var lastday = flag.Bool("lastday", false, "do not export to the given end timestamp but up until one day before (midnight) - cron use.")
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
	// TODO: Implement new check - now writing to a directory.
	if !*force {
		if dir, err := IsDirectory(*exportDirectory); err != nil || !dir {
			log.Fatal("directory ", *exportDirectory, " does not exist - aborting.")
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

func IsDirectory(path string) (bool, error) {
	if fileInfo, err := os.Stat(path); err != nil {
		return false, err
	} else {
		return fileInfo.IsDir(), nil
	}
}

func getCanonicalFilename(ts int64) string {
	date := time.Unix(ts, 0)
	filename := fmt.Sprintf("%04d%02d%02d.txt", date.Year(), date.Month(), date.Day())
	return filepath.Join(*exportDirectory, filename)
}

func getStartOfDayTimestamp(ts int64) int64 {
	date := time.Unix(ts, 0)
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0,
		time.UTC).Unix()
}

func getEndOfDayTimestamp(ts int64) int64 {
	return getStartOfDayTimestamp(ts) + SecsPerDay - 1
}

func exportRange(start int64, end int64, filename string) (err error) {
	if !*force {
		if _, err := os.Stat(filename); err == nil {
			if *verbose {
				log.Printf("File %s already exists - skipping.", filename)
			}
			return err
		}
	}
	timeReadings, err := dbclient.GetFrequenciesBetween(*meterID,
		time.Unix(start, 0), time.Unix(end, 0))
	if err != nil {
		err = fmt.Errorf("Failed to query database: ", err.Error())
		return
	}
	sort.Sort(defluxio.ByTimestamp(timeReadings))

	log.Printf("Exporting %d readings [%d, %d] into file %s",
		len(timeReadings), start, end, filename)

	if tsve, err := defluxio.NewTsvExporter(filename); err != nil {
		return fmt.Errorf("Cannot create exporter for file %s ", filename)
	} else {
		if err = tsve.ExportDataset(timeReadings); err != nil {
			return fmt.Errorf("Failed to export dataset: %s ", err.Error())
		}
	}
	return
}

func main() {
	if *verbose {
		log.Printf("Attempting to export from meter %s\n", *meterID)
	}

	// Cron usage: In "lastday" mode, subtract a day from the end date.
	// The reasoning behind this is that a cronjob can use the date +%s
	// command to get the current timestamp and we export only up to the
	// day before - the current day is not yet fully passed.
	if *lastday {
		*endTimestamp -= SecsPerDay
	}

	for ts := getStartOfDayTimestamp(*startTimestamp); ts < getEndOfDayTimestamp(*endTimestamp); ts += SecsPerDay {
		exportRange(
			getStartOfDayTimestamp(ts),
			getEndOfDayTimestamp(ts),
			getCanonicalFilename(ts),
		)
	}

}
