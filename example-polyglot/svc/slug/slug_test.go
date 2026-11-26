package slug

import (
	"testing"

	"example.com/slugs/claims"
)

func TestASlugIsLowerCaseWordsJoined(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_4NQQ, "lower-and-joined", "no-edge-separator")
	for title, want := range map[string]string{"Hello, World!": "hello-world", "  --Go--  ": "go"} {
		if got := Of(title); got != want {
			t.Errorf("Of(%q) = %q, want %q", title, got, want)
		}
	}
}
