// Code generated from claims by the charter gate. DO NOT EDIT.
// Regenerate: make claims

package claims

import "github.com/leroux/charter"

// CH_001 is CH-001: A reviewer can judge a change by its entries and their tests, not by reading the code
//
// Every clause of an active decision is a scenario, every scenario has a
// test that names it, and the gate fails on the first one that does not, so
// the reviewer reads the area document diff and the gate's verdict and
// opens the code only where the review points. The gate is static and runs
// in seconds, and its verdict is the same on every machine.
var CH_001 = &charter.Spec{ID: "CH-001", Title: "A reviewer can judge a change by its entries and their tests, not by reading the code", Status: charter.Active}

// CH_002 is CH-002: A builder reaches the code and the tests from an id, and the gate names the fix for every problem
//
// From an entry's id the tooling lists the code that realizes it, by file,
// line and declaration, and the tests that verify it, by file, test and
// scenario. Every problem the gate reports names a file, a line and what to
// do. Symbols are generated, so a citation of an entry that does not exist
// fails to compile.
var CH_002 = &charter.Spec{ID: "CH-002", Title: "A builder reaches the code and the tests from an id, and the gate names the fix for every problem", Status: charter.Active}

// CH_003 is CH-003: The record is stable: ids never change or get reused, retired entries keep their history, and changes are dated
//
// An entry keeps its id for life; a retired one stays in its document with
// its status and its changes, and its symbol still builds, so the strip can
// follow at its own pace. Two entries with one id, anywhere in the charter,
// fail. A change to an entry is a dated line in its Changes.
var CH_003 = &charter.Spec{ID: "CH-003", Title: "The record is stable: ids never change or get reused, retired entries keep their history, and changes are dated", Status: charter.Active}

// CH_004 is CH-004: No surface exists that an entry does not own, and the baseline only shrinks
//
// Every environment variable, build tag, contract operation and event kind
// the tree exposes is named by an entry's Accounts, or listed with a reason
// in a baseline that predates the charter. The gate fails a surface neither
// names, a baseline line whose surface is gone or now accounted for, and a
// read of an environment variable that goes through anything but a named
// constant.
var CH_004 = &charter.Spec{ID: "CH-004", Title: "No surface exists that an entry does not own, and the baseline only shrinks", Status: charter.Active}

// CH_005 is CH-005: An external system's fact is recorded with its source and date, and no active entry rests on one found false
//
// An observation records a fact about a system the code does not own, with
// its subject, where it was read or when a probe last saw it true, and the
// entries that rest on it. Probes check observations and nothing else; tests
// cite entries and never observations. A broken observation under an active
// entry fails the gate.
var CH_005 = &charter.Spec{ID: "CH-005", Title: "An external system's fact is recorded with its source and date, and no active entry rests on one found false", Status: charter.Active}

// CH_006 is CH-006: The charter specifies itself, and its own gate is green
//
// The tooling's decisions are entries in this directory, in the form the
// charter asks of any adopter, realized by the tooling and verified by its
// tests, and checked by a gate of their own over the tooling alone.
var CH_006 = &charter.Spec{ID: "CH-006", Title: "The charter specifies itself, and its own gate is green", Status: charter.Active}

// CH_010 is CH-010: An entry is a heading with id and title, then field lines, text, Why, Scenarios, Guidance and Changes in order
//
// An area document is a markdown file whose header is a title and an
// optional Code line, and whose entries are headings of the form
// "## ID Title". The lines after the heading, up to a blank line, are
// fields: Status, Kind, Serves, See and Accounts, each at most once; the
// values a Status takes are CH-013's. The body is paragraphs: the text
// first, then paragraphs opened by Why, Scenarios, Guidance and Changes, in
// that order and each at most once; the text and Why are required, and an
// active entry lists at least one scenario. A scenario is a numbered item
// with a handle in backticks, lower-case kebab, unique in the entry, and a
// sentence; a change is a dated item; a continuation line is indented. Every
// form problem names the document, the line and the form to write. The
// header, before the first entry, may hold prose, a vocabulary list, tables
// and diagrams that orient the reader; it is not parsed beyond its title and
// Code line and has no section headings, and a line inside a fenced block,
// in the header or in an entry, is text and never a heading.
var CH_010 = &charter.Spec{ID: "CH-010", Title: "An entry is a heading with id and title, then field lines, text, Why, Scenarios, Guidance and Changes in order", Status: charter.Active}

