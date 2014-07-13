// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package defluxio

import (
	"fmt"
	"github.com/influxdb/influxdb-go"
	"log"
	"net/http"
	"time"
)

type MeterReading struct {
	MeterID string
	Reading Reading
}

type DBClient struct {
	client       *influxdb.Client
	serverconfig *InfluxDBConfig
}

func NewDBClient(serverConfig *InfluxDBConfig) (*DBClient, error) {
	retval := new(DBClient)
	var err error
	retval.client, err = influxdb.NewClient(&influxdb.ClientConfig{
		Host: fmt.Sprintf("%s:%d", serverConfig.Host,
			serverConfig.Port),
		Username:   serverConfig.User,
		Password:   serverConfig.Pass,
		Database:   serverConfig.Database,
		HttpClient: http.DefaultClient,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create InfluxDB client:", err.Error())
	}
	retval.client.DisableCompression()
	// Save config for later use
	retval.serverconfig = serverConfig
	log.Println("Created DB client.")
	return retval, nil
}

func (dbc DBClient) MkDBPusher(dbchannel chan MeterReading) (func(), error) {
	log.Println("Getting list of databases.")
	dbs, err := dbc.client.GetDatabaseList()
	if err != nil {
		log.Println("acquired: ", len(dbs))

		return nil, fmt.Errorf("Cannot retrieve list of InfluxDB databases:", err.Error())
	}
	foundDatabase := false
	for idx := range dbs {
		name := dbs[idx]["name"]
		log.Printf("found database %s", name)
		if name == dbc.serverconfig.Database {
			foundDatabase = true
		}
	}
	if !foundDatabase {
		log.Printf("Did not find database %s - attempting to create it",
			dbc.serverconfig.Database)
		if err := dbc.client.CreateDatabase(dbc.serverconfig.Database); err != nil {
			return nil, fmt.Errorf("Failed to create database ", dbc.serverconfig.Database)
		}
	}

	return func() {
		for {
			meterreading, ok := <-dbchannel
			if !ok {
				log.Fatal("Cannot read from internal channel - aborting")
			}
			//log.Printf("Pushing reading %v", meterreading.Reading)
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
	}, nil
}

func (dbc DBClient) GetLastFrequency(meterID string) (MeterReading, error) {
	retval := MeterReading{}
	querystr := fmt.Sprintf("select time, frequency from %s limit 1", meterID)
	series, err := dbc.client.Query(querystr)
	if err != nil {
		return retval, fmt.Errorf("Failed query:", err.Error())
	}
	// TODO: Check length of return value
	fmt.Printf("%#v\n", series)
	fmt.Printf("%i\n", len(series))
	// TODO: Cleanup
	for key, val := range series {
		fmt.Printf("%i: %s\n", key, val)
		timestamp := time.Unix(0, int64(val.Points[0][0].(float64))*
			int64(time.Millisecond))
		frequency := val.Points[0][2].(float64)
		fmt.Printf("timestamp %v: %f\n", timestamp, frequency)
		retval = MeterReading{val.Name, Reading{timestamp, frequency}}
	}
	return retval, nil
}
