package bloom

import (
	"sync"
)

// Counting go holds counter rather than bit (0/1)
type CountingFilter struct {
	numCounters uint
	numHashes   uint
	counters    []uint8 //8-bit counters
	mu          sync.RWMutex
}

// Constructor for new Counting Filter
func NewCountingFilter(numCounters, numHashes uint) *CountingFilter {
	if numCounters < 1 {
		numCounters = 1
	}
	if numHashes < 1 {
		numHashes = 1
	}

	return &CountingFilter{
		numCounters: numCounters,
		numHashes:   numHashes,
		counters:    make([]uint8, numCounters),
	}
}

// Same init logic as Standard Filter
func NewCountingWithEstimatedParams(datasize int, fp float64) *CountingFilter {
	numCounters, numHashes := EstimateParameters(datasize, fp) // From ./bloomGo.go
	return NewCountingFilter(numCounters, numHashes)
}

// method for Counting Filter to get location --> modulo index for the counters array
func (cf *CountingFilter) location(hashes [4]uint64, i uint) uint {
	return uint(getLocation(hashes, i) % uint64(cf.numCounters)) // getLocation from ./bloomGo.go
}

// ------Core Methods of Counting Filter Below ---------=
func (cf *CountingFilter) Add(data []byte) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	hashes := baseHashes(data) // From ./bloomGo.go [using Murmur.go]
	for i := uint(0); i < cf.numHashes; i++ {
		idx := cf.location(hashes, i)
		// overflow warning: no more than 255
		if cf.counters[idx] < 255 {
			cf.counters[idx]++
		}
	}
}

// check if all required counter are > 0
func (cf *CountingFilter) Check(data []byte) bool {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	hashes := baseHashes(data)
	for i := uint(0); i < cf.numHashes; i++ {
		idx := cf.location(hashes, i)
		if cf.counters[idx] == 0 {
			return false //definitely not in the set
		}
	}
	return true // could be false positive
}

// decrements the counters fo rht given data
// true is removed successfully else false
func (cf *CountingFilter) Remove(data []byte) bool {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	// check if it is presented. Cannot remove non-existed data
	hashes := baseHashes(data)
	for i := uint(0); i < cf.numHashes; i++ {
		idx := cf.location(hashes, i)
		if cf.counters[idx] == 0 {
			// data is not presented
			return false // remove failed
		}
	}

	// it does exist
	for i := uint(0); i < cf.numHashes; i++ {
		idx := cf.location(hashes, i)
		cf.counters[idx]--
	}
	return true
}
