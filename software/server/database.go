// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package main

import (
	"fmt"
	"github.com/influxdb/influxdb-go"
	"log"
	"net/http"
)

type MeterReading struct {
	MeterID string
	Reading Reading
}

type DBClient struct {
	client    *influxdb.Client
	DbChannel chan MeterReading
}

func NewDBClient() (*DBClient, error) {
	retval := new(DBClient)
	var err error
	retval.client, err = influxdb.NewClient(&influxdb.ClientConfig{
		Host: fmt.Sprintf("%s:%d", Cfg.InfluxDB.Host,
			Cfg.InfluxDB.Port),
		Username:   Cfg.InfluxDB.User,
		Password:   Cfg.InfluxDB.Pass,
		Database:   Cfg.InfluxDB.Database,
		HttpClient: http.DefaultClient,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create InfluxDB client:", err.Error())
	}
	retval.client.DisableCompression()
	dbs, err := retval.client.GetDatabaseList()
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve list of InfluxDB databases:", err.Error())
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
		if err := retval.client.CreateDatabase(Cfg.InfluxDB.Database); err != nil {
			return nil, fmt.Errorf("Failed to create database ", Cfg.InfluxDB.Database)
		}
	}
	retval.DbChannel = make(chan MeterReading)
	//log.Printf("Ready to push data into the database.")
	return retval, nil
}

func (dbc DBClient) mkDBPusher() func() {
	//	if client == nil {
	//		log.Fatal("InfluxDB client not initialized - aborting")
	//	}
	fmt.Printf("mkDBPusher")
	for {
		meterreading, ok := <-dbc.DbChannel
		if !ok {
			log.Fatal("Cannot read from internal channel - aborting")
		}
		log.Printf("Pushing reading %v", meterreading.Reading)
		series := &influxdb.Series{
			Name:    meterreading.MeterID,
			Columns: []string{"time", "frequency"},
			Points: [][]interface{}{
				[]interface{}{meterreading.Reading.Timestamp.Unix() * 1000,
					meterreading.Reading.Value},
			},
		}
		if err := dbc.client.WriteSeries([]*influxdb.Series{series}); err != nil {
			log.Printf("Failed to store data: ", err)
		}
	}
}
