// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defluxio

import (
	"testing"
)

func getValidMeters() (valid_meters []Meter) {
	valid_meters = []Meter{
		{
			Rank:      0,
			ID:        "valid",
			Key:       "valid",
			Name:      "Valid",
			Location:  "Here.",
			Cache:     MakeReadingCache(10),
			CacheSize: 10,
		},
		{
			Rank:      1,
			ID:        "valid",
			Key:       "valid",
			Name:      "Valid",
			Location:  "Here.",
			Cache:     MakeReadingCache(10),
			CacheSize: 10,
		},
	}
	return valid_meters
}

func TestMeterValid(t *testing.T) {
	valid := Meter{
		ID:       "valid",
		Key:      "valid",
		Name:     "Valid",
		Location: "Here.",
	}
	invalidID := Meter{
		Key:      "valid",
		Name:     "Valid",
		Location: "Here.",
	}
	invalidKey := Meter{
		ID:       "valid",
		Name:     "Valid",
		Location: "Here.",
	}
	invalidName := Meter{
		ID:       "valid",
		Key:      "valid",
		Location: "Here.",
	}
	invalidLocation := Meter{
		ID:   "valid",
		Key:  "valid",
		Name: "Valid",
	}
	if !valid.IsValid() {
		t.Error("Expected the good meter entry to be valid")
	}
	if invalidID.IsValid() || invalidKey.IsValid() ||
		invalidName.IsValid() || invalidLocation.IsValid() {
		t.Error("Expected the bad meter entry to be invalid")
	}
}

func TestMeterEquality(t *testing.T) {
	valid_meters := getValidMeters()
	if valid_meters[0] != valid_meters[0] {
		t.Error("Same meter is assumed to be different")
	}
	if valid_meters[0] == valid_meters[1] {
		t.Error("Different meters are assumed to be equal")
	}
}

func TestMeterCollectionValid(t *testing.T) {
	valid_meters := getValidMeters()
	mc0 := MeterCollection{Meters: valid_meters}
	if !mc0.IsValid() {
		t.Error("Meter Collection is not valid")
	}

	invalid_meters := append(valid_meters, valid_meters[0])
	mc1 := MeterCollection{Meters: invalid_meters}
	if mc1.IsValid() {
		t.Error("Meter Collection was accepted although it was not valid")
	}
}

func TestMeterCaching(t *testing.T) {
	valid_meters := getValidMeters()
	r1 := Reading{Value: 1.0}
	valid_meters[0].AppendReading(r1)
	last_reading, _ := valid_meters[0].Cache.LastReading()
	if r1 != last_reading {
		t.Error("Meter does not provide last reading correctly.")
	}
}
