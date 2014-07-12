// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package defluxio

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

func NewDBClient(serverConfig *ServerConfigurationData) (*DBClient, error) {
	retval := new(DBClient)
	var err error
	retval.client, err = influxdb.NewClient(&influxdb.ClientConfig{
		Host: fmt.Sprintf("%s:%d", serverConfig.InfluxDB.Host,
			serverConfig.InfluxDB.Port),
		Username:   serverConfig.InfluxDB.User,
		Password:   serverConfig.InfluxDB.Pass,
		Database:   serverConfig.InfluxDB.Database,
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
		if name == serverConfig.InfluxDB.Database {
			foundDatabase = true
		}
	}
	if !foundDatabase {
		log.Printf("Did not find database %s - attempting to create it",
			serverConfig.InfluxDB.Database)
		if err := retval.client.CreateDatabase(serverConfig.InfluxDB.Database); err != nil {
			return nil, fmt.Errorf("Failed to create database ", serverConfig.InfluxDB.Database)
		}
	}
	retval.DbChannel = make(chan MeterReading)
	//log.Printf("Ready to push data into the database.")
	return retval, nil
}

func (dbc DBClient) MkDBPusher() func() {
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

func (dbc DBClient) GetAllValues(meterID string) ([]MeterReading, error) {
	querystr := fmt.Sprintf("select time, frequency from %s", meterID)
	series, err := dbc.client.Query(querystr)
	if err != nil {
		return nil, fmt.Errorf("Failed query:", err.Error())
	}
	fmt.Printf("%v", series)
	return nil, nil
}
