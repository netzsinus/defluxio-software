// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package defluxio

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

// TODO: Move MeterID out of MeterReading - too much duplication
// of the meter name. New type "MeterTimeseries?"
type MeterReading struct {
	MeterID string
	Reading Reading
}

// ByTimestamp implements sort.Interface for []MeterReading
// based on the timestamp field of a reading.
type ByTimestamp []MeterReading

func (a ByTimestamp) Len() int      { return len(a) }
func (a ByTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool {
	return a[i].Reading.Timestamp.Unix() < a[j].Reading.Timestamp.Unix()
}

type DBClient struct {
	client       client.Client
	serverconfig *InfluxDBConfig
}

func NewDBClient(serverConfig *InfluxDBConfig) (*DBClient, error) {
	retval := new(DBClient)
	var err error
	retval.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", serverConfig.Host,
			serverConfig.Port),
		Username: serverConfig.User,
		Password: serverConfig.Pass,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create InfluxDB client: %s", err.Error())
	}

	// Save config for later use
	retval.serverconfig = serverConfig
	return retval, nil
}

func (dbc DBClient) MkDBPusher(dbchannel chan MeterReading) (func(), error) {
	log.Println("Getting list of databases:")
	response, err := dbc.client.Query(
		client.Query{
			Command:  "SHOW DATABASES",
			Database: dbc.serverconfig.Database,
		})
	if err != nil {
		log.Fatal(err)
	}
	if err == nil && response.Error() != nil {
		return nil, fmt.Errorf("Cannot retrieve list of InfluxDB databases: %s", response.Error())
	}
	log.Printf("Show databases response: %v", response.Results)
	foundDatabase := false
	// holy shit this is ugly
	for _, result := range response.Results {
		for _, row := range result.Series {
			if row.Name == "databases" {
				for _, values := range row.Values {
					for _, database := range values {
						log.Printf("found database: %s", database)
						if database == dbc.serverconfig.Database {
							foundDatabase = true
						}
					}
				}
			}
		}
	}

	if !foundDatabase {
		log.Fatalf("Did not find database \"%s\" - please create it",
			dbc.serverconfig.Database)
	}
	return func() {
		for {
			log.Printf("Not implemented: Push1 callback")
			_, ok := <-dbchannel
			meterreading, ok := <-dbchannel
			if !ok {
				log.Fatal("Cannot read from internal channel - aborting")
			}
			log.Printf("Pushing reading %v", meterreading.Reading)
			// Create a new point batch
			bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "frequency",
				Precision: "ms",
			})

			// Create a point and add to batch
			tags := map[string]string{"meterid": meterreading.MeterID}
			fields := map[string]interface{}{
				"value":     meterreading.Reading.Value,
				"timestamp": meterreading.Reading.Timestamp,
			}
			pt, err := client.NewPoint("frequency", tags, fields, time.Now())
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
			bp.AddPoint(pt)

			// Write the batch
			dbc.client.Write(bp)

			//series := &influxdb.Series{
			//	Name:    meterreading.MeterID,
			//	Columns: []string{"time", "frequency"},
			//	Points: [][]interface{}{
			//		[]interface{}{meterreading.Reading.Timestamp.Unix() * 1000,
			//			meterreading.Reading.Value},
			//	},
			//}
			//if err := dbc.client.WriteSeries([]*influxdb.Series{series}); err != nil {
			//	log.Printf("Failed to store data: %s", err)
			//}
		}
	}, nil
}

func (dbc DBClient) points2meterreadings(name string,
	points [][]interface{}) (retval []MeterReading) {
	for _, val := range points {
		timestamp := time.Unix(0, int64(val[0].(float64))*
			int64(time.Millisecond))
		frequency := val[2].(float64)
		//fmt.Printf("timestamp %v: %f\n", timestamp, frequency)
		retval = append(retval, MeterReading{name, Reading{timestamp, frequency}})
	}
	return retval
}

func (dbc DBClient) GetFrequenciesBetween(meterID string,
	start time.Time, end time.Time) (retval []MeterReading, err error) {
	log.Printf("GetFrequenciesBetween not implemented!")
	return nil, nil
	//querystr := fmt.Sprintf("select time, frequency from %s where time > %ds and time < %ds", meterID, start.Unix(), end.Unix())
	//series, err := dbc.client.Query(querystr)
	//if err != nil {
	//	return retval, fmt.Errorf("Failed query: %s", err.Error())
	//}
	//if len(series) == 0 {
	//	return retval, fmt.Errorf("No dataset received from database.")
	//}
	//// Debug: Print raw data points
	////fmt.Printf("%#v\n", series[0].Points)
	//retval = dbc.points2meterreadings(series[0].Name, series[0].Points)
	//return retval, nil
}

func (dbc DBClient) GetLastFrequencies(meterID string, amount int) ([]MeterReading, error) {
	retval := []MeterReading{}
	//querystr := fmt.Sprintf("select time, frequency from %s limit %d",
	//	meterID, amount)
	//series, err := dbc.client.Query(querystr)
	//if err != nil {
	//	return retval, fmt.Errorf("Failed query: %s", err.Error())
	//}
	//if len(series[0].Points) != amount {
	//	return retval, fmt.Errorf("Received invalid number of readings: Expected %d, got ", len(series))
	//}
	//// Debug: Print raw data points
	////fmt.Printf("%#v\n", series[0].Points)
	//retval = dbc.points2meterreadings(series[0].Name, series[0].Points)
	log.Printf("GerLastFrequencies not implemented")
	return retval, nil
}

func (dbc DBClient) GetLastFrequency(meterID string) (MeterReading,
	error) {
	readings, error := dbc.GetLastFrequencies(meterID, 1)
	if error != nil {
		return MeterReading{}, error
	} else {
		return readings[0], error
	}
}