// CH_011 is CH-011: A need lives in needs.md and has no Kind, Serves or Guidance
//
// An entry of the document named needs.md is a need. A need serves nothing,
// carries no kind and no guidance; each of those is a form problem.
var CH_011 = &charter.Spec{ID: "CH-011", Title: "A need lives in needs.md and has no Kind, Serves or Guidance", Status: charter.Active}

// CH_012 is CH-012: An id is declared once; a new one is its prefix, the day it was taken and a random code
//
// An id is a prefix and then either a number, as SW-118, or the day it was
// taken and a code, as SW-261001-K7Q2; a heading whose id is neither is a
// form problem. The day is the UTC date as six digits, year, month and
// day; the code is four characters of BCDFGHJKMNPQRSTVWXYZ23456789, the
// consonants and digits that spell no word and are not read as one
// another. The numbered ids stand as they are, and next gives every new
// id, of the second kind and one the tree does not hold, as CH-040 rules.
// Which ids may name a spec file of the earlier form is CH-014's. Ids are
// unique across the area documents and the spec files of the earlier form;
// a second declaration fails naming both places and saying that the entry
// that has not landed takes a new id. The entries of the area documents are
// ordered by prefix, and within one the numbered ids first, by number, then
// the dated ones by date and code, and the generated Go and TypeScript
// module declare them in that order.
var CH_012 = &charter.Spec{ID: "CH-012", Title: "An id is declared once; a new one is its prefix, the day it was taken and a random code", Status: charter.Active}

// CH_013 is CH-013: Status is planned, active or deprecated, every status generates a symbol, and coverage is asked of active entries only
//
// An entry's Status is planned, active or deprecated, and any other is a
// form problem. A planned entry is ruled and not built; its symbol is
// generated, a test may cite it, and the gate asks nothing else of it. An
// active entry must be covered. A deprecated entry is retired; its symbol
// is generated, and the gate asks nothing of it. How the generated files
// mark a planned or a deprecated entry is CH-020's.
var CH_013 = &charter.Spec{ID: "CH-013", Title: "Status is planned, active or deprecated, every status generates a symbol, and coverage is asked of active entries only", Status: charter.Active}

// CH_014 is CH-014: The earlier form, one file per spec, is read beside the area documents until it is migrated
//
// A spec file in the spec directory, named for its numbered id, never a
// dated one, with a lens and a slug, is parsed into the same registry as an
// area entry: its fields in their fixed order, its Realizes counting as
// Serves, and its links pointing at the target's current file. The directory
// may be absent.
var CH_014 = &charter.Spec{ID: "CH-014", Title: "The earlier form, one file per spec, is read beside the area documents until it is migrated", Status: charter.Active}

// CH_015 is CH-015: Serves and See name existing entries, a need serves nothing, and nothing names itself
//
// Every id under Serves or See, and under the earlier form's Realizes,
// Concerns and Related, names an entry the charter holds; an entry naming
// itself fails; a stakeholder spec of the earlier form realizes nothing. An
// observation's Supports names existing entries too.
var CH_015 = &charter.Spec{ID: "CH-015", Title: "Serves and See name existing entries, a need serves nothing, and nothing names itself", Status: charter.Active}

// CH_020 is CH-020: Every entry has a generated Go and TypeScript symbol, its id with underscores for hyphens, its title the doc comment
//
// The entries package's generated Go declares one variable per entry of any
// status, named by the id with each hyphen as an underscore, holding the id,
// the title, the status and the kind, with the title and the text as its
// doc comment; a planned entry is marked planned, a deprecated one carries a
// deprecation notice naming its document. The TypeScript module maps the
// same symbols to the ids, grouped by prefix, with the title as a doc
// comment. The order both declare the entries in is CH-012's. Each file must
// equal its rendering, or the gate names the command to run; a variable of
// each name makes the gate write the file instead.
var CH_020 = &charter.Spec{ID: "CH-020", Title: "Every entry has a generated Go and TypeScript symbol, its id with underscores for hyphens, its title the doc comment", Status: charter.Active}

