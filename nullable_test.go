package maybe_test

import (
	"encoding/json"
	"testing"

	"github.com/iamBelugaa/maybe"
)

// TestNullableJSON verifies that Nullable[T] properly marshals to and from JSON.
// It checks both null and non-null cases.
func TestNullableJSON(t *testing.T) {
	t.Run("Marshal/Unmarshal", func(t *testing.T) {
		// Marshal a valid Nullable[int] -> should produce "42"
		n := maybe.NullableOf(42)
		data, err := json.Marshal(n)
		if err != nil || string(data) != "42" {
			t.Errorf("Marshal valid failed: got %s, err: %v", data, err)
		}

		// Marshal a null Nullable[int] -> should produce "null"
		null := maybe.Null[int]()
		data, err = json.Marshal(null)
		if err != nil || string(data) != "null" {
			t.Errorf("Marshal null failed: got %s, err: %v", data, err)
		}

		// Unmarshal "3.14" into Nullable[float64]
		var unm maybe.Nullable[float64]
		err = json.Unmarshal([]byte("3.14"), &unm)
		val, ok := unm.Extract()
		if err != nil || !unm.IsValid() || !ok || val != 3.14 {
			t.Errorf("Unmarshal float64 failed: got %v, ok: %v, err: %v", val, ok, err)
		}
	})
}

// TestNullableConversion checks conversion from Nullable to Option type.
func TestNullableConversion(t *testing.T) {
	t.Run("ToOption", func(t *testing.T) {
		// Converting a null Nullable[string] should yield None
		null := maybe.Null[string]()
		opt1 := null.ToOption()
		if opt1.IsSome() {
			t.Error("Null -> Option should result in None")
		}

		// Converting a valid Nullable[int] should yield Some with same value
		valid := maybe.NullableOf(42)
		opt2 := valid.ToOption()
		if !opt2.IsSome() || opt2.Unwrap() != 42 {
			t.Error("Valid Nullable -> Option conversion failed")
		}
	})
}
