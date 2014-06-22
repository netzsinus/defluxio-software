// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.

package main

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

type ErrorMessage struct {
	Id      string
	Message string
}

/* Accepts a new reading. Format: {"Timestamp":<ISO8601>,"Value":342.2}
 */
func submitReading(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// extract and check credentials.
	// TODO: Refactore this into a separate middleware for chaining.
	vars := mux.Vars(r)
	meterid := vars["meter"]
	apikey := r.Header["X-Api-Key"][0]
	credentials := fmt.Sprintf("%s:%s", meterid, apikey)
	invalidCredentials := true
	for _, key := range Cfg.API.Keys {
		if key == credentials {
			invalidCredentials = false
			break
		}
	}
	if invalidCredentials {
		errormessage := ErrorMessage{"invalidcredentials", "given credentials are not valid"}
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
		errormessage := ErrorMessage{"invalidformat", "Cannot decode data - invalid format"}
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
		errormessage := ErrorMessage{"timedeviation", "Timestamp deviates too much"}
		errorbytes, _ := json.Marshal(errormessage)
		http.Error(w,
			string(errorbytes),
			http.StatusBadRequest)
		return
	}

	// Check if value is within plausibility range
	if reading.Value < 47.5 || reading.Value > 52 {
		errormessage := ErrorMessage{"timedeviation", "Reading value is not plausible"}
		errorbytes, _ := json.Marshal(errormessage)
		http.Error(w,
			string(errorbytes),
			http.StatusBadRequest)
		return
	}

	// round reading.Value to 4 decimal places
	reading.Value = float64(int(reading.Value*10000)) / 10000
	log.Println("Accepted", meterid, ":", reading.Timestamp, "-", reading.Value)

	// finally: wrap everything again and forward update to all connected clients.
	updateMessage, uerr := json.Marshal(reading)
	if uerr != nil {
		panic(uerr)
	}
	h.broadcast <- []byte(updateMessage)
}
