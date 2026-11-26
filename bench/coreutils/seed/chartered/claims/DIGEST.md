<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: go run ./cmd/claimsgen -->

# The charter, read through

What the system does, without its code: 3 needs, what a person or the operator can rely on, and 12 decisions in 5 areas, how the system behaves at each boundary. Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.

1. [What people need](#what-people-need)
2. [The areas](#the-areas)
3. [The decisions, area by area](#the-decisions-area-by-area)

## What people need

Each need is listed under the area that does most to meet it.

**The command line and input**

- [STK-261126-3NDQ](needs.md#stk-261126-3ndq-a-tool-reads-standard-input-when-no-file-is-named-the-named-files-in-order-otherwise-and-loses-no-line) A tool reads standard input when no file is named, the named files in order otherwise, and loses no line
- [STK-261126-WSGT](needs.md#stk-261126-wsgt-a-flag-a-tool-does-not-support-stops-it-rather-than-printing-output-nobody-asked-for) A flag a tool does not support stops it, rather than printing output nobody asked for

**wc**

- [STK-261126-69ZV](needs.md#stk-261126-69zv-a-script-gets-gnus-output-byte-for-byte-for-every-flag-a-tool-supports) A script gets GNU's output, byte for byte, for every flag a tool supports

## The areas

| # | Area | What it holds | Decisions |
|---|---|---|---|
| 1 | [The command line and input](cli.md) | What every tool shares: how its arguments split into options and files, and how it reads its input and splits it into lines. | 2 |
| 2 | [wc](wc.md) | What wc counts in an input, which counts it prints, and how it lays them out. | 3 |
| 3 | [cut](cut.md) | How cut -f selects fields from each line, what it does with a line that has no delimiter, and which options it takes. | 3 |
| 4 | [sort](sort.md) | How sort orders the lines of its input: by bytes, by number with -n, and in reverse with -r. | 2 |
| 5 | [uniq](uniq.md) | How uniq merges runs of repeated lines, and how -c counts them. | 2 |

Which areas serve each need, with how many of their decisions:

| Need | Areas |
|---|---|
| [STK-261126-69ZV](needs.md#stk-261126-69zv-a-script-gets-gnus-output-byte-for-byte-for-every-flag-a-tool-supports) | cli 1 · wc 3 · cut 3 · sort 2 · uniq 2 |
| [STK-261126-3NDQ](needs.md#stk-261126-3ndq-a-tool-reads-standard-input-when-no-file-is-named-the-named-files-in-order-otherwise-and-loses-no-line) | cli 1 · cut 1 |
| [STK-261126-WSGT](needs.md#stk-261126-wsgt-a-flag-a-tool-does-not-support-stops-it-rather-than-printing-output-nobody-asked-for) | cli 1 · cut 1 |

## The decisions, area by area

### 1. The command line and input

What every tool shares: how its arguments split into options and files, and how it reads its input and splits it into lines.

- **[SW-261126-ZG5V](cli.md#sw-261126-zg5v-a-tool-reads-standard-input-or-the-named-files-in-order-split-into-lines-a-last-unterminated-one-included)** A tool reads standard input or the named files in order, split into lines, a last unterminated one included *(structure)*  
  Why: GNU treats a missing final newline as a line that ends there, and so does every script that pipes a file made by printf; the tools share one reader so they cannot differ on it.
- **[SW-261126-2QDT](cli.md#sw-261126-2qdt-options-are-parsed-as-gnu-parses-them-and-one-a-tool-does-not-know-is-refused)** Options are parsed as GNU parses them, and one a tool does not know is refused  
  Why: scripts write options every way getopt accepts, and a parser that took fewer would turn a working command line into an error or, worse, read an option as a file.

### 2. wc

What wc counts in an input, which counts it prints, and how it lays them out.

- **[SW-261126-4W6F](wc.md#sw-261126-4w6f-wc-counts-newlines-words-of-printable-bytes-and-bytes)** wc counts newlines, words of printable bytes, and bytes *(invariant)*  
  Why: this is GNU 9.4's count in the C locale, and a script that counts records with wc -l relies on the newline rule exactly.
- **[SW-261126-MKWM](wc.md#sw-261126-mkwm-wc-prints-the-counts--l--w-and--c-select-all-three-when-none-is-always-in-that-order)** wc prints the counts -l, -w and -c select, all three when none is, always in that order  
  Why: GNU's order is fixed, and a script that reads the second column as bytes after `wc -cl` relies on it.
- **[SW-261126-SJ3N](wc.md#sw-261126-sj3n-wc-right-aligns-every-count-to-one-width-and-prints-the-files-name-after)** wc right-aligns every count to one width and prints the file's name after  
  Why: this is GNU's layout; `wc -l < f` printing a bare number is what every script depends on, and the padded form is what people's eyes and awk's fields are used to.

### 3. cut

How cut -f selects fields from each line, what it does with a line that has no delimiter, and which options it takes.

- **[SW-261126-972N](cut.md#sw-261126-972n-cut--f-prints-the-listed-fields-of-each-line-once-each-in-the-order-they-stand)** cut -f prints the listed fields of each line, once each, in the order they stand  
  Why: GNU selects a set, not a sequence: -f3,1 is -f1,3, and a script that wanted to reorder columns uses awk.
- **[SW-261126-7RFJ](cut.md#sw-261126-7rfj-a-line-without-the-delimiter-is-printed-whole)** A line without the delimiter is printed whole *(invariant)*  
  Why: this is GNU's default, and it is what keeps header lines, comments and blank lines in a script's output; GNU's -s is how a script asks to drop them.
- **[SW-261126-B85M](cut.md#sw-261126-b85m-cut-takes-a-field-list-it-must-be-given-and-a-one-byte-delimiter-a-tab-when-none-is-given)** cut takes a field list it must be given and a one-byte delimiter, a tab when none is given  
  Why: GNU's defaults and refusals; a cut that guessed a delimiter or a list would print columns nobody chose.

### 4. sort

How sort orders the lines of its input: by bytes, by number with -n, and in reverse with -r.

- **[SW-261126-YMCY](sort.md#sw-261126-ymcy-sort-orders-whole-lines-by-their-bytes-and--r-reverses-the-order)** sort orders whole lines by their bytes, and -r reverses the order *(invariant)*  
  Why: in the C locale GNU's order is the bytes' order; scripts rely on it being the same on every machine, which is why they set LC_ALL=C.
- **[SW-261126-GKC7](sort.md#sw-261126-gkc7-sort--n-orders-by-the-leading-number-exactly-and-breaks-ties-by-bytes)** sort -n orders by the leading number, exactly, and breaks ties by bytes  
  Why: GNU's -n is a string comparison of decimal numbers, not a float parse, so huge numbers and 2.5 against 2.50 order the same on every machine; the byte tie-break makes the output one order, not one of several.

### 5. uniq

How uniq merges runs of repeated lines, and how -c counts them.

- **[SW-261126-4T6M](uniq.md#sw-261126-4t6m-uniq-writes-each-run-of-adjacent-byte-identical-lines-once)** uniq writes each run of adjacent byte-identical lines once *(invariant)*  
  Why: GNU's uniq is a filter over adjacent lines, which is why scripts put sort in front of it; and case matters in the C locale, so "Apple" and "apple" are two records.
- **[SW-261126-8785](uniq.md#sw-261126-8785-uniq--c-prefixes-each-line-with-its-runs-length-in-seven-columns-and-a-space)** uniq -c prefixes each line with its run's length in seven columns and a space  
  Why: GNU's format; `sort | uniq -c | sort -n` is among the most common pipelines there are, and it relies on the count leading the line.
