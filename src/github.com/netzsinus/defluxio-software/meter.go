// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defluxio

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type Meter struct {
	Rank      uint16
	ID        string
	Key       string
	Name      string
	Location  string
	Cache     *ReadingCache `json:"-"` // do not export
	CacheSize uint32
}

type Meters []*Meter

var BestMeter *Meter //TODO: Refactor this to be private

func isEmpty(s string) bool {
	return s == ""
}

func (m *Meter) IsValid() bool {
	return !isEmpty(m.ID) && !isEmpty(m.Key) && !isEmpty(m.Name) && !isEmpty(m.Location)
}

func (m *Meter) AppendReading(r Reading) {
	m.Cache.AddReading(r)
}

func (mc Meters) StartBestMeterUpdater(timeout time.Duration) {
	BestMeter = mc.SelectBestMeter(timeout)
	ticker := time.NewTicker(time.Second * timeout / 2)
	go func() {
		for _ = range ticker.C {
			BestMeter = mc.SelectBestMeter(timeout)
			if BestMeter == nil {
				log.Printf("No active meters available, deactivating BestMeter")
			}
		}
	}()
}

func (mc Meters) GetBestMeter() (m *Meter) {
	if BestMeter == nil {
		// previously, no best meter was selected, but maybe we just
		// received an update.
		BestMeter = mc.SelectBestMeter(time.Second * 5)
		if BestMeter != nil {
			log.Println("Activated BestMeter - connected web clients will " +
				"receive updates")
		}
	}
	return BestMeter
}

func (mc Meters) SelectBestMeter(timeout time.Duration) (m *Meter) {
	sort.Sort(ByRank{mc})
	for _, m := range mc {
		if b, _ := m.activeWithinLast(time.Second * timeout); b {
			return m
		}
	}
	return nil
}

func (m *Meter) activeWithinLast(deadline time.Duration) (bool, error) {
	r, err := m.Cache.LastReading()
	if err != nil {
		return false, err
	} else {
		return time.Since(r.Timestamp) < deadline, nil
	}
}

func (mc Meters) IsValid() bool {
	ranks := []uint16{}
	// Check individual meters
	for _, m := range mc {
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

func (ms Meters) String() string {
	retval := "Available Meters:"
	for _, m := range ms {
		retval = fmt.Sprintf("%s\n* %s - %s (%d)", retval, m.ID, m.Name, m.Rank)
	}
	return retval
}

/************  Make meters sortable ******************/

func (m Meters) Len() int      { return len(m) }
func (m Meters) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

type ByRank struct{ Meters }
type ByName struct{ Meters }
type ByLastUpdate struct{ Meters }

func (s ByRank) Less(i, j int) bool {
	return s.Meters[i].Rank < s.Meters[j].Rank
}

func (s ByName) Less(i, j int) bool {
	return s.Meters[i].Name < s.Meters[j].Name
}

func (s ByLastUpdate) Less(i, j int) bool {
	it, ierr := s.Meters[i].Cache.LastReading()
	jt, jerr := s.Meters[j].Cache.LastReading()
	if ierr != nil || jerr != nil {
		log.Println("Failed to acquire last reading in sorting predicate.")
		return false
	} else {
		return it.Timestamp.Unix() < jt.Timestamp.Unix()
	}
}

/************  Helpers below  ******************/

func intInSlice(i uint16, s []uint16) bool {
	for _, b := range s {
		if b == i {
			return true
		}
	}
	return false
}