// CH_021 is CH-021: Production code anchors an entry by a directive in a declaration's or package's doc comment, or a TypeScript line comment
//
// A Go directive is a comment line "//charter:realizes ID, ..." or
// "//charter:concerns ID, ...", with or without a space after the slashes,
// in the doc comment of a package clause or of a declaration; the anchor
// records the kind, the file, the line and the declaration it documents. A
// TypeScript directive is the same text in a line comment of a file that is
// not a test; the anchor records the file and the line.
var CH_021 = &charter.Spec{ID: "CH-021", Title: "Production code anchors an entry by a directive in a declaration's or package's doc comment, or a TypeScript line comment", Status: charter.Active}

// CH_022 is CH-022: An anchor names an entry that exists and is not planned, and concerns never counts toward coverage
//
// A directive naming an id that is not an entry fails; one naming a planned
// entry fails, saying to set it active as its work begins; a directive with
// no id fails. A concerns anchor is navigation: an active decision with only
// a concerns anchor is unrealized.
var CH_022 = &charter.Spec{ID: "CH-022", Title: "An anchor names an entry that exists and is not planned, and concerns never counts toward coverage", Status: charter.Active}

// CH_023 is CH-023: A test cites an entry by a Verifies statement in a test function, naming by handle the scenarios it exercises
//
// A Go citation is a call of Verifies inside a function whose name starts
// with Test, whose second argument is the entry's symbol and whose further
// arguments are string literals naming the entry's scenarios. A TypeScript
// citation is a call of verifies inside a test call, with the same shape.
// The calls open the test, or a subtest it runs: they are its first
// statements, one after another, and a call that anything precedes or
// encloses, or in a function go test never runs, such as Testx, fails. A TypeScript call in a test
// marked skip, fixme, fail or todo, or to a verifies the test file declares
// itself, fails. A call anywhere else, a spec named by anything but its
// symbol, a handle that is not a literal, or a handle the entry does not
// list, each fail naming the file and the line. Any other reference to an entry from a test
// file is a reference, not a citation.
var CH_023 = &charter.Spec{ID: "CH-023", Title: "A test cites an entry by a Verifies statement in a test function, naming by handle the scenarios it exercises", Status: charter.Active}

// CH_024 is CH-024: Verifies and verifies mark the running test with each citation, "ID" or "ID/handle", in the runner's structured output
//
// The Go marker is a test attribute per citation, which go test's JSON
// stream reports as an attr event on the test before its result; the
// helper a TypeScript test runner calls may annotate the running test the
// same way. A nil spec is refused at run time. The report reads these events.
var CH_024 = &charter.Spec{ID: "CH-024", Title: "Verifies and verifies mark the running test with each citation, \"ID\" or \"ID/handle\", in the runner's structured output", Status: charter.Active}

// CH_025 is CH-025: The entries package holds only observations and the generated file, and no entry is assigned outside its declaration
//
// A hand-written file of the entries package declares observations as
// addressed composite literals with keyed, literal fields and nothing else:
// no function, no other declaration, no spec written in Go. No test file
// assigns to an entry or a field of one, and no test dot-imports the
// package.
var CH_025 = &charter.Spec{ID: "CH-025", Title: "The entries package holds only observations and the generated file, and no entry is assigned outside its declaration", Status: charter.Active}

// CH_030 is CH-030: An active decision is anchored and each scenario cited; a waiver absorbs a missing test, never a missing anchor
//
// An active decision with no realizes anchor fails, and nothing waives it;
// that a concerns anchor is not one is CH-022's. An active entry with
// scenarios fails for each scenario no citation names, naming the document,
// the line and the call to write; one without scenarios fails when no
// citation names it at all. The waiver file lists, one per line with a
// reason after a tab, an id or a scenario as ID/handle, never a glob; a
// matching line absorbs the missing test; a line that absorbs nothing, gives
// no reason or is a glob fails.
var CH_030 = &charter.Spec{ID: "CH-030", Title: "An active decision is anchored and each scenario cited; a waiver absorbs a missing test, never a missing anchor", Status: charter.Active, Kind: "invariant"}

