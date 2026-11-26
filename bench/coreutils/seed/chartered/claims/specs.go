// Code generated from claims by the charter gate. DO NOT EDIT.
// Regenerate: go run ./cmd/claimsgen

package claims

import "github.com/leroux/charter"

// STK_261126_3NDQ is STK-261126-3NDQ: A tool reads standard input when no file is named, the named files in order otherwise, and loses no line
//
// With no file named a tool reads standard input; with files named it reads
// each in the order given. Every line of the input reaches the output it
// belongs in, the last one included when it lacks its newline.
var STK_261126_3NDQ = &charter.Spec{ID: "STK-261126-3NDQ", Title: "A tool reads standard input when no file is named, the named files in order otherwise, and loses no line", Status: charter.Active}

// STK_261126_69ZV is STK-261126-69ZV: A script gets GNU's output, byte for byte, for every flag a tool supports
//
// For every flag a tool supports, its standard output and exit status are
// GNU's for the same arguments and input: the same bytes, the same
// padding, the same order. A tool supports fewer flags than GNU, never a
// different meaning for one.
var STK_261126_69ZV = &charter.Spec{ID: "STK-261126-69ZV", Title: "A script gets GNU's output, byte for byte, for every flag a tool supports", Status: charter.Active}

// STK_261126_WSGT is STK-261126-WSGT: A flag a tool does not support stops it, rather than printing output nobody asked for
//
// Given an option it does not know, a tool writes nothing to standard
// output, names the option on standard error and exits with a failing
// status.
var STK_261126_WSGT = &charter.Spec{ID: "STK-261126-WSGT", Title: "A flag a tool does not support stops it, rather than printing output nobody asked for", Status: charter.Active}

// SW_261126_2QDT is SW-261126-2QDT: Options are parsed as GNU parses them, and one a tool does not know is refused
//
// Short options may be clustered (-lw) and may come before, between or
// after the file operands; "--" ends the options and "-" alone is an
// operand. An option that takes a value takes the rest of its cluster
// (-d:) or the next argument (-d :). A long option is --name or
// --name=value. An option the tool does not know, or a value missing,
// is an error in GNU's words, and the tool stops before reading input:
// exit status 1, or 2 for sort, as GNU's.
var SW_261126_2QDT = &charter.Spec{ID: "SW-261126-2QDT", Title: "Options are parsed as GNU parses them, and one a tool does not know is refused", Status: charter.Active}

// SW_261126_4T6M is SW-261126-4T6M: uniq writes each run of adjacent byte-identical lines once
//
// Lines that are the same, byte for byte, and adjacent form a run, and uniq
// writes the run's line once. Equal lines with another line between them
// are separate runs. Lines that differ only in case differ.
var SW_261126_4T6M = &charter.Spec{ID: "SW-261126-4T6M", Title: "uniq writes each run of adjacent byte-identical lines once", Status: charter.Active, Kind: "invariant"}

// SW_261126_4W6F is SW-261126-4W6F: wc counts newlines, words of printable bytes, and bytes
//
// The line count is the number of newline bytes, so a last line without
// its newline is not counted. A word is a run of bytes between white space
// (space, tab, newline, vertical tab, form feed, carriage return) holding
// at least one printable byte (0x21 to 0x7e); a byte that is neither
// printable nor white space neither starts a word nor ends one. The byte
// count is the input's length.
var SW_261126_4W6F = &charter.Spec{ID: "SW-261126-4W6F", Title: "wc counts newlines, words of printable bytes, and bytes", Status: charter.Active, Kind: "invariant"}

// SW_261126_7RFJ is SW-261126-7RFJ: A line without the delimiter is printed whole
//
// cut prints a line that holds no delimiter as it stands, whatever the
// field list.
var SW_261126_7RFJ = &charter.Spec{ID: "SW-261126-7RFJ", Title: "A line without the delimiter is printed whole", Status: charter.Active, Kind: "invariant"}

// SW_261126_8785 is SW-261126-8785: uniq -c prefixes each line with its run's length in seven columns and a space
//
// With -c, each line uniq writes is preceded by the length of its run,
// right-aligned in seven columns, and a space.
var SW_261126_8785 = &charter.Spec{ID: "SW-261126-8785", Title: "uniq -c prefixes each line with its run's length in seven columns and a space", Status: charter.Active}

