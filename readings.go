// (C) 2014 Mathias Dalheimer <md@gonium.net>.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defluxio

import (
	"container/ring"
	"errors"
	"fmt"
	"time"
)

// Reading type. This is the type we use to represent a single
// measurement. Currently, this only supports frequency measurements -
// but it can be extended in the future.

type Reading struct {
	Timestamp time.Time
	Value     float64
}

func (r Reading) String() (retval string) {
	return fmt.Sprintf("(%s: %2.6f)", r.Timestamp, r.Value)
}

// The cache implements a cache of readings for a single meter.
type ReadingCache struct {
	Cache *ring.Ring
}

func MakeReadingCache(size int) (r ReadingCache) {
	r = ReadingCache{
		Cache: ring.New(size),
	}
	return r
}

func (c *ReadingCache) AddReading(r Reading) {
	c.Cache.Value = r
	c.Cache = c.Cache.Next()
}

func (c *ReadingCache) LastReading() (r Reading, err error) {
	foo := c.Cache.Prev()
	if foo.Value == nil {
		return r, errors.New("no element in cache")
	} else {
		r = foo.Value.(Reading)
		return r, nil
	}
}

func (c *ReadingCache) AllReadings() (r []Reading) {
	r = make([]Reading, 0, 10)
	i := 0
	c.Cache.Do(func(x interface{}) {
		if x != nil {
			//fmt.Printf("adding element %d: %s\n", i, x.(Reading))
			r = append(r, x.(Reading))
			i++
		}
	})
	//fmt.Printf("Returning %d elements\n", len(r))
	return r
}

func (c *ReadingCache) NumElements() (retval int) {
	retval = 0
	c.Cache.Do(func(x interface{}) {
		if x != nil {
			retval += 1
		}
	})
	return retval
}

func (c ReadingCache) String() (retval string) {
	retval = "["
	c.Cache.Do(func(x interface{}) {
		if x != nil {
			retval = fmt.Sprintf("%s %s", retval, x)
		}
	})
	retval = retval + "]"
	return retval
}
