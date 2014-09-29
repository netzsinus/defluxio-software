// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.

package defluxio

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

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
func MkSubmitReadingHandler(dbchannel chan MeterReading, serverConfig *ServerConfiguration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// extract and check credentials.
		// TODO: Refactore this into a separate middleware for chaining.
		vars := mux.Vars(r)
		meterid := vars["meter"]
		apikey := r.Header["X-Api-Key"][0]
		//credentials := fmt.Sprintf("%s:%s", meterid, apikey)
		invalidCredentials := true
		for _, meter := range serverConfig.Meters {
			if meter.ID == meterid && meter.Key == apikey {
				invalidCredentials = false
				break
			}
		}
		//for _, key := range serverConfig.API.Keys {
		//	if key == credentials {
		//		invalidCredentials = false
		//		break
		//	}
		//}
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

		if dbchannel != nil {
			// Push the new reading into the database
			meterReading := MeterReading{meterid, reading}
			dbchannel <- meterReading
		}

		// Update the meter cache
		for idx := range serverConfig.Meters {
			if serverConfig.Meters[idx].ID == meterid {
				serverConfig.Meters[idx].AppendReading(reading)
				break
			}
		}

		// Push the last reading to the clients, but only for the active
		// meter with the lowest rank. A meter is active if it has been
		// sending data within the last five seconds. Therefore: Sort the
		// meters by the last reading timestamp, then pick the one with the
		// lowest rank. If the current reading has been sent by this meter,
		// push the reading to the websocket clients.
		// Note: This is a costly operation, but right now we don't have
		// many meters. If the number of meters increases, change this
		// algorithm.
		if BestMeter.ID == meterid {
			// The update we have received is from the best meter
			//log.Println("Received update from highest ranking meter", BestMeter.ID)
			// wrap everything again and forward update to all connected clients.
			updateMessage, uerr := json.Marshal(reading)
			if uerr != nil {
				log.Printf("Cannot marshal update message for websocket consumers: ", uerr)
			}
			H.broadcast <- []byte(updateMessage)
		}

	})
}
