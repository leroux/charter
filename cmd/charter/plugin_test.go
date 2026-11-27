package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

const repo = "../.."

func readJSON(t *testing.T, path string, v any) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(repo, path))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}

// skillFiles are the plugin's skills by directory name.
func skillFiles(t *testing.T, dir string) map[string]string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(repo, dir, "skills", "*", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		out[filepath.Base(filepath.Dir(p))] = string(b)
	}
	return out
}

var frontName = regexp.MustCompile(`(?m)\A---\n(?:.*\n)*?name: (\S+)\n`)

func TestThePluginShipsTheSixSkills(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_5JMD, "plugin-lists-the-skills")
	var market struct {
		Name    string
		Plugins []struct{ Name, Source string }
	}
	readJSON(t, ".claude-plugin/marketplace.json", &market)
	if len(market.Plugins) != 1 || market.Plugins[0].Name != "charter" {
		t.Fatalf("the marketplace names %+v, want the one plugin charter", market.Plugins)
	}
	dir := strings.TrimPrefix(market.Plugins[0].Source, "./")
	var plugin struct{ Name string }
	readJSON(t, filepath.Join(dir, ".claude-plugin", "plugin.json"), &plugin)
	if plugin.Name != "charter" {
		t.Errorf("the plugin's manifest names %q", plugin.Name)
	}
	skills := skillFiles(t, dir)
	var names []string
	for n, src := range skills {
		names = append(names, n)
		if m := frontName.FindStringSubmatch(src); m == nil || m[1] != n {
			t.Errorf("%s/SKILL.md names itself %v, not its directory", n, m)
		}
	}
	sort.Strings(names)
	if want := "claims-audit claims-draft claims-implement claims-narrate claims-probe claims-review"; strings.Join(names, " ") != want {
		t.Errorf("the plugin holds %v, want %s", names, want)
	}
}

func TestNoSkillNamesOneRepositorysLayout(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_5JMD, "no-fixed-layout")
	for n, src := range skillFiles(t, "plugin") {
		for _, fixed := range []string{"cmd/claimsgen", "claimsgen", "claims/README.md", "go run ./cmd/charter", "go run ./charter", "superpowers:", "torii", "ratelimit", ".github/"} {
			if strings.Contains(src, fixed) {
				t.Errorf("%s names %q, a part of one repository's layout", n, fixed)
			}
		}
		if !strings.Contains(src, "charter.toml") {
			t.Errorf("%s does not say the layout is read from charter.toml", n)
		}
	}
}

func TestTheReadmeGivesTheClaudeMdSection(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_5JMD, "claude-md-template")
	b, err := os.ReadFile(filepath.Join(repo, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)
	i := strings.Index(src, "```markdown\n## Charter\n")
	if i < 0 {
		t.Fatal("the README holds no CLAUDE.md section as a markdown template")
	}
	block := src[i:]
	block = block[:strings.Index(block[3:], "```")+3]
	for _, want := range []string{"claims/charter.toml", "charter check", "charter gen", "charter next"} {
		if !strings.Contains(block, want) {
			t.Errorf("the CLAUDE.md template does not name %q", want)
		}
	}
}
