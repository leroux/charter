#!/usr/bin/env python3
"""Write ../questions.jsonl and reference.json from QUESTIONS below.

Run once by the benchmark's author. The reference answers come from
Python's tomllib (the standard library's TOML 1.0 decoder, Python 3.11+),
frozen into reference.json so grading never depends on the Python that
runs it; EXPECT holds the answer the author worked out by hand from the
specification, and qgen refuses to write anything when the two disagree.
"""

import json
import sys
import tomllib
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import tomltest as T  # noqa: E402

VALUE = ("What does toml-decoder output for {where}? Answer `TYPE VALUE`, TYPE as toml-test tags it "
         "(string, integer, float, bool, datetime, datetime-local, date-local, time-local), "
         "e.g. `integer 42`; or `invalid` if toml-decoder refuses the document.")
VALID = "Does toml-decoder (in its default mode) accept this document? Answer `valid` or `invalid`."

# id, document, path to the value (None: ask whether the document is valid),
# the hand-worked answer, the feature it is about.
QUESTIONS = [
    ("q01", 'a = 1\nb = "two"\na = 3\n', ["a"], "invalid", "subset"),
    ("q02", '"a.b" = 1\na.b = 2\n', ["a", "b"], "integer 2", "subset"),
    ("q03", "s = 'C:\\new\\table'\n", ["s"], "string C:\\new\\table", "subset"),
    ("q04", "n = -0\n", ["n"], "integer 0", "subset"),
    ("q05", "[a.b]\nx = 1\n[a]\ny = 2\n", None, "valid", "subset"),
    ("q06", "a.b = 1\n[a]\nc = 2\n", None, "invalid", "subset"),
    ("q07", "port = 0080\n", None, "invalid", "subset"),
    ("q08", 'k = "\\x41"\n', None, "invalid", "subset"),
    ("q09", 't = True\n', None, "invalid", "subset"),
    ("q10", 's = "a#b" # c\n', ["s"], "string a#b", "subset"),
    ("q11", "x = 1e3\n", ["x"], "float 1000", "float"),
    ("q12", "y = .5\n", None, "invalid", "float"),
    ("q13", "h = 0xff\n", ["h"], "integer 255", "int-bases"),
    ("q14", 's = """\nab\\\n   cd"""\n', ["s"], "string abcd", "multiline"),
    ("q15", 'a = [1, "x", [2]]\n', ["a", 1], "string x", "array"),
    ("q16", "p = { x = 1 }\np.y = 2\n", None, "invalid", "inline-table"),
    ("q17", "[[f]]\nn = 1\n[[f]]\nn = 2\n", ["f", 1, "n"], "integer 2", "array-of-tables"),
    ("q18", "d = 1979-05-27 07:32:00Z\n", ["d"], "datetime 1979-05-27T07:32:00Z", "datetime"),
    ("q19", "d = 1979-05-27\n", ["d"], "date-local 1979-05-27", "date-time-local"),
    ("q20", "t = { a = 1, }\n", None, "invalid", "inline-table"),
]


def where(path):
    out = ""
    for p in path:
        out += f"[{p}]" if isinstance(p, int) else (f".{p}" if out else p)
    return "the value at `" + out + "`"


def reference(doc, path):
    try:
        v = T.tag_python(tomllib.loads(doc))
    except tomllib.TOMLDecodeError:
        return {"valid": False}
    if path is None:
        return {"valid": True}
    for p in path:
        v = v[p]
    return {"valid": True, "type": v["type"], "value": v["value"]}


def same(a, b):
    if " " not in a or " " not in b:
        return a == b
    (ta, va), (tb, vb) = a.split(" ", 1), b.split(" ", 1)
    return T.compare({"type": ta, "value": va}, {"type": tb, "value": vb}) is None


def main():
    here = Path(__file__).resolve().parent
    qs, refs = [], {}
    for qid, doc, path, expect, feature in QUESTIONS:
        ref = reference(doc, path)
        ans = ("valid" if ref["valid"] else "invalid") if path is None or not ref["valid"] else \
            f"{ref['type']} {ref['value']}"
        if not same(ans, expect):
            sys.exit(f"{qid}: tomllib says {ans!r}, the author {expect!r}")
        text = (f"Given this TOML document:\n```toml\n{doc}```\n" +
                (VALID if path is None else VALUE.format(where=where(path))))
        qs.append({"id": qid, "question": text})
        refs[qid] = {"doc": doc, "path": path, "feature": feature, "reference": ref, "answer": expect}
    (here.parent / "questions.jsonl").write_text("".join(json.dumps(q) + "\n" for q in qs))
    (here / "reference.json").write_text(json.dumps(refs, indent=1) + "\n")
    print(len(qs), "questions")


if __name__ == "__main__":
    main()
