<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: go run ./cmd/claimsgen -->

# The charter, read through

What the system does, without its code: 4 needs, what a person or the operator can rely on, and 24 decisions in 5 areas, how the system behaves at each boundary. Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.

1. [What people need](#what-people-need)
2. [The areas](#the-areas)
3. [The decisions, area by area](#the-decisions-area-by-area)

## What people need

Each need is listed under the area that does most to meet it.

**Paths**

- [STK-261126-ZQJR](needs.md#stk-261126-zqjr-a-path-returns-the-value-the-documented-syntax-names-and-nothing-when-it-names-none) A path returns the value the documented syntax names, and nothing when it names none

**Results**

- [STK-261126-TFYB](needs.md#stk-261126-tfyb-a-caller-reads-a-result-as-the-go-value-it-needs-and-a-big-integer-keeps-every-digit) A caller reads a result as the Go value it needs, and a big integer keeps every digit

**Validity**

- [STK-261126-8FW8](needs.md#stk-261126-8fw8-get-never-panics-whatever-the-document-or-the-path) Get never panics, whatever the document or the path
- [STK-261126-94D2](needs.md#stk-261126-94d2-a-caller-can-tell-whether-a-document-is-valid-json-before-trusting-it) A caller can tell whether a document is valid JSON before trusting it

## The areas

| # | Area | What it holds | Decisions |
|---|---|---|---|
| 1 | [Paths](paths.md) | How Get reads a path: the keys and indices it is made of, the wildcards and escapes a key may hold, the pipe that chains one path onto the result of another, and the `..` prefix that reads JSON lines. | 6 |
| 2 | [Arrays](arrays.md) | What the `#` component does in a path: the length of an array, a path read from every element, and the queries that pick the elements that match. | 4 |
| 3 | [Modifiers](modifiers.md) | The `@` components of a path: the built-in modifiers, the ones a caller adds, and the switch that turns them off. | 3 |
| 4 | [Results](results.md) | What a caller reads off a Result: its type and raw text, and the conversions to Go values, the walks over arrays and objects, and the order Less puts results in. | 8 |
| 5 | [Validity](validity.md) | What gjson does with a document it cannot trust: Valid to ask, a reader that does not validate and never panics, and the byte and many-path forms of Get that answer as Get does. | 3 |

Which areas serve each need, with how many of their decisions:

| Need | Areas |
|---|---|
| [STK-261126-ZQJR](needs.md#stk-261126-zqjr-a-path-returns-the-value-the-documented-syntax-names-and-nothing-when-it-names-none) | paths 6 · arrays 4 · modifiers 3 · results 2 · validity 1 |
| [STK-261126-8FW8](needs.md#stk-261126-8fw8-get-never-panics-whatever-the-document-or-the-path) | validity 1 |
| [STK-261126-TFYB](needs.md#stk-261126-tfyb-a-caller-reads-a-result-as-the-go-value-it-needs-and-a-big-integer-keeps-every-digit) | results 6 |
| [STK-261126-94D2](needs.md#stk-261126-94d2-a-caller-can-tell-whether-a-document-is-valid-json-before-trusting-it) | validity 2 |

## The decisions, area by area

### 1. Paths

How Get reads a path: the keys and indices it is made of, the wildcards and escapes a key may hold, the pipe that chains one path onto the result of another, and the `..` prefix that reads JSON lines.

- **[SW-261126-T9SN](paths.md#sw-261126-t9sn-a-path-is-keys-separated-by-dots-a-key-naming-an-object-member-and-a-number-an-array-index)** A path is keys separated by dots, a key naming an object member and a number an array index  
  Why: dot paths are the whole interface; a caller predicts the answer from the document alone, and a missing value is told apart from a present one by Exists, not by a guess.
- **[SW-261126-4HQP](paths.md#sw-261126-4hqp-a-key-with--or--matches-the-first-member-in-document-order-whose-key-fits-the-pattern)** A key with `*` or `?` matches the first member, in document order, whose key fits the pattern  
  Why: a wildcard lets a caller reach a value whose key it knows only in part; one value, the first, keeps a wildcard path returning the same kind of result as any other path.
- **[SW-261126-FX9X](paths.md#sw-261126-fx9x-a-backslash-makes-the-next-character-of-a-path-a-literal-part-of-the-key)** A backslash makes the next character of a path a literal part of the key  
  Why: JSON keys may hold dots and wildcard characters, and a caller must be able to name every key a document can have.
- **[SW-261126-8RR3](paths.md#sw-261126-8rr3-a-key-matches-exactly-case-and-all-and-a-key-the-document-repeats-gives-its-first-value)** A key matches exactly, case and all, and a key the document repeats gives its first value *(invariant)*  
  Why: JSON keys are case-sensitive and a document may hold both `id` and `ID`; folding case would make one unreachable.
- **[SW-261126-SQKQ](paths.md#sw-261126-sqkq-a-pipe-applies-the-path-after-it-to-the-result-of-the-path-before-it)** A pipe applies the path after it to the result of the path before it  
  Why: a pipe lets a caller run a path on the output of a modifier or of an array path, which a dot cannot express, since a dot after `#.first` reads into each element.
- **[SW-261126-ZQDQ](paths.md#sw-261126-zqdq-a-path-starting-with--reads-the-document-as-json-lines-one-value-a-line-as-if-they-were-an-array)** A path starting with `..` reads the document as JSON lines, one value a line, as if they were an array  
  Why: logs and streams are JSON lines; reading them with the array syntax a caller already knows saves splitting the text first.

### 2. Arrays

What the `#` component does in a path: the length of an array, a path read from every element, and the queries that pick the elements that match.

- **[SW-261126-6FGV](arrays.md#sw-261126-6fgv--as-the-last-component-of-a-path-returns-the-length-of-the-array-before-it)** `#` as the last component of a path returns the length of the array before it  
  Why: the count is the question a caller asks of an array most often, and answering it without materialising the elements keeps it as cheap as any other lookup.
- **[SW-261126-V9J9](arrays.md#sw-261126-v9j9-path-reads-the-path-from-every-element-and-returns-the-values-found-in-order-as-an-array)** `#.path` reads the path from every element and returns the values found, in order, as an array  
  Why: pulling one field out of every record is the array lookup callers need most, and skipping the records that lack it keeps the answer a list of real values.
- **[SW-261126-2TYJ](arrays.md#sw-261126-2tyj-query-returns-the-first-element-that-matches-the-query-and-query-every-one-as-an-array)** `#[query]` returns the first element that matches the query, and `#[query]#` every one as an array  
  Why: finding a record by a field is how a caller reaches into a list of objects without knowing positions; the two forms are the one-record and every-record questions.
- **[SW-261126-XMB5](arrays.md#sw-261126-xmb5-a-query-compares-by-the-elements-type-with--------and-matches-strings-with--and-)** A query compares by the element's type, with `==`, `=`, `!=`, `<`, `<=`, `>`, `>=`, and matches strings with `%` and `!%`  
  Why: one operator set over every JSON type keeps a query readable as a condition, and comparing by the element's own type means `age>45` means a number comparison where the data holds numbers.

### 3. Modifiers

The `@` components of a path: the built-in modifiers, the ones a caller adds, and the switch that turns them off.

- **[SW-261126-H629](modifiers.md#sw-261126-h629-a-path-starting-with-name-runs-that-modifier-and-reads-the-rest-of-the-path-from-its-output)** A path starting with `@name` runs that modifier and reads the rest of the path from its output  
  Why: a modifier reshapes the JSON where the path syntax alone cannot, and treating its output as a fresh document lets every other part of the syntax apply to it.
- **[SW-261126-4Z6X](modifiers.md#sw-261126-4z6x-a-caller-adds-a-modifier-by-name-and-a-name-no-modifier-has-gives-a-result-that-does-not-exist)** A caller adds a modifier by name, and a name no modifier has gives a result that does not exist  
  Why: a caller's own reshaping belongs in its own code, not in gjson, and an unknown name that quietly returned the input would hide a typo.
- **[SW-261126-Z3CC](modifiers.md#sw-261126-z3cc-disablemodifiers-makes--an-ordinary-character-of-a-key)** DisableModifiers makes `@` an ordinary character of a key  
  Why: a document whose keys start with `@`, as JSON-LD's do, must stay reachable by a caller who never uses modifiers.

### 4. Results

What a caller reads off a Result: its type and raw text, and the conversions to Go values, the walks over arrays and objects, and the order Less puts results in.

- **[SW-261126-PKGR](results.md#sw-261126-pkgr-string-returns-a-string-unescaped-an-integer-as-the-document-wrote-it-and-other-numbers-as-the-shortest-float)** String returns a string unescaped, an integer as the document wrote it, and other numbers as the shortest float  
  Why: String is the conversion callers use most; returning an integer's own digits keeps an ID beyond 2^53 intact, which a float cannot hold, and a string's text is the caller's data, whitespace included.
- **[SW-261126-5VFD](results.md#sw-261126-5vfd-int-and-uint-read-every-digit-of-a-64-bit-integer-and-float-reads-a-number-or-a-numeric-string)** Int and Uint read every digit of a 64-bit integer, and Float reads a number or a numeric string  
  Why: integers in JSON are often IDs and counters wider than a float64's 53 bits, and a caller asking for an integer means the integer the document wrote.
- **[SW-261126-ZGV3](results.md#sw-261126-zgv3-bool-is-true-for-true-a-number-other-than-zero-and-a-string-other-than-empty-0-and-false)** Bool is true for true, a number other than zero, and a string other than empty, `0` and `false`  
  Why: configuration documents spell flags as booleans, numbers and strings alike; a lenient reading lets a caller ask the one question whatever the spelling.
- **[SW-261126-26XP](results.md#sw-261126-26xp-exists-is-false-only-for-a-value-the-path-did-not-find-and-a-json-null-exists)** Exists is false only for a value the path did not find, and a JSON null exists  
  Why: null and absent mean different things to a caller patching a record, and the Type alone cannot tell them apart.
- **[SW-261126-NKTP](results.md#sw-261126-nktp-array-map-and-foreach-walk-an-arrays-elements-or-an-objects-members-in-document-order)** Array, Map and ForEach walk an array's elements or an object's members in document order  
  Why: the walks are how a caller goes from one lookup to all the values under it, and document order is the order the author meant.
- **[SW-261126-NRBX](results.md#sw-261126-nrbx-value-returns-the-result-as-the-go-value-encodingjson-would-decode-it-into-an-interface)** Value returns the result as the Go value encoding/json would decode it into an interface  
  Why: a caller handing a value on to code that takes interface{} expects the shapes encoding/json produces, and gets them without a second decode.
- **[SW-261126-Z23S](results.md#sw-261126-z23s-parse-reads-a-whole-document-as-a-result-and-resultget-reads-a-path-inside-one)** Parse reads a whole document as a result, and Result.Get reads a path inside one  
  Why: a caller that reads many paths from one value parses once and walks from there, instead of naming the whole path every time.
- **[SW-261126-4WX9](results.md#sw-261126-4wx9-less-orders-results-by-type-null-false-number-string-true-json-and-within-a-type-by-value)** Less orders results by type, Null, False, Number, String, True, JSON, and within a type by value  
  Why: a caller sorting results of mixed types needs one total order that does not depend on how the document spelled them.

### 5. Validity

What gjson does with a document it cannot trust: Valid to ask, a reader that does not validate and never panics, and the byte and many-path forms of Get that answer as Get does.

- **[SW-261126-KNNZ](validity.md#sw-261126-knnz-valid-reports-whether-a-document-is-exactly-one-json-value-as-the-json-grammar-has-it)** Valid reports whether a document is exactly one JSON value, as the JSON grammar has it  
  Why: Get does not validate; Valid is the one place a caller learns whether a document is well formed, and it must be strict for the answer to be worth asking.
- **[SW-261126-9CTW](validity.md#sw-261126-9ctw-get-does-not-validate-on-a-malformed-document-it-returns-what-it-finds-and-it-never-panics)** Get does not validate: on a malformed document it returns what it finds, and it never panics *(invariant)*  
  Why: validating would make every lookup read the whole document, the cost gjson exists to avoid; the caller who needs validation asks Valid, and the one promise on bad input is that the program survives it.
- **[SW-261126-3FK5](validity.md#sw-261126-3fk5-getbytes-and-getmany-return-what-get-returns-for-the-same-document-and-paths)** GetBytes and GetMany return what Get returns for the same document and paths  
  Why: the byte and many-path forms exist for speed, and a caller swaps one for another only if the answers stay the same.
