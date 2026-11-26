package toml

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// decodes decodes src and fails the test unless it decodes to want.
func decodes(t *testing.T, src string, want map[string]any) {
	t.Helper()
	got, err := Decode([]byte(src))
	if err != nil {
		t.Fatalf("Decode(%q): %v", src, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Decode(%q)\n got %#v\nwant %#v", src, got, want)
	}
}

// refused decodes src and fails the test unless it is refused with a
// *ParseError whose message contains each of words. It returns the error.
func refused(t *testing.T, src string, words ...string) *ParseError {
	t.Helper()
	got, err := Decode([]byte(src))
	if err == nil {
		t.Fatalf("Decode(%q) = %#v, want it refused", src, got)
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("Decode(%q): error %T, want *ParseError", src, err)
	}
	for _, w := range words {
		if !strings.Contains(pe.Msg, w) {
			t.Errorf("Decode(%q): message %q does not name %q", src, pe.Msg, w)
		}
	}
	return pe
}
