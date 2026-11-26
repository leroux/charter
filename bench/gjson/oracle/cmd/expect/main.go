// Command expect prints the reference answer to each question in the key:
// gjson at upstream 8e88233 (oracle/ref), run on the question's document
// and path. The grader compares the agent's answers with these.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

type question struct {
	ID   string `json:"id"`
	JSON string `json:"json"`
	Path string `json:"path"`
	Ask  string `json:"ask"` // "string": the result's String(); "exists": whether it exists
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	var qs []question
	if err := json.Unmarshal(data, &qs); err != nil {
		panic(err)
	}
	out := map[string]string{}
	for _, q := range qs {
		r := gjson.Get(q.JSON, q.Path)
		if q.Ask == "exists" {
			out[q.ID] = fmt.Sprint(r.Exists())
		} else {
			out[q.ID] = r.String()
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.Encode(out)
}
