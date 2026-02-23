package bloom

import (
	"fmt"
)

// Filter interface --> standard operations that every bloom filter variant must support
type Filter interface {
	// Add inserts data into the filter
	Add(data []byte)

	// Check --> true (false positive) if data is presented, false if it is DEFINITELY NOT
	Check(data []byte) bool
}

// Constant for out factory
const (
	TypeStandard = "standard"
	TypeScalable = "scalable"
	TypeCounting = "counting"
)

// Config --> what are needed for running a new Bloom Filter
type Config struct {
	Type          string  // type of filter
	ExpectedItems int     // Expected Number of items to be inserted
	FalsePosRate  float64 // Desired False Positive Rate(e.g., 0.01)
	GrowthFactor  float64 //Parameter only for 'Scalable Filter'
}

// Factory Design Pattern
// NewFilter is the factory. It returns the interface.
func NewFilter(cfg Config) (Filter, error) {
	switch cfg.Type {
	case TypeStandard:
		return NewStandardWithEstimatedParams(cfg.ExpectedItems, cfg.FalsePosRate), nil
	case TypeScalable:
		return NewScalableFilterWithEstimatedParams(cfg.ExpectedItems, cfg.FalsePosRate, cfg.GrowthFactor), nil
	case TypeCounting:
		return NewCountingWithEstimatedParams(cfg.ExpectedItems, cfg.FalsePosRate), nil
	default:
		return nil, fmt.Errorf("Unknown bloom filter type: %s", cfg.Type)
	}

}