// SW_261126_972N is SW-261126-972N: cut -f prints the listed fields of each line, once each, in the order they stand
//
// A field list is ranges separated by commas: N, N-M, N- (to the end of
// the line) and -M (from the first), fields numbered from one. A line's
// fields are the pieces between delimiters; cut prints the fields any
// range selects, each once and in the order they stand in the line, joined
// by the delimiter. A field past the end of the line is not printed, and a
// line none of whose fields are selected prints as an empty line.
var SW_261126_972N = &charter.Spec{ID: "SW-261126-972N", Title: "cut -f prints the listed fields of each line, once each, in the order they stand", Status: charter.Active}

// SW_261126_B85M is SW-261126-B85M: cut takes a field list it must be given and a one-byte delimiter, a tab when none is given
//
// -f names the field list and is required; without it cut stops with
// "you must specify a list of bytes, characters, or fields". -d names the
// delimiter, one byte; a longer one stops cut with "the delimiter must be a
// single character". Without -d the delimiter is a tab.
var SW_261126_B85M = &charter.Spec{ID: "SW-261126-B85M", Title: "cut takes a field list it must be given and a one-byte delimiter, a tab when none is given", Status: charter.Active}

// SW_261126_GKC7 is SW-261126-GKC7: sort -n orders by the leading number, exactly, and breaks ties by bytes
//
// -n reads each line's leading number: blanks, an optional minus, digits
// and an optional point with more digits; a plus sign, an exponent or a
// thousands separator ends it. A line that does not start with a number
// reads as zero, as does negative zero. Numbers compare exactly, however
// many digits they have. Lines whose numbers are equal compare by bytes,
// and -r reverses that comparison too.
var SW_261126_GKC7 = &charter.Spec{ID: "SW-261126-GKC7", Title: "sort -n orders by the leading number, exactly, and breaks ties by bytes", Status: charter.Active}

// SW_261126_MKWM is SW-261126-MKWM: wc prints the counts -l, -w and -c select, all three when none is, always in that order
//
// -l selects the line count, -w the word count and -c the byte count; with
// none of them all three are printed. Whatever order the options come in,
// the counts are printed lines, words, bytes.
var SW_261126_MKWM = &charter.Spec{ID: "SW-261126-MKWM", Title: "wc prints the counts -l, -w and -c select, all three when none is, always in that order", Status: charter.Active}

// SW_261126_SJ3N is SW-261126-SJ3N: wc right-aligns every count to one width and prints the file's name after
//
// Counts are separated by a space and right-aligned to one width: 1 when a
// single count of a single input is printed; else at least 7 when an input
// is standard input or not a regular file, otherwise the number of digits
// of the inputs' total size in bytes. A file's name follows its counts
// after a space; standard input has no name.
var SW_261126_SJ3N = &charter.Spec{ID: "SW-261126-SJ3N", Title: "wc right-aligns every count to one width and prints the file's name after", Status: charter.Active}

// SW_261126_YMCY is SW-261126-YMCY: sort orders whole lines by their bytes, and -r reverses the order
//
// Without -n, lines compare as unsigned byte strings, as the C locale's
// collation does: uppercase before lowercase, a prefix before the longer
// line. -r reverses the comparison. sort reads every input before it
// prints, the lines of all its files together.
var SW_261126_YMCY = &charter.Spec{ID: "SW-261126-YMCY", Title: "sort orders whole lines by their bytes, and -r reverses the order", Status: charter.Active, Kind: "invariant"}

// SW_261126_ZG5V is SW-261126-ZG5V: A tool reads standard input or the named files in order, split into lines, a last unterminated one included
//
// A tool with no file operand reads standard input; with file operands it
// reads each file in the order named. A tool splits its input on newline
// bytes into lines; a last line without a newline is a line all the same,
// and an empty input has no lines. Every line a tool prints ends with a
// newline.
var SW_261126_ZG5V = &charter.Spec{ID: "SW-261126-ZG5V", Title: "A tool reads standard input or the named files in order, split into lines, a last unterminated one included", Status: charter.Active, Kind: "structure"}
