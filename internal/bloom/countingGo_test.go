package bloom

import (
	"testing"
)

func TestCountingFilterBasic(t *testing.T) {
	// 100 items, 1% false positive rate

	cf := NewCountingWithEstimatedParams(100, 0.01)

	// 1. Test Add and Check
	cf.Add([]byte("apple"))
	cf.Add([]byte("orange"))

	if !cf.Check([]byte("apple")) {
		t.Error("Expected to find apple")
	}
	if !cf.Check([]byte("orange")) {
		t.Error("Expected to find orange")
	}
	if cf.Check([]byte("watermelon")) {
		t.Error("Did not expect to find watermelon (false positive is unlikely now)")
	}

	// 2. Test remove
	success := cf.Remove([]byte("apple"))
	if !success {
		t.Error("Expected Removal of 'apple' to return True")
	}
	// Recheck to confirm apple is not there
	if cf.Check([]byte("apple")) {
		t.Error("Apple should be removed")
	}

	if !cf.Check([]byte("orange")) {
		t.Error("Orange should still be in the filter")
	}

	// 3. Test removing something not presented
	success = cf.Remove([]byte("banana"))
	if success {
		t.Error("Expected removal of 'banana' to return false (it does not exist)")
	}

}
