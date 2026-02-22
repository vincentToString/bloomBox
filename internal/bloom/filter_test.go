package bloom

import (
	"testing"
)

// tests that our NewFilter factory return the correct registered types
func TestFactory(t *testing.T) {

	// 1. TypeStandard
	stdCfg := Config{Type: TypeStandard, ExpectedItems: 100, FalsePosRate: 0.01}
	filter1, err := NewFilter(stdCfg)

	if err != nil {
		t.Fatalf("Fatory failed for standard filter: %v", err)
	}

	if _, ok := filter1.(*StandardFilter); !ok {
		t.Error("Factory did not return a *StandardFilter for TypeStandard")
	}

	// 2. TypeScalable
	scaleCfg := Config{Type: TypeScalable, ExpectedItems: 100, FalsePosRate: 0.01, GrowthFactor: 2.0}
	filter2, err := NewFilter(scaleCfg)

	if err != nil {
		t.Fatalf("Factory failed for scalable: %v", err)
	}
	if _, ok := filter2.(*ScalableFilter); !ok {
		t.Error("Factory did not return a *ScalableFilter for TypeScalable")
	}

	// 3. Test Unknown
	badCfg := Config{Type: "magic_bloom"}
	_, err = NewFilter(badCfg)
	if err == nil {
		t.Error("Expected an error when requesting an unknown filter type")
	}
}
