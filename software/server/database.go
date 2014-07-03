// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package main

import (
	"fmt"
	"github.com/influxdb/influxdb-go"
	"log"
	"net/http"
)

var dbChannel = make(chan MeterReading)
var client *influxdb.Client

func InitDBConnector() {
	var err error
	client, err = influxdb.NewClient(&influxdb.ClientConfig{
		Host: fmt.Sprintf("%s:%d", Cfg.InfluxDB.Host,
			Cfg.InfluxDB.Port),
		Username:   Cfg.InfluxDB.User,
		Password:   Cfg.InfluxDB.Pass,
		Database:   Cfg.InfluxDB.Database,
		HttpClient: http.DefaultClient,
	})
	if err != nil {
		log.Fatal("Cannot create InfluxDB client:", err.Error())
	}
	client.DisableCompression()
	dbs, err := client.GetDatabaseList()
	if err != nil {
		log.Fatal("Cannot retrieve list of InfluxDB databases:", err.Error())
	}
	foundDatabase := false
	for idx := range dbs {
		name := dbs[idx]["name"]
		log.Printf("found database %s", name)
		if name == Cfg.InfluxDB.Database {
			foundDatabase = true
		}
	}
	if !foundDatabase {
		log.Printf("Did not find database %s - attempting to create it",
			Cfg.InfluxDB.Database)
		if err := client.CreateDatabase(Cfg.InfluxDB.Database); err != nil {
			log.Fatal("Failed to create database ", Cfg.InfluxDB.Database)
		}
	}
}

func DBPusher() {
	if client == nil {
		log.Fatal("InfluxDB client not initialized - aborting")
	}
	log.Println("Ready to push values into the database")
	for {
		meterreading, ok := <-dbChannel
		if !ok {
			log.Fatal("Cannot read from internal channel - aborting")
		}
		//log.Printf("Pushing reading %v", reading)
		series := &influxdb.Series{
			Name:    meterreading.MeterID,
			Columns: []string{"time", "frequency"},
			Points: [][]interface{}{
				[]interface{}{meterreading.Reading.Timestamp.Unix() * 1000,
					meterreading.Reading.Value},
			},
		}
		if err := client.WriteSeries([]*influxdb.Series{series}); err != nil {
			log.Printf("Failed to store data: ", err)
		}
	}
}
