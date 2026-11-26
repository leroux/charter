package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims/gate"
)

func TestWritesTheRegistryAsJSON(t *testing.T) {
	out := filepath.Join(t.TempDir(), "claims.json")
	var stderr bytes.Buffer
	if code := Claimsmap([]string{"-format", "json", "-o", out}, &stderr, self); code != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Charter string `json:"charter"`
		Source  string `json:"source"`
		Areas   []struct {
			File    string `json:"file"`
			Title   string `json:"title"`
			Summary string `json:"summary"`
		} `json:"areas"`
		Waivers []struct {
			Pattern string `json:"pattern"`
			Reason  string `json:"reason"`
		} `json:"waivers"`
		Entries []struct {
			Name      string   `json:"name"`
			Level     string   `json:"level"`
			Form      string   `json:"form"`
			Lens      string   `json:"lens"`
			File      string   `json:"file"`
			Supports  []string `json:"supports"`
			Accounts  []string `json:"accounts"`
			Serves    []string `json:"serves"`
			Realizes  []string `json:"realizes"`
			Scenarios []struct {
				Handle string `json:"handle"`
			} `json:"scenarios"`
			Refs []struct {
				File    string   `json:"file"`
				Func    string   `json:"func"`
				Line    int      `json:"line"`
				Handles []string `json:"handles"`
			} `json:"refs"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	reg, err := chartertest.Read(gate.Options)
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Entries) == 0 {
		t.Fatal("the registry read as empty")
	}
	if len(got.Entries) != len(reg.Entries) {
		t.Fatalf("%d entries in the JSON, %d in the registry", len(got.Entries), len(reg.Entries))
	}
	if got.Source == "" || got.Charter != "self" {
		t.Errorf("source %q, charter %q, want a source and self", got.Source, got.Charter)
	}
	if len(got.Areas) != len(reg.Areas) || len(got.Waivers) != len(reg.Waivers) {
		t.Errorf("%d areas and %d waivers, want %d and %d", len(got.Areas), len(got.Waivers), len(reg.Areas), len(reg.Waivers))
	}
	for i, a := range got.Areas {
		if a.Title == "" || a.File != reg.Areas[i].File || a.Summary != reg.Areas[i].Summary {
			t.Errorf("area %d = %+v, want %s with its title", i, a, reg.Areas[i].File)
		}
	}
	cited := false
	for i, e := range reg.Entries {
		g := got.Entries[i]
		if g.Name != e.Name || g.Level != e.Level || g.File != e.File || g.Form != e.Form || g.Lens != e.Lens {
			t.Errorf("entry %d: got %s/%s/%s/%s/%s, want %s/%s/%s/%s/%s", i, g.Name, g.Level, g.File, g.Form, g.Lens, e.Name, e.Level, e.File, e.Form, e.Lens)
		}
		if len(g.Supports) != len(e.Supports) || len(g.Refs) != len(e.Refs) || len(g.Scenarios) != len(e.Scenarios) {
			t.Errorf("%s: %d supports, %d refs and %d scenarios, want %d, %d and %d", e.Name, len(g.Supports), len(g.Refs), len(g.Scenarios), len(e.Supports), len(e.Refs), len(e.Scenarios))
		}
		if len(g.Serves) != len(e.Serves) || len(g.Realizes) != len(e.Realizes) {
			t.Errorf("%s: %d serves and %d realizes, want %d and %d", e.Name, len(g.Serves), len(g.Realizes), len(e.Serves), len(e.Realizes))
		}
		for j, r := range e.Refs {
			if j < len(g.Refs) && (g.Refs[j].Line != r.Line || len(g.Refs[j].Handles) != len(r.Handles)) {
				t.Errorf("%s ref %d: line %d with %d handles, want %d with %d", e.Name, j, g.Refs[j].Line, len(g.Refs[j].Handles), r.Line, len(r.Handles))
			}
		}
		if len(e.Refs) > 0 {
			cited = true
		}
	}
	if !cited {
		t.Error("no entry carries a citing test")
	}
	if null := regexp.MustCompile(`"[a-z]+":null`).FindString(strings.ReplaceAll(string(b), " ", "")); null != "" {
		t.Errorf("the JSON carries %s; empty arrays must be written as []", null)
	}
}

func TestExportsTheChartersOwnCharter(t *testing.T) {
	out := filepath.Join(t.TempDir(), "charter.json")
	var stderr bytes.Buffer
	if code := Claimsmap([]string{"-charter", "self", "-o", out}, &stderr, self); code != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Charter string `json:"charter"`
		Source  string `json:"source"`
		Entries []struct {
			ID string `json:"id"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if got.Charter != "self" || !strings.HasPrefix(got.Source, "claims") {
		t.Errorf("charter %q, source %q, want self from claims", got.Charter, got.Source)
	}
	if len(got.Entries) == 0 || !strings.HasPrefix(got.Entries[0].ID, "CH-") {
		t.Errorf("entries = %v, want the CH entries", got.Entries)
	}
	if code := Claimsmap([]string{"-charter", "other", "-o", out}, &stderr, self); code == 0 {
		t.Error("an unknown charter succeeded")
	}
}

func TestRejectsAnUnknownFormatAndWritesNothing(t *testing.T) {
	out := filepath.Join(t.TempDir(), "claims.json")
	var stderr bytes.Buffer
	if code := Claimsmap([]string{"-format", "html", "-o", out}, &stderr, self); code == 0 {
		t.Fatal("an unknown format succeeded")
	}
	if !strings.Contains(stderr.String(), `"html"`) {
		t.Errorf("stderr = %q, want it to name the format", stderr.String())
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("the output file was written")
	}
}

// self is the charter's own charter, the one the tests export.
var self = Charter{Name: "self", Dir: "claims", Options: gate.Options}
