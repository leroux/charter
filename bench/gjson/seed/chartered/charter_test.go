package gjson

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/tidwall/gjson/claims"
)

// readmeJSON is the document the README's path examples read.
const readmeJSON = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44},
    {"first": "Roger", "last": "Craig", "age": 68},
    {"first": "Jane", "last": "Murphy", "age": 47}
  ]
}`

// expect fails the test when path on json does not return want as its
// String.
func expect(t *testing.T, json, path, want string) {
	t.Helper()
	if got := Get(json, path).String(); got != want {
		t.Errorf("Get(%s): got %q, want %q", path, got, want)
	}
}

// absent fails the test when path on json returns a result that exists.
func absent(t *testing.T, json, path string) {
	t.Helper()
	if r := Get(json, path); r.Exists() {
		t.Errorf("Get(%s): got %q, want a result that does not exist", path, r.Raw)
	}
}

// matches returns the first names a query on the friends selects.
func matches(query string) string {
	return Get(readmeJSON, "friends.#["+query+"]#.first").String()
}

func TestNeedReadmePaths(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_ZQJR, "readme-paths")
	for path, want := range map[string]string{
		"name.last":                        "Anderson",
		"age":                              "37",
		"children":                         `["Sara","Alex","Jack"]`,
		"children.#":                       "3",
		"children.1":                       "Alex",
		"child*.2":                         "Jack",
		"c?ildren.0":                       "Sara",
		`fav\.movie`:                       "Deer Hunter",
		"friends.#.first":                  `["Dale","Roger","Jane"]`,
		"friends.1.last":                   "Craig",
		`friends.#[last=="Murphy"].first`:  "Dale",
		`friends.#[last=="Murphy"]#.first`: `["Dale","Jane"]`,
		`friends.#[age>45]#.last`:          `["Craig","Murphy"]`,
		`friends.#[first%"D*"].last`:       "Murphy",
		`friends.#[first!%"D*"].last`:      "Craig",
		"children|@reverse":                `["Jack","Alex","Sara"]`,
		"children|@reverse|0":              "Jack",
	} {
		expect(t, readmeJSON, path, want)
	}
}

func TestNeedMissingIsAbsent(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_ZQJR, "missing-is-absent")
	for _, path := range []string{"nickname", "name.middle", "children.3", "friends.9.first", "age.years"} {
		absent(t, readmeJSON, path)
	}
}

// noPanic calls f and fails the test, naming the input, if it panics.
func noPanic(t *testing.T, json, path string, f func()) {
	t.Helper()
	defer func() {
		if v := recover(); v != nil {
			t.Fatalf("panic %v on document %q, path %q", v, json, path)
		}
	}()
	f()
}

// pathForms are paths of every form the syntax has.
var pathForms = []string{
	"", "a", "a.b", "a.0", "a.#", "a.#.b", "a.#[b=1]", "a.#[b>1]#.c", `a.#[b%"x*"]`,
	"a*", "a?.b", `a\.b`, "a|b", "@reverse", "@ugly|a", `@pretty:{"indent":"  "}`,
	"..#", "..0.a", "#", "0", "a.", ".", "|", "@", "#[", "#[a", `#[a="`, "\\",
}

func TestNeedRandomInputNeverPanics(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_8FW8, "random-input")
	chartertest.Verifies(t, claims.SW_261126_9CTW, "no-panic")
	rnd := rand.New(rand.NewSource(1))
	alphabet := []byte(`{}[]":,.#*?\|@ 0123456789abtrue-falsnul`)
	b := make([]byte, 64)
	for i := 0; i < 20000; i++ {
		n := rnd.Intn(len(b))
		for j := 0; j < n; j++ {
			b[j] = alphabet[rnd.Intn(len(alphabet))]
		}
		doc := string(b[:n])
		path := pathForms[rnd.Intn(len(pathForms))]
		noPanic(t, doc, path, func() { Get(doc, path) })
	}
}

func TestNeedTruncatedNeverPanics(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_8FW8, "truncated")
	for n := 0; n <= len(readmeJSON); n++ {
		doc := readmeJSON[:n]
		for _, path := range append(pathForms, "name.last", "friends.#.first", `friends.#[age>45]#.last`, "children|@reverse|0") {
			noPanic(t, doc, path, func() { Get(doc, path) })
		}
	}
}

func TestNeedGoValues(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_TFYB, "as-go-values")
	doc := `{"s":"hi","n":2.5,"b":true,"a":[1,2],"o":{"k":"v"}}`
	if Get(doc, "s").String() != "hi" || Get(doc, "n").Float() != 2.5 || !Get(doc, "b").Bool() {
		t.Error("the string, number or boolean did not read back as its Go value")
	}
	if a := Get(doc, "a").Array(); len(a) != 2 || a[1].Int() != 2 {
		t.Errorf("the array read back as %v", a)
	}
	if m := Get(doc, "o").Map(); len(m) != 1 || m["k"].String() != "v" {
		t.Errorf("the object read back as %v", m)
	}
}

