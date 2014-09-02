// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defluxio

import (
	"testing"
)

func TestCacheAdding(t *testing.T) {
	r1 := Reading{Value: 1.0}
	r2 := Reading{Value: 2.0}
	r3 := Reading{Value: 3.0}

	rc := MakeReadingCache(10)
	rc.AddReading(r1)
	rc.AddReading(r2)
	rc.AddReading(r3)
	if rc.NumElements() != 3 {
		t.Errorf("Expected %d elements, got %d.", 3, rc.NumElements())
	}
}

func TestCacheGetLast(t *testing.T) {
	rc := MakeReadingCache(10)
	if _, e := rc.LastReading(); e == nil {
		t.Errorf("Expected error when querying empty cache, but got none")
	}
	r1 := Reading{Value: 1.0}
	r2 := Reading{Value: 2.0}
	rc.AddReading(r1)
	if r, _ := rc.LastReading(); r != r1 {
		t.Errorf("Retrieving last reading: Expected %s, got %s", r1, r)
	}
	rc.AddReading(r2)
	if r, _ := rc.LastReading(); r != r2 {
		t.Errorf("Retrieving last reading: Expected %s, got %s", r2, r)
	}
	for i := 0; i < 100; i++ {
		rc.AddReading(r2)
	}
	if r, _ := rc.LastReading(); r != r2 {
		t.Errorf("Retrieving last reading: Expected %s, got %s", r2, r)
	}
	if l := rc.NumElements(); l != 10 {
		t.Errorf("Counting number of elements: Expected %s, got %s", 10, l)
	}
}

func TestCacheGetAll(t *testing.T) {
	rc := MakeReadingCache(10)
	if readings := rc.AllReadings(); len(readings) != 0 {
		t.Errorf("Expected empty cache, but got some elements: %s", readings)
	}
	r1 := Reading{Value: 1.0}
	for i := 0; i < 100; i++ {
		rc.AddReading(r1)
	}
	if readings := rc.AllReadings(); len(readings) != 10 {
		t.Errorf("Expected 10 elements in cache, but got %d", len(readings))
	}
}
