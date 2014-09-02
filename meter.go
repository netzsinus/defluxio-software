// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defluxio

type Meter struct {
	Rank     uint16
	ID       string
	Key      string
	Name     string
	Location string
}

type MeterCollection struct {
	Meters []Meter
}

func isEmpty(s string) bool {
	return s == ""
}

func (m *Meter) IsValid() bool {
	return !isEmpty(m.ID) && !isEmpty(m.Key) && !isEmpty(m.Name) && !isEmpty(m.Location)
}

func intInSlice(i uint16, s []uint16) bool {
	for _, b := range s {
		if b == i {
			return true
		}
	}
	return false
}

func (mc *MeterCollection) IsValid() bool {
	ranks := []uint16{}
	// Check individual meters
	for _, m := range mc.Meters {
		if !m.IsValid() {
			return false
		}
		// has the rank been previously seen?
		if intInSlice(m.Rank, ranks) {
			return false
		} else {
			ranks = append(ranks, m.Rank)
		}
	}
	return true
}
