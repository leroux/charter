package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// decoder runs the command as toml-test does: the document on stdin, the
// tagged JSON on stdout.
func decoder(t *testing.T, doc string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = run(strings.NewReader(doc), &out, &errOut)
	return out.String(), errOut.String(), code
}

// tagged decodes the command's stdout, failing unless it exited 0.
func tagged(t *testing.T, doc string) map[string]any {
	t.Helper()
	out, errOut, code := decoder(t, doc)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out)
	}
	return v
}

func val(typ, v string) map[string]any { return map[string]any{"type": typ, "value": v} }

func TestAConfigFileDecodes(t *testing.T) {
	doc := `# The service's configuration
title = "TOML \"example\""
"quoted key" = 'C:\temp'
owner.name = "Tom"

[server]
port = 8080  # the default
debug = false

[server.limits]
max-conns = -1
`
	want := map[string]any{
		"title":      val("string", `TOML "example"`),
		"quoted key": val("string", `C:\temp`),
		"owner":      map[string]any{"name": val("string", "Tom")},
		"server": map[string]any{
			"port":   val("integer", "8080"),
			"debug":  val("bool", "false"),
			"limits": map[string]any{"max-conns": val("integer", "-1")},
		},
	}
	if got := tagged(t, doc); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestTypesAreKept(t *testing.T) {
	got := tagged(t, "a = \"42\"\nb = 42\n")
	if !reflect.DeepEqual(got["a"], val("string", "42")) || !reflect.DeepEqual(got["b"], val("integer", "42")) {
		t.Errorf("got %v, want a string 42 and an integer 42", got)
	}
}

func TestARefusalSaysWhere(t *testing.T) {
	out, errOut, code := decoder(t, "name = 'a'\nport = 1\nname = 'b'\n")
	if code == 0 || out != "" {
		t.Errorf("exit %d, stdout %q; want a failing exit and nothing on stdout", code, out)
	}
	if !strings.Contains(errOut, "line 3, column 1") {
		t.Errorf("stderr %q does not say line 3, column 1", errOut)
	}
}

func TestAnUnbuiltFormIsRefused(t *testing.T) {
	out, errOut, code := decoder(t, "pi = 3.14\n")
	if code == 0 || out != "" || !strings.Contains(errOut, "floats") {
		t.Errorf("exit %d, stdout %q, stderr %q; want a refusal naming floats", code, out, errOut)
	}
}

func TestValuesAreTagged(t *testing.T) {
	out, _, code := decoder(t, "s = 'x'\ni = -3\nb = true\n[t]\n")
	want := `{"b":{"type":"bool","value":"true"},"i":{"type":"integer","value":"-3"},"s":{"type":"string","value":"x"},"t":{}}`
	var v any
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil {
		t.Fatalf("exit %d, stdout %q", code, out)
	}
	if got, _ := json.Marshal(v); string(got) != want {
		t.Errorf("stdout %s, want %s", got, want)
	}
}

func TestARefusalExitsOne(t *testing.T) {
	out, errOut, code := decoder(t, "a = \n")
	if code != 1 || out != "" || !strings.HasPrefix(errOut, "toml: line 1") {
		t.Errorf("exit %d, stdout %q, stderr %q; want exit 1, the error on stderr only", code, out, errOut)
	}
}
