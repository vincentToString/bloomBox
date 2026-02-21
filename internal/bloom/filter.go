package bloom

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
}