// CH_031 is CH-031: A need needs no anchor and is verified only from the stakeholder's side
//
// A need is realized by the decisions that serve it, so the gate asks no
// anchor of it. When the options name the test files on the stakeholder's
// side, a citation of a need from any other file fails naming the file, and
// does not count; when they name none, any test file counts.
var CH_031 = &charter.Spec{ID: "CH-031", Title: "A need needs no anchor and is verified only from the stakeholder's side", Status: charter.Active}

// CH_032 is CH-032: A probe cites observations only, and a test never cites an observation
//
// A probe is a test file whose build constraint names the probe tag on a
// positive path. A probe referencing an entry fails, and its references are
// no citations; a test file that is not a probe referencing an observation
// fails. Without a probe tag in the options, observations are one
// configuration problem and no probe rule runs.
var CH_032 = &charter.Spec{ID: "CH-032", Title: "A probe cites observations only, and a test never cites an observation", Status: charter.Active}

// CH_033 is CH-033: Every surface is enumerated, and accounted for by an entry or listed in a baseline of only live, unaccounted surfaces
//
// The surfaces are the exported string constants of the env packages, every
// build tag but the probe tag and version tags, every operation of the
// OpenAPI files, and every kind of the events file; an adopter may add
// kinds. A surface no entry accounts for and the baseline does not list
// fails; an entry accounting for a surface nothing enumerates fails; a
// baseline line whose surface is gone, or which an entry now accounts for,
// or which appears twice, or which lacks a reason, fails. An enumerator
// configured by a path that finds nothing is a misconfiguration; zero build
// tags is a valid tree. With the update variable set the gate rewrites the
// baseline with every unaccounted surface, keeping the reasons it had. A
// deprecated entry's Accounts own nothing: what it accounted for is a live
// entry's to claim, or the baseline's.
var CH_033 = &charter.Spec{ID: "CH-033", Title: "Every surface is enumerated, and accounted for by an entry or listed in a baseline of only live, unaccounted surfaces", Status: charter.Active, Kind: "invariant"}

// CH_034 is CH-034: An environment variable is read through a named constant, or under a dynamic directive, and no literal names one
//
// Outside the env packages and outside test files, a call of a configured
// reader takes a constant of an env package as its first argument, or the
// line carries the dynamic directive. A string literal equal to a declared
// variable's value, or matching a configured pattern, fails.
var CH_034 = &charter.Spec{ID: "CH-034", Title: "An environment variable is read through a named constant, or under a dynamic directive, and no literal names one", Status: charter.Active, Kind: "invariant"}

// CH_035 is CH-035: Observations carry their status's fields, observed ones a date and a probe; one broken under an active entry fails
//
// Every observation names a subject. An observed one has a calendar date and
// at least one citing probe; a documented one has a source; any date given
// is a calendar date. A broken observation supporting an active entry fails;
// under a planned or a deprecated entry it does not.
var CH_035 = &charter.Spec{ID: "CH-035", Title: "Observations carry their status's fields, observed ones a date and a probe; one broken under an active entry fails", Status: charter.Active}

// CH_036 is CH-036: The walk covers every Go and TypeScript file under the root within the scope, skipping what the options ignore
//
// The walk parses every Go file under the root and, when TypeScript is
// configured, lists every .ts and .tsx file that is not a declaration file.
// It skips .git, node_modules, vendor and testdata, hidden and underscore
// directories, nested modules, every Ignore pattern matched against the
// path and the base name, and everything outside Only when Only names
// paths.
var CH_036 = &charter.Spec{ID: "CH-036", Title: "The walk covers every Go and TypeScript file under the root within the scope, skipping what the options ignore", Status: charter.Active}

// CH_037 is CH-037: The gate reports every problem, sorted, with its file, line, rule and fix, and passes on a clean tree
//
// Check reports each problem as one test error; Problems returns the same
// list, sorted. The options resolve the module root from the calling file or
// the working directory, require the entries package inside the module, and
// validate the TypeScript settings.
var CH_037 = &charter.Spec{ID: "CH-037", Title: "The gate reports every problem, sorted, with its file, line, rule and fix, and passes on a clean tree", Status: charter.Active}