func TestNeedBigIntegerDigits(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_TFYB, "big-integer-digits")
	r := Get(`{"id":9007199254740993}`, "id")
	if r.Int() != 9007199254740993 || r.String() != "9007199254740993" {
		t.Errorf("got %d and %q", r.Int(), r.String())
	}
}

func TestNeedCheckFirst(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_94D2, "check-first")
	if !Valid(`{"a":[1,2]}`) || Valid(`{"a":[1,2],}`) {
		t.Error("Valid did not tell the well-formed document from the one with a trailing comma")
	}
}

func TestPathMembersAndIndices(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_T9SN, "member", "index", "missing", "empty-path")
	expect(t, `{"name":{"first":"Tom","last":"Anderson"}}`, "name.last", "Anderson")
	expect(t, `{"children":["Sara","Alex","Jack"]}`, "children.1", "Alex")
	absent(t, `{"a":{"b":1}}`, "a.c")
	absent(t, `{"a":[1,2,3]}`, "a.3")
	absent(t, `{"a":"text"}`, "a.b")
	absent(t, `{"a":1}`, "")
}

func TestPathWildcards(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_4HQP, "star", "question", "first-match", "no-match")
	doc := `{"children":["Sara","Alex","Jack"]}`
	expect(t, doc, "child*.2", "Jack")
	expect(t, doc, "c?ildren.0", "Sara")
	r := Get(`{"abc":1,"abd":2}`, "ab*")
	if r.Type != Number || r.Int() != 1 {
		t.Errorf("ab*: got %s %q, want the number 1", r.Type, r.Raw)
	}
	absent(t, `{"abc":1,"abd":2}`, "x*")
}

func TestPathEscapes(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_FX9X, "escaped-dot", "escaped-wildcard")
	expect(t, `{"fav.movie":"Deer Hunter"}`, `fav\.movie`, "Deer Hunter")
	expect(t, `{"keyXv":"x","key*v":"star"}`, `key\*v`, "star")
	absent(t, `{"keyXv":"x"}`, `key\*v`)
	expect(t, `{"keyvX":"x","keyv?":"q"}`, `keyv\?`, "q")
}

func TestPathKeysExact(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_8RR3, "case-sensitive", "escaped-key", "duplicate-first")
	absent(t, `{"Name":1}`, "name")
	expect(t, `{"Name":1}`, "Name", "1")
	expect(t, `{"a\"b":"q"}`, `a"b`, "q")
	doc := `{"name":"Alex","name":"Peter"}`
	expect(t, doc, "name", "Alex")
	if got := Parse(doc).Map()["name"].String(); got != "Alex" {
		t.Errorf("Map holds %q under name, want Alex", got)
	}
}

func TestPathPipes(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_SQKQ, "pipe", "pipe-after-modifier", "chaining-disabled")
	expect(t, readmeJSON, "friends|0|first", "Dale")
	expect(t, readmeJSON, "friends|@reverse|0|age", "47")
	DisableChaining = true
	defer func() { DisableChaining = false }()
	expect(t, `{"a|b":"pipe","a":{"b":"dot"}}`, "a|b", "pipe")
}

func TestPathJSONLines(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_ZQDQ, "lines-array", "for-each-line")
	lines := "true\nfalse\n{\"name\":\"tom\"}\n[1,2]\n{\"name\":\"janet\"}\nnull\n12.5\n"
	expect(t, lines, "..#", "7")
	expect(t, lines, "..0", "true")
	expect(t, lines, "..2.name", "tom")
	expect(t, lines, "..#.name", `["tom","janet"]`)
	var seen []string
	ForEachLine(lines, func(line Result) bool {
		seen = append(seen, line.Raw)
		return len(seen) < 3
	})
	if strings.Join(seen, " ") != `true false {"name":"tom"}` {
		t.Errorf("ForEachLine visited %q", seen)
	}
}

func TestArrayCount(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_6FGV, "count", "count-empty", "count-not-array")
	r := Get(`{"children":["Sara","Alex","Jack"]}`, "children.#")
	if r.Type != Number || r.Int() != 3 {
		t.Errorf("children.#: got %s %v, want the number 3", r.Type, r.Num)
	}
	expect(t, `[]`, "#", "0")
	absent(t, `{"a":{"x":1}}`, "a.#")
}

func TestArrayEachElement(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_V9J9, "each", "each-skips-missing", "each-none")
	expect(t, readmeJSON, "friends.#.first", `["Dale","Roger","Jane"]`)
	doc := `{"a":[{"b":1},{"c":2},{"b":3}]}`
	expect(t, doc, "a.#.b", "[1,3]")
	if r := Get(doc, "a.#.z"); !r.Exists() || r.Raw != "[]" {
		t.Errorf("a.#.z: got %q, want [] that exists", r.Raw)
	}
}

