package bloom

import (
	"sync"
)

// Scalable filter grows dynamically. It keep a list of standard filters and add new one when full.
type ScalableFilter struct {
	filters         []*StandardFilter
	mu              sync.RWMutex
	currentCapacity int // capacity of current (newest) active filter
	falsePosRate    float64
	growthFactor    float64 // float to support 1.5x, 1.6x growth
	itemsAdded      int     // items added into the current active filter
}

// Constructor for Bloom Scalable Filter
func NewScalableFilterWithEstimatedParams(initialCapacity int, fpr float64, growthFactor float64) *ScalableFilter {
	if initialCapacity <= 0 {
		initialCapacity = 1000 // default
	}

	if growthFactor <= 1.0 {
		growthFactor = 2.0 //default
	}

	firstFilter := NewStandardWithEstimatedParams(initialCapacity, fpr)

	return &ScalableFilter{
		filters:         []*StandardFilter{firstFilter},
		currentCapacity: initialCapacity,
		falsePosRate:    fpr,
		growthFactor:    growthFactor,
		itemsAdded:      0,
	}
}

// Add Data to the Scalable Filter. If "full", add new one
func (sf *ScalableFilter) Add(data []byte) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	// Check if current (the latest) filter is full
	if sf.itemsAdded >= sf.currentCapacity {
		// Time to scale

		// 1. Tighten the FP rate for the new Filter
		/*
			In a sequance of filter, FP(total) = 1 - (1 - p)^k
			If Small FP --> FP(total) = k * p
			If 'p' remains unchange --> each filter always has the same rate
			--> FP(total) is k times worse
		*/
		sf.falsePosRate *= 0.9

		// 2. Calculate new capacity and update it
		newCapacity := int(float64(sf.currentCapacity) * sf.growthFactor)
		sf.currentCapacity = newCapacity
		sf.itemsAdded = 0 //reset

		// 3. Create and append this new filter
		sf.filters = append(sf.filters, NewStandardWithEstimatedParams(sf.currentCapacity, sf.falsePosRate))

	}

	// Final step --> always add to the latest created filter
	sf.filters[len(sf.filters)-1].Add(data)
	sf.itemsAdded++

}

// Check searched across all internal standard filters
// Search sequantially (Reverse Order), anytime a filter return "True", we return True
func (sf *ScalableFilter) Check(data []byte) bool {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	//Newer Created Filters first. (More items and newer items)
	for i := len(sf.filters) - 1; i >= 0; i-- {
		// Just use standard filter's own check
		if sf.filters[i].Check(data) {
			return true
		}
	}
	return false
}
