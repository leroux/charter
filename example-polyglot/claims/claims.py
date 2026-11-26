# Code generated from claims by charter gen. DO NOT EDIT.
"""The charter's record as Python: a symbol per entry and per observation,
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

#: STK_261126_78RN is STK-261126-78RN: A title gives the same slug from the service and from the tool
STK_261126_78RN = Spec("STK-261126-78RN", "A title gives the same slug from the service and from the tool", "active", "")

#: SW_261126_4NQQ is SW-261126-4NQQ: A slug is the title's letters and digits in lower case, its words joined by one separator
SW_261126_4NQQ = Spec("SW-261126-4NQQ", "A slug is the title's letters and digits in lower case, its words joined by one separator", "active", "invariant")

#: SW_261126_X6HB is SW-261126-X6HB: The tool takes another separator from SLUGS_SEPARATOR
SW_261126_X6HB = Spec("SW-261126-X6HB", "The tool takes another separator from SLUGS_SEPARATOR", "active", "")
