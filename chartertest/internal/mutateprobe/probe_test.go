package mutateprobe

import (
	"os"
	"testing"
)

func TestClampHoldsTheBounds(t *testing.T) {
	if Clamp(-5, 0, 10) != 0 || Clamp(15, 0, 10) != 10 || Clamp(5, 0, 10) != 5 {
		t.Fatal("Clamp does not hold the bounds")
	}
}

// TestLabelIsNamed passes whatever Label decides: the hollow scenario.
func TestLabelIsNamed(t *testing.T) {
	if Label(2) == "" {
		t.Fatal("Label is empty")
	}
}

// TestUnsteady fails when the mutate test asks it to, as a flaky test would.
func TestUnsteady(t *testing.T) {
	if os.Getenv("MUTATEPROBE_UNSTEADY") != "" {
		t.Fatal("unsteady")
	}
}

// TestLitter leaves a file where the mutate test asks it to.
func TestLitter(t *testing.T) {
	if p := os.Getenv("MUTATEPROBE_LITTER"); p != "" {
		if err := os.WriteFile(p, []byte("litter\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestClampAndLabel runs a subtest that notices a broken Clamp beside one
// that notices nothing.
func TestClampAndLabel(t *testing.T) {
	t.Run("clamp holds", func(t *testing.T) {
		if Clamp(5, 0, 10) != 5 {
			t.Fatal("Clamp moved a value inside the bounds")
		}
	})
	t.Run("label", func(t *testing.T) {
		if Label(2) == "" {
			t.Fatal("Label is empty")
		}
	})
}
