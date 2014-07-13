// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.

package defluxio

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Reading struct {
	Timestamp time.Time
	Value     float64
}

type APIClientErrorMessage struct {
	Id      string
	Message string
}

type StatusCode int32

const (
	STATUS_OK               StatusCode = 200
	STATUS_INTERNAL_ERROR              = 500
	STATUS_NO_DB_CONNECTION            = 501 // TODO: Implement keepalive
	STATUS_NO_UPDATES                  = 502
)

type ServerStatusData struct {
	Code           StatusCode
	LastValueAdded time.Time
	Message        string
}

var lastUpdate time.Time

func ServerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Check when the last value was received
	maxTimeDeviation, _ := time.ParseDuration("30s")
	deviation := lastUpdate.Sub(time.Now())
	if (deviation > maxTimeDeviation) || (deviation < -maxTimeDeviation) {
		statusmessage := ServerStatusData{
			STATUS_NO_UPDATES,
			lastUpdate,
			"No updates received."}
		statusbytes, _ := json.Marshal(statusmessage)
		w.Write(statusbytes)
		return
	}
	// TODO: Add other checks here.

	// Finally: No check failed, return OK.
	statusmessage := ServerStatusData{
		STATUS_OK,
		lastUpdate,
		"All is good."}
	statusbytes, _ := json.Marshal(statusmessage)
	w.Write(statusbytes)
}

/* Accepts a new reading. Format: {"Timestamp":<ISO8601>,"Value":342.2}
 */
func MkSubmitReadingHandler(dbclient *DBClient, serverConfig *ServerConfigurationData) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// extract and check credentials.
		// TODO: Refactore this into a separate middleware for chaining.
		vars := mux.Vars(r)
		meterid := vars["meter"]
		apikey := r.Header["X-Api-Key"][0]
		credentials := fmt.Sprintf("%s:%s", meterid, apikey)
		invalidCredentials := true
		for _, key := range serverConfig.API.Keys {
			if key == credentials {
				invalidCredentials = false
				break
			}
		}
		if invalidCredentials {
			errormessage := APIClientErrorMessage{"invalidcredentials", "given credentials are not valid"}
			errorbytes, _ := json.Marshal(errormessage)
			http.Error(w,
				string(errorbytes),
				http.StatusUnauthorized)
			return
		}

		// Decode the body and check for errors.
		decoder := json.NewDecoder(r.Body)
		var reading Reading
		err := decoder.Decode(&reading)
		if err != nil {
			errormessage := APIClientErrorMessage{"invalidformat", "Cannot decode data - invalid format"}
			errorbytes, _ := json.Marshal(errormessage)
			http.Error(w,
				string(errorbytes),
				http.StatusBadRequest)
			return
		}

		// Check whether timestamp is current
		maxTimeDeviation, _ := time.ParseDuration("30s")
		deviation := reading.Timestamp.Sub(time.Now())
		if (deviation > maxTimeDeviation) || (deviation < -maxTimeDeviation) {
			errormessage := APIClientErrorMessage{"timedeviation", "Timestamp deviates too much"}
			errorbytes, _ := json.Marshal(errormessage)
			http.Error(w,
				string(errorbytes),
				http.StatusBadRequest)
			return
		}

		// Check if value is within plausibility range
		if reading.Value < 47.5 || reading.Value > 52 {
			errormessage := APIClientErrorMessage{"timedeviation", "Reading value is not plausible"}
			errorbytes, _ := json.Marshal(errormessage)
			http.Error(w, string(errorbytes), http.StatusBadRequest)
			return
		}

		// round reading.Value to 4 decimal places
		//	reading.Value = float64(int(reading.Value*1000)) / 1000
		//log.Println("Accepted", meterid, ":", reading.Timestamp, "-", reading.Value)
		lastUpdate = reading.Timestamp

		if serverConfig.InfluxDB.Enabled {
			// Push the new reading into the database
			meterReading := MeterReading{meterid, reading}
			if dbclient == nil {
				log.Fatal("DBClient not correctly initialized!")
			} else {
				dbclient.DbChannel <- meterReading
			}
		}

		// finally: wrap everything again and forward update to all connected clients.
		updateMessage, uerr := json.Marshal(reading)
		if uerr != nil {
			panic(uerr)
		}
		H.broadcast <- []byte(updateMessage)
	})
}
