// Command toml-decoder speaks toml-test's decoder protocol: it reads a TOML
// document on stdin and writes it on stdout as tagged JSON, every value an
// object {"type": ..., "value": ...}. A document it refuses is reported on
// stderr with a non-zero exit, and nothing is written on stdout.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"example.com/toml"
)

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr))
}

// run is the command: it returns the exit code.
func run(in io.Reader, out, errOut io.Writer) int {
	data, err := io.ReadAll(in)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	doc, err := toml.Decode(data)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tag(doc)); err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	return 0
}

// tag turns a decoded value into toml-test's tagged JSON.
func tag(v any) any {
	switch v := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, e := range v {
			out[k] = tag(e)
		}
		return out
	case string:
		return map[string]string{"type": "string", "value": v}
	case int64:
		return map[string]string{"type": "integer", "value": strconv.FormatInt(v, 10)}
	case bool:
		return map[string]string{"type": "bool", "value": strconv.FormatBool(v)}
	default:
		panic(fmt.Sprintf("toml-decoder: no tag for %T", v))
	}
}