func TestArrayQueries(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2TYJ, "first", "all", "no-match")
	expect(t, readmeJSON, `friends.#[last=="Murphy"].first`, "Dale")
	expect(t, readmeJSON, `friends.#[last=="Murphy"]#.first`, `["Dale","Jane"]`)
	absent(t, readmeJSON, `friends.#[age>100]`)
	absent(t, readmeJSON, `friends.#[age>100]#`)
}

func TestArrayQueryOperators(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_XMB5, "strings", "numbers", "like", "booleans")
	for query, want := range map[string]string{
		`first>"Dale"`:  `["Roger","Jane"]`,
		`first=="Dale"`: `["Dale"]`,
		`first="Dale"`:  `["Dale"]`,
		`age>45`:        `["Roger","Jane"]`,
		`age<=44`:       `["Dale"]`,
		`age!=44`:       `["Roger","Jane"]`,
		`first%"D*"`:    `["Dale"]`,
		`first!%"D*"`:   `["Roger","Jane"]`,
	} {
		if got := matches(query); got != want {
			t.Errorf("#[%s]#: got %s, want %s", query, got, want)
		}
	}
	doc := `{"a":[{"n":1,"ok":true},{"n":2,"ok":false},{"n":3,"ok":true}]}`
	expect(t, doc, `a.#[ok=true]#.n`, "[1,3]")
	expect(t, doc, `a.#[ok!=true]#.n`, "[2]")
}

func TestModifiersBuiltIn(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_H629, "reverse", "ugly", "pretty")
	expect(t, readmeJSON, "children|@reverse", `["Jack","Alex","Sara"]`)
	expect(t, `{"a":"x"}`, "a|@reverse", "x")
	expect(t, "{\n  \"a\": [1, 2],\n  \"b\": {\"c\": true}\n}", "@ugly", `{"a":[1,2],"b":{"c":true}}`)
	pretty := Get(`{"b":1,"a":2}`, `@pretty:{"sortKeys":true}`).String()
	if pretty != "{\n  \"a\": 2,\n  \"b\": 1\n}\n" {
		t.Errorf("@pretty sorted: got %q", pretty)
	}
	expect(t, pretty, "@ugly", `{"a":2,"b":1}`)
}

func TestModifiersCustom(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_4Z6X, "custom", "exists", "unknown")
	upper := func(json, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(json)
		}
		return json
	}
	AddModifier("case", upper)
	expect(t, readmeJSON, "children|@case:upper", `["SARA","ALEX","JACK"]`)
	if !ModifierExists("case", upper) || ModifierExists("nope", upper) {
		t.Error("ModifierExists did not tell a registered name from an unregistered one")
	}
	absent(t, readmeJSON, "@nope")
}

func TestModifiersDisabled(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_Z3CC, "disabled")
	DisableModifiers = true
	defer func() { DisableModifiers = false }()
	expect(t, `{"@ugly":1}`, "@ugly", "1")
}

func TestResultString(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_PKGR, "string-decoded", "integer-as-written", "float-shortest", "literals-and-json")
	expect(t, `{"s":"  café\n "}`, "s", "  café\n ")
	expect(t, `{"n":9007199254740993}`, "n", "9007199254740993")
	if r := Get(`{"n":1.50}`, "n"); r.String() != "1.5" || r.Raw != "1.50" {
		t.Errorf("1.50: got String %q, Raw %q", r.String(), r.Raw)
	}
	doc := `{"t":true,"f":false,"z":null,"o":{"a": 1}}`
	expect(t, doc, "t", "true")
	expect(t, doc, "f", "false")
	expect(t, doc, "z", "")
	expect(t, doc, "o", `{"a": 1}`)
}

func TestResultIntegers(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_5VFD, "int64-exact", "uint64-exact", "from-string")
	if n := Get(`[-9223372036854775808]`, "0").Int(); n != -9223372036854775808 {
		t.Errorf("Int: got %d", n)
	}
	if n := Get(`[18446744073709551615]`, "0").Uint(); n != 18446744073709551615 {
		t.Errorf("Uint: got %d", n)
	}
	s := Get(`["12"]`, "0")
	if s.Int() != 12 || s.Float() != 12 || Parse("true").Float() != 1 {
		t.Errorf("a numeric string read as %d and %v", s.Int(), s.Float())
	}
}

func TestResultBool(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_ZGV3, "literals", "numbers", "strings")
	for raw, want := range map[string]bool{
		"true": true, "false": false, "null": false,
		"1": true, "-2": true, "0": false,
		`"yes"`: true, `"true"`: true, `""`: false, `"0"`: false, `"false"`: false,
	} {
		if got := Parse(raw).Bool(); got != want {
			t.Errorf("Bool of %s: got %v, want %v", raw, got, want)
		}
	}
}