// CH_040 is CH-040: claimsgen generates the files, lints, gives new ids and lists an entry's anchors, on any charter
//
// Without a command claimsgen writes the generated Go and the TypeScript
// module from the documents alone, so it works when a generated file is
// stale or missing. lint prints the gate's problems and then the advice.
// next prints a new id for a prefix, the day's UTC date and a code drawn at
// random, drawn again when the tree already holds that id, or as many as the
// count after the prefix asks for, each different. anchors prints, per id,
// the document, the title, the status, the scenarios, the realizing and
// concerning sites with their declarations, the verifying tests with the
// handles each names, and the relations. A repository's own command passes
// the charters it keeps, each by name with its gate's options; the charter
// flag selects one by name, the first when none is named, and refuses a name
// it was not given. A charter whose options name no root is read from the
// module above the working directory, never the module of the tooling the
// command was built with.
var CH_040 = &charter.Spec{ID: "CH-040", Title: "claimsgen generates the files, lints, gives new ids and lists an entry's anchors, on any charter", Status: charter.Active}

// CH_041 is CH-041: Lint advises on what the gate does not fail
//
// Lint runs the whole gate first and prints every problem it fails on,
// then one line counting the open questions of the question file, then
// its advice: a decision that serves nothing, a planned entry each of
// whose scenarios is verified, a deprecated entry
// still anchored or cited, an active decision realized only by a package
// clause unless its kind is structure, and a citation of an entry with
// scenarios that names none. Then it checks the style a reader would
// stumble on, by the bars of the README's Style, in the live entries of
// the area documents and the needs: a word retired.txt in the entries
// directory retires; a title over twenty words; a text over six sentences;
// a scenario sentence over thirty words; a Why under twenty-five words that
// names a story, a ruling or a review; and a test whose citations name
// scenarios of two or more live decisions and no need. A retired word is
// read in a title, a text, a Why, a scenario's sentence and a document's
// header, never in a code span or a double-quoted string, and is kept
// where the line or the sentence holding it holds one of the word's
// phrases. The file holds a word a line: a pattern matched at word
// boundaries whatever the case, then -> and the word to write, then | and
// the phrases, separated by semicolons. Each style check that finds
// anything prints one line, after the rest, counting what it found and
// naming the command that lists it; with -style, lint lists every finding
// at its line in their place, a retired word with the word to write and a
// test once, at the test.
var CH_041 = &charter.Spec{ID: "CH-041", Title: "Lint advises on what the gate does not fail", Status: charter.Active}

// CH_042 is CH-042: A report joins each test's markers to its scenarios and says which passed, failed, were skipped or have none
//
// The report runs the Go suite with structured output, reads the citation
// markers each test emitted, and writes, per entry and per scenario, the
// citing tests and their result: passed when a citing test passed and none
// failed, failed when one failed, skipped when every one was skipped, none
// when no citing test ran. An active scenario, or an active entry without
// scenarios, that did not pass is unverified, and the report fails on any.
// Each scenario is shown with who ruled it, and the active scenarios are
// counted by it: the owner's, a delegate's, an agent's provisional. A
// markdown rendering opens with that count, lists the unverified ones
// first, then every entry. The markers TypeScript test runners emit are
// not yet joined.
var CH_042 = &charter.Spec{ID: "CH-042", Title: "A report joins each test's markers to its scenarios and says which passed, failed, were skipped or have none", Status: charter.Active}

// CH_261005_6ZBS is CH-261005-6ZBS: The digest reads the charter in order: needs by area, an area map, then each area's decisions with their why
//
// The digest is one generated markdown file. It lists each active need
// under the area whose active decisions most often reach it through what
// they serve, the area earlier in the reading order on a tie, and the
// needs no active decision reaches apart. A map gives each area in the
// reading order with the first sentence of its header and its count of
// active and planned decisions, then, for each need, the areas that serve
// it with how many decisions. Then, area by area, each active decision in
// the order of its document: its id linked to its heading, its title, its
// Kind, and the first sentence of its Why; planned and deprecated ones are
// counted, not listed. The reading order is a file listing the area
// documents one a line; one it leaves out, or a line that names none or
// names one twice, fails the gate. Generation writes the digest beside the
// generated symbols, and the gate fails a digest that is missing or differs
// from the documents, naming the command that writes it; the gate writes
// it instead when CHARTER_UPDATE_DIGEST is set.
var CH_261005_6ZBS = &charter.Spec{ID: "CH-261005-6ZBS", Title: "The digest reads the charter in order: needs by area, an area map, then each area's decisions with their why", Status: charter.Active}

