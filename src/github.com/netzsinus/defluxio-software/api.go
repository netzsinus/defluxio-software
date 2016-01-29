// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.

package defluxio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
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
		if invalidCredentials {
			errormessage := APIClientErrorMessage{"invalidcredentials", "given credentials are not valid"}
			errorbytes, _ := json.Marshal(errormessage)
			http.Error(w,
				string(errorbytes),
				http.StatusUnauthorized)
			return
		}
		// http://stackoverflow.com/questions/23070876/reading-body-of-http-request-without-modifying-request-state
		body, _ := ioutil.ReadAll(r.Body)
		reader1 := bytes.NewBuffer(body)
		reader2 := bytes.NewBuffer(body)
		// Decode the body and check for errors.
		decoder := json.NewDecoder(reader1)
		var reading Reading
		format_invalid := false
		errormsgs := ""
		if err := decoder.Decode(&reading); err != nil {
			errormsgs = err.Error()
			// Standard RFC3339 decoding did not work - maybe this is an
			// embedded client. Is the timestamp represented as a unix
			// timestamp w/ milliseconds?
			type UnixReading struct {
				Timestamp float64
				Value     float64
			}
			var unixreading UnixReading
			decoder = json.NewDecoder(reader2)
			if uerr := decoder.Decode(&unixreading); uerr != nil {
				format_invalid = true
				errormsgs = fmt.Sprintf("%s; %s", errormsgs, uerr.Error())
			} else {
				// convert the float to a unix timestamp w/ nanosecond
				// resolution
				secs := math.Trunc(unixreading.Timestamp)
				nsecs := (unixreading.Timestamp - secs) * math.Pow(10, 9)
				reading.Timestamp = time.Unix(int64(secs), int64(nsecs))
				reading.Value = unixreading.Value
				log.Println("Reconstructed reading from unix timestamp: ", reading)
			}
			if format_invalid {
				msg := fmt.Sprintf("Cannot decode data - invalid format: %s", errormsgs)
				errormessage := APIClientErrorMessage{"invalidformat", msg}
				errorbytes, _ := json.Marshal(errormessage)
				http.Error(w,
					string(errorbytes),
					http.StatusBadRequest)
				return
			}
		}

		// Check whether timestamp is current
		maxTimeDeviation, _ := time.ParseDuration("30s")
		deviation := reading.Timestamp.Sub(time.Now())
		if (deviation > maxTimeDeviation) || (deviation < -maxTimeDeviation) {
			errormessage := APIClientErrorMessage{"timedeviation", "Timestamp deviates more than 30s from server time"}
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
