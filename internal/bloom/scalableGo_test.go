package bloom

import (
	"fmt"
	"testing"
)

// test basic Add and Check
func TestScalableFilterBasic(t *testing.T) {
	sf := NewScalableFilterWithEstimatedParams(100, 0.01, 2.0)

	sf.Add([]byte("apple"))
	sf.Add([]byte("orange"))
	// both apple and orange should returned true(although Bloom Filter's true will not be reliable
	// due to false positive)

	if !sf.Check([]byte("apple")) {
		t.Error("Expected apple to be found")
	}

	if !sf.Check([]byte("orange")) {
		t.Error("Expected orange to be found")
	}

	if sf.Check([]byte("watermelon")) {
		t.Error("Did nt expect watermelon to be found (False positive very unlikely here as we have 100 c and 2 i)")
	}
}

// test scaling function when capacity is reached
func TestScalableFilterScaling(t *testing.T) {
	initialSize := 10
	sf := NewScalableFilterWithEstimatedParams(initialSize, 0.01, 2.0)

	if len(sf.filters) != 1 {
		t.Fatalf("Expected 1 internal filter, got %d", len(sf.filters))
	}

	// Add initialSize + 1 items to trigger scaling one time
	for i := 0; i <= initialSize; i++ {
		sf.Add([]byte(fmt.Sprintf("item_%d", i)))
	}

	// we should see 2 filters now
	if len(sf.filters) != 2 {
		t.Fatalf("Expected 2 internal filter, got %d", len(sf.filters))
	}

	// Verity all inserted items
	for i := 0; i <= initialSize; i++ {
		if !sf.Check([]byte(fmt.Sprintf("item_%d", i))) {
			t.Errorf("Could not find the inserted item: %d", i)
		}

	}
}