// CH_261005_VFJR is CH-261005-VFJR: mutate breaks an entry's realizing code and says, per scenario, whether the tests that cite it noticed
//
// mutate takes active entries by id, or the active entries whose realizing
// Go functions a diff since a revision touched, and makes small changes
// inside those functions alone: a comparison or a logical operator flipped,
// a condition negated or its branch emptied, a call dropped, an error
// returned as nil, a boolean flipped, an integer moved by one, a break
// swapped with a continue; a diff-scoped run changes only the lines the
// diff changed. A mutant a caller wrote, an exact replacement in a file,
// runs beside them. Each mutant runs only the tests that cite its entry,
// through an overlay so no file of the tree is written: it is killed when
// one of them fails or times out, survives when all pass, and is counted
// apart when it does not build. The citing tests first run unmutated, and
// an entry whose tests fail then is reported unsteady and none of its
// mutants runs. The result gives, per entry, the mutants killed, survived
// and unbuilt, the survivors with their place and change, and the
// realizing anchors it could not mutate; per scenario, the mutants its
// tests killed, a scenario whose tests killed none being hollow. A run
// after which the tree holds a file it did not hold before fails, naming
// the file; a survivor never fails it.
var CH_261005_VFJR = &charter.Spec{ID: "CH-261005-VFJR", Title: "mutate breaks an entry's realizing code and says, per scenario, whether the tests that cite it noticed", Status: charter.Active}

// CH_261009_3FHW is CH-261009-3FHW: A reader tells what was ruled from what an agent chose, and an entry's change is recorded or named
//
// Every scenario says who ruled it: the owner, with what the ruling cites;
// a named delegate; or an agent whose choice nobody has ruled on, which a
// scenario with nothing written is. The registry and the report carry it,
// and the report counts the active scenarios by who ruled them. A change to
// the record that its Changes lines do not say, an entry removed, unbuilt
// or retired without a line, text or a scenario changed without one, or a
// ruling lost, is named against the revision it was made from.
var CH_261009_3FHW = &charter.Spec{ID: "CH-261009-3FHW", Title: "A reader tells what was ruled from what an agent chose, and an entry's change is recorded or named", Status: charter.Active}

// CH_261009_4B3B is CH-261009-4B3B: The question file holds what nobody ruled, each question naming its entry, who asked and the provisional answer until answered
//
// questions.md in the entries directory holds one question per heading,
// "## Q-ID question", with field lines Asked-by, Entry (an id, or an id and
// a scenario as ID/handle), Provisional (the answer taken meanwhile) and,
// once a person has answered, Answer, opening with the date and who
// answered. A question without Asked-by, Entry or Provisional, naming an
// entry or scenario the record does not hold, with an Answer that does not
// open with a date, or with an id declared twice, is a form problem naming
// the line. A question is never deleted. Lint prints one line, at the
// question file's first line, saying how many questions are open and naming
// the command that lists them; with none open it prints none. claimsgen
// questions lists each open question with its id, the entry it is about,
// its text, the provisional answer and who asked, then how many are open;
// with -all it lists the answered ones too, each with its answer.
var CH_261009_4B3B = &charter.Spec{ID: "CH-261009-4B3B", Title: "The question file holds what nobody ruled, each question naming its entry, who asked and the provisional answer until answered", Status: charter.Active}

// CH_261009_QQC8 is CH-261009-QQC8: claimsgen diff names what a change did to the record that its Changes lines do not say
//
// Given a git revision, claimsgen diff reads the area documents as they
// were there and as they are, and names, by id: an entry removed; an active
// entry made planned; an entry deprecated, or its text, Kind or a scenario
// changed or removed, with no Changes line the revision lacked; and a
// scenario whose ruling was the owner's or a delegate's and is now an
// agent's. An entry only set active, or changed with a new dated line, is
// not named. It exits 1 when it names anything.
var CH_261009_QQC8 = &charter.Spec{ID: "CH-261009-QQC8", Title: "claimsgen diff names what a change did to the record that its Changes lines do not say", Status: charter.Active}