func TestResultExists(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_26XP, "null-exists", "missing-not")
	doc := `{"a":null}`
	if r := Get(doc, "a"); !r.Exists() || r.Type != Null {
		t.Errorf("a: exists %v, type %s", r.Exists(), r.Type)
	}
	if r := Get(doc, "b"); r.Exists() || r.Type != Null {
		t.Errorf("b: exists %v, type %s", r.Exists(), r.Type)
	}
}

func TestResultWalks(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_NKTP, "array-of-scalar", "map", "for-each-order")
	if a := Parse("1").Array(); len(a) != 1 || a[0].Int() != 1 {
		t.Errorf("Array of 1: %v", a)
	}
	if a := Get(`{}`, "x").Array(); len(a) != 0 {
		t.Errorf("Array of nothing: %v", a)
	}
	m := Parse(`{"a":1,"b":"x"}`).Map()
	if len(m) != 2 || m["a"].Int() != 1 || m["b"].String() != "x" {
		t.Errorf("Map: %v", m)
	}
	if len(Parse(`"s"`).Map()) != 0 {
		t.Error("Map of a string is not empty")
	}
	var keys []string
	Parse(`{"z":1,"a":2,"m":3}`).ForEach(func(k, _ Result) bool {
		keys = append(keys, k.String())
		return len(keys) < 2
	})
	if strings.Join(keys, ",") != "z,a" {
		t.Errorf("ForEach visited %v", keys)
	}
}

func TestResultValue(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_NRBX, "value-types")
	v, ok := Parse(`{"b":true,"n":1.5,"s":"x","z":null,"a":[1]}`).Value().(map[string]interface{})
	if !ok {
		t.Fatal("Value of an object is not a map")
	}
	if v["b"] != true || v["n"] != 1.5 || v["s"] != "x" || v["z"] != nil {
		t.Errorf("Value: %v", v)
	}
	if a, ok := v["a"].([]interface{}); !ok || len(a) != 1 || a[0] != 1.0 {
		t.Errorf("Value of the array: %#v", v["a"])
	}
}

func TestResultParse(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_Z23S, "parse-value", "result-get")
	if Parse(`  {"a":1} `).Type != JSON || Parse("12").Type != Number || Parse(`"x"`).Type != String {
		t.Error("Parse did not type the values")
	}
	if got := Parse(`{"name":{"last":"Anderson"}}`).Get("name").Get("last").String(); got != "Anderson" {
		t.Errorf("Result.Get: got %q", got)
	}
}

func TestResultLess(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_4WX9, "by-type", "by-value")
	null, f, tr := Parse("null"), Parse("false"), Parse("true")
	if !null.Less(f, true) || !f.Less(tr, true) || tr.Less(f, true) {
		t.Error("types out of order")
	}
	upper, lower := Parse(`"ABC"`), Parse(`"abc"`)
	if !upper.Less(lower, true) || upper.Less(lower, false) {
		t.Error("case order wrong")
	}
	if !Parse("123").Less(Parse("456"), true) {
		t.Error("123 does not sort before 456")
	}
}

func TestValidity(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_KNNZ, "valid", "invalid")
	for _, doc := range []string{"0", `{"a":[1,2]}`, `"x"`, " true "} {
		if !Valid(doc) || !ValidBytes([]byte(doc)) {
			t.Errorf("%q is not valid", doc)
		}
	}
	for _, doc := range []string{"00", `{"a":1,}`, `"x`, "1 2", ""} {
		if Valid(doc) {
			t.Errorf("%q is valid", doc)
		}
	}
}

func TestGetIsLenient(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_9CTW, "lenient")
	expect(t, `{"array":["a","b",0],}`, "array.1", "b")
}

func TestGetBytesAndMany(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_3FK5, "bytes-same", "many-same")
	paths := []string{"name.last", "age", "children", "friends.#.first", `friends.#[age>45]#.last`, "nope"}
	for _, p := range paths {
		b := []byte(readmeJSON)
		r := GetBytes(b, p)
		s := r.String()
		for i := range b {
			b[i] = 'x'
		}
		if want := Get(readmeJSON, p); r.Type != want.Type || s != want.String() || r.String() != s {
			t.Errorf("GetBytes(%s): got %q, want %q", p, r.String(), want.String())
		}
	}
	many := GetMany(readmeJSON, paths[:3]...)
	for i, r := range many {
		if r.String() != Get(readmeJSON, paths[i]).String() {
			t.Errorf("GetMany[%d]: got %q", i, r.String())
		}
	}
	if len(many) != 3 {
		t.Errorf("GetMany returned %d results", len(many))
	}
}
