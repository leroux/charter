# wc

Code: cmd/wc

What wc counts in an input, which counts it prints, and how it lays them
out.

## SW-261126-4W6F wc counts newlines, words of printable bytes, and bytes

Status: active
Kind: invariant
Serves: STK-261126-69ZV

The line count is the number of newline bytes, so a last line without
its newline is not counted. A word is a run of bytes between white space
(space, tab, newline, vertical tab, form feed, carriage return) holding
at least one printable byte (0x21 to 0x7e); a byte that is neither
printable nor white space neither starts a word nor ends one. The byte
count is the input's length.

Why: this is GNU 9.4's count in the C locale, and a script that counts
records with wc -l relies on the newline rule exactly.

Scenarios:
1. `lines-are-newlines` (agent, provisional): "a\nb" counts one line.
2. `words-printable` (agent, provisional): "\x01" alone is no word,
   "a\x01b c" is two, and tabs and carriage returns separate words.
3. `bytes-are-length` (agent, provisional): the byte count is the
   input's length, newlines and high bytes included.

## SW-261126-MKWM wc prints the counts -l, -w and -c select, all three when none is, always in that order

Status: active
Serves: STK-261126-69ZV

-l selects the line count, -w the word count and -c the byte count; with
none of them all three are printed. Whatever order the options come in,
the counts are printed lines, words, bytes.

Why: GNU's order is fixed, and a script that reads the second column as
bytes after `wc -cl` relies on it.

Scenarios:
1. `none-is-all` (agent, provisional): with no option the line has the
   lines, words and bytes.
2. `order-fixed` (agent, provisional): -cl prints the lines, then the
   bytes.

## SW-261126-SJ3N wc right-aligns every count to one width and prints the file's name after

Status: active
Serves: STK-261126-69ZV

Counts are separated by a space and right-aligned to one width: 1 when a
single count of a single input is printed; else at least 7 when an input
is standard input or not a regular file, otherwise the number of digits
of the inputs' total size in bytes. A file's name follows its counts
after a space; standard input has no name.

Why: this is GNU's layout; `wc -l < f` printing a bare number is what
every script depends on, and the padded form is what people's eyes and
awk's fields are used to.

Scenarios:
1. `single-count-bare` (agent, provisional): wc -l on standard input
   prints the bare number.
2. `stdin-width-seven` (agent, provisional): wc on standard input pads
   each count to seven columns.
3. `file-width-from-size` (agent, provisional): wc on a file of twelve
   bytes pads to two columns.
4. `name-after` (agent, provisional): wc -l on a file prints the count,
   a space and the name.