// CH_261009_SNX2 is CH-261009-SNX2: lint -baseline fails when a style check's count exceeds its count in the style baseline, naming both counts
//
// lint -baseline reads style-baseline.txt in the entries directory, a line
// per style check holding its name and a whole count, blank lines and lines
// opening with # skipped. The checks are retired-words,
// titles-over-twenty-words, texts-over-six-sentences,
// scenario-sentences-over-thirty-words, pointing-whys and shared-tests, and
// a check the file leaves out has a baseline of zero. It runs the six style
// checks and, for each whose count exceeds its baseline, prints one line at
// the check's line in the file, or at its first line for a check it leaves
// out, naming both counts and what to do, and exits non-zero; when none grew
// it prints nothing. A line naming an unknown check or one counted already,
// or whose count is not a whole number, is a lint line of its own, and a
// missing file is an error saying there is no baseline. The gate's problems
// are not part of it, and -baseline with -style is refused.
var CH_261009_SNX2 = &charter.Spec{ID: "CH-261009-SNX2", Title: "lint -baseline fails when a style check's count exceeds its count in the style baseline, naming both counts", Status: charter.Active, Kind: "bound"}

// CH_261009_VY5C is CH-261009-VY5C: A scenario says who ruled it: the owner with what it cites, a named delegate, or, unwritten, an agent, provisional
//
// A scenario's handle may be followed, in parentheses, by who ruled it:
// "owner", with what the ruling cites after a comma; "delegate" with the
// delegate's name and, after a comma, what it cites; or "agent,
// provisional". Nothing written is an agent's, provisional, and the
// registry says it was not written. Anything else is a form problem naming
// the line and the three forms. The provenance is part of the scenario as
// the registry gives it.
var CH_261009_VY5C = &charter.Spec{ID: "CH-261009-VY5C", Title: "A scenario says who ruled it: the owner with what it cites, a named delegate, or, unwritten, an agent, provisional", Status: charter.Active}

// CH_261015_9ZTY is CH-261015-9ZTY: Either charter exports as one file, which the viewer shows as linked pages in the charter's reading order
//
// The export writes either charter to one file: each area document in the
// reading order with its header, each entry with its fields, text, Why,
// scenarios and who ruled each, Guidance and Changes, its anchors and its
// citations by handle, the waivers, the observations and the gate's
// problems. A charter the gate fails on is still exported, its problems
// printed apart. The viewer reads that file, or one dropped on the page,
// and names no charter's areas or ids, so the charter's own charter reads
// as the product's does. Its pages walk the areas in the reading order,
// place each need under its area as the digest does, show each area
// document with coverage beside each entry and each entry with what it
// serves and what serves it, and search the record, an exact id opening its
// entry; each page is a link that opens it again. A scenario shows tested
// when a test cites it, waived when a waiver accepts it and untested
// otherwise, and passed or failed only from a report beside the file.
//
// Planned: CH-261015-9ZTY is ruled and not built; see claims/tooling.md.
var CH_261015_9ZTY = &charter.Spec{ID: "CH-261015-9ZTY", Title: "Either charter exports as one file, which the viewer shows as linked pages in the charter's reading order", Status: charter.Planned}

// CH_261126_NS7T is CH-261126-NS7T: An adopter keeps a charter in their own repository with the tooling as a dependency, nothing of it copied
//
// A repository takes the charter up by requiring this module: it writes its
// entries in its own claims directory, its gate as a test that passes its
// options to the gate, and a claimsgen of its own whose main passes the same
// options to the command line. Nothing of the tooling is copied into it,
// and the example adopter in this repository is kept so, its gate run with
// the tooling's tests.
var CH_261126_NS7T = &charter.Spec{ID: "CH-261126-NS7T", Title: "An adopter keeps a charter in their own repository with the tooling as a dependency, nothing of it copied", Status: charter.Active}
