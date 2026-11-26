// Package python is the charter's Python language: a reader that runs
// python3 over the standard ast module to turn Python source into facts,
// and a generator that writes claims/claims.py, whose verifies and observes
// are the marker a unittest or pytest test calls.
package python

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/record"
)

//go:embed reader.py
var readerScript string

// Interpreter is the python3 the reader runs.
var Interpreter = "python3"

// ClaimsModule is the name tests import the generated module by.
const ClaimsModule = "claims"

// Reader reads Python files into facts.
type Reader struct {
	Root string
	// Dir is the directory, slash-relative to the root, the Python tests
	// run in: module names are taken from it.
	Dir string
	// EnvFile is the constants module's file, slash-relative to the root,
	// or "".
	EnvFile string
	// ClaimsDir is the claims directory, whose generated file is not read.
	ClaimsDir string
	// Table names the configuration's Python table in a problem.
	Table string
}

// ErrNoInterpreter is wrapped when python3 cannot be run.
var ErrNoInterpreter = errors.New("python3 is not on PATH")

// Read runs the reader script over the files. The script, reader.py,
// embedded here, is where Python's anchors, citations and environment reads
// are decided; this charter reads no Python of its own, so they are anchored
// on the call that runs it.
//
//charter:realizes CH-261126-PQ7X, CH-021, CH-023, CH-034
func (r Reader) Read(files []string) (*facts.Set, error) {
	var keep []string
	for _, f := range files {
		if r.ClaimsDir != "" && strings.HasPrefix(f, r.ClaimsDir+"/") {
			continue
		}
		keep = append(keep, f)
	}
	if len(keep) == 0 {
		return &facts.Set{}, nil
	}
	bin, err := exec.LookPath(Interpreter)
	if err != nil {
		return nil, fmt.Errorf("%s: the charter reads Python and %w; install it, or remove the [python] table", r.Table, ErrNoInterpreter)
	}
	req, err := json.Marshal(map[string]any{
		"root": r.Root, "files": keep, "claims_module": ClaimsModule,
		"env_file": r.EnvFile, "env_module": moduleName(r.Dir, r.EnvFile),
	})
	if err != nil {
		return nil, err
	}
	// -I: the tree's own files are never imported in place of the
	// standard library the script reads with.
	cmd := exec.Command(bin, "-I", "-c", readerScript)
	cmd.Stdin = bytes.NewReader(req)
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("the Python reader failed: %v\n%s", err, strings.TrimSpace(errOut.String()))
	}
	var reply struct {
		Facts *facts.Set `json:"facts"`
		Error string     `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &reply); err != nil {
		return nil, fmt.Errorf("the Python reader wrote something other than JSON: %w: %s", err, out.String())
	}
	if reply.Error != "" {
		return nil, errors.New(reply.Error)
	}
	return reply.Facts, nil
}

// moduleName is the dotted module a file is imported as from dir.
func moduleName(dir, file string) string {
	if file == "" {
		return ""
	}
	rel := strings.TrimPrefix(file, strings.Trim(dir, "/")+"/")
	if dir == "" || dir == "." {
		rel = file
	}
	return strings.ReplaceAll(strings.TrimSuffix(rel, path.Ext(rel)), "/", ".")
}

// Generate renders claims.py: a symbol per entry and per observation, the
// types they hold, and verifies and observes. regen is the command named
// in its header.
//
//charter:realizes CH-020
func Generate(rec *record.Record, regen string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# Code generated from %s by %s. DO NOT EDIT.\n", rec.Dir, regen)
	b.WriteString(marker)
	for _, e := range rec.Entries {
		fmt.Fprintf(&b, "\n#: %s is %s: %s\n", e.Symbol, e.ID, oneLine(e.Title))
		if e.Status == "deprecated" {
			fmt.Fprintf(&b, "#: Deprecated: %s is retired; see its Changes in %s. Add no new anchors.\n", e.ID, e.File)
		}
		fmt.Fprintf(&b, "%s = Spec(%s, %s, %q, %s)\n", e.Symbol, pyString(e.ID), pyString(e.Title), e.Status, pyString(e.Kind))
	}
	for _, o := range rec.Observations {
		fmt.Fprintf(&b, "\n#: %s is an observation of %s: %s\n%s = Observation(%s, %s, %s)\n", o.Name, o.Subject, oneLine(o.Text), o.Name, pyString(o.Name), pyString(o.Subject), pyString(o.Status))
	}
	return b.Bytes()
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// pyString is s as a Python string literal: JSON's escapes are Python's.
func pyString(s string) string {
	q, _ := json.Marshal(s)
	return string(q)
}

// marker is the generated module's own code.
//
//charter:realizes CH-024
const marker = `"""The charter's record as Python: a symbol per entry and per observation,
and the verifies and observes a test opens with. With CHARTER_EVENTS set,
each citation and its test's outcome is appended to that file as a JSON line
once the test ends. Under pytest, list this module in pytest_plugins so its
hook sees each test's outcome."""

import json
import os
import sys
import threading

__all__ = ["Spec", "Observation", "verifies", "observes"]


class Spec:
    """An entry of the record: a need or a decision."""

    __slots__ = ("id", "title", "status", "kind")

    def __init__(self, id, title, status, kind):
        object.__setattr__(self, "id", id)
        object.__setattr__(self, "title", title)
        object.__setattr__(self, "status", status)
        object.__setattr__(self, "kind", kind)

    def __setattr__(self, name, value):
        raise AttributeError("an entry is changed in its document, not by a test")


class Observation:
    """A fact about a system the code does not own."""

    __slots__ = ("name", "subject", "status")

    def __init__(self, name, subject, status):
        object.__setattr__(self, "name", name)
        object.__setattr__(self, "subject", subject)
        object.__setattr__(self, "status", status)

    def __setattr__(self, name, value):
        raise AttributeError("an observation is changed in its document, not by a test")


_lock = threading.Lock()
_pending = {}  # pytest node id -> [(id, handle, file)]


def _write(id, handles, test, file, outcome):
    path = os.environ.get("CHARTER_EVENTS", "")
    if not path:
        return
    with _lock, open(path, "a", encoding="utf-8") as f:
        for h in handles or [""]:
            line = {"id": id, "test": test, "file": file, "outcome": outcome}
            if h:
                line["handle"] = h
            # One write per line, opened to append: a line is never split by
            # another process's.
            f.write(json.dumps(line) + "\n")


def _unittest_outcome(tc):
    outcome = getattr(tc, "_outcome", None)
    result = getattr(outcome, "result", None)
    if result is not None:
        if any(t is tc for t, _ in getattr(result, "skipped", [])):
            return "skipped"
        if any(t is tc for t, _ in list(result.failures) + list(result.errors)):
            return "failed"
    if outcome is not None and not getattr(outcome, "success", True):
        return "failed"
    return "passed"


def _record(tc, id, handles):
    file = sys._getframe(2).f_code.co_filename
    if not os.environ.get("CHARTER_EVENTS", ""):
        return
    if tc is None:
        node = os.environ.get("PYTEST_CURRENT_TEST", "").rsplit(" ", 1)[0]
        _pending.setdefault(node, []).append((id, handles, file))
        return
    tc.addCleanup(lambda: _write(id, handles, tc.id(), file, _unittest_outcome(tc)))


def verifies(tc, spec, *handles):
    """Opens a test that verifies an entry, naming by handle the scenarios it
    exercises: tc is the unittest.TestCase, or None under pytest."""
    if not isinstance(spec, Spec):
        raise TypeError("claims: verifies needs a spec")
    _record(tc, spec.id, list(handles))


def observes(tc, observation):
    """Opens a probe that checks an observation against the live system."""
    if not isinstance(observation, Observation):
        raise TypeError("claims: observes needs an observation")
    _record(tc, observation.name, [])


def pytest_runtest_logreport(report):
    """The pytest hook: a test's call report, or a failed or skipped setup,
    gives the outcome of the citations it made."""
    if report.when == "call" or (report.when == "setup" and not report.passed):
        outcome = "passed" if report.passed else "skipped" if report.skipped else "failed"
        for id, handles, file in _pending.pop(report.nodeid, []):
            _write(id, handles, report.nodeid, file, outcome)
`
