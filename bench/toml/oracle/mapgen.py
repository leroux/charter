#!/usr/bin/env python3
"""Write features.txt: which feature each toml-test case checks.

Run once by the benchmark's author; the oracle reads features.txt, the
mapping written down. A case belongs to the latest feature (in task order)
it needs: a valid case needs every feature its document uses, so it can
pass only once all of them are built; an invalid case is filed under the
feature its directory names, or the latest feature its document uses.
Everything else is the seed's subset, the spec. Cases listed in OVERRIDES
were read by hand.
"""

import json
import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import tomltest as T  # noqa: E402

# feature -> the task that asks for it. 0 is the seed's subset (the spec).
FEATURES = {
    "subset": 0,
    "float": 1,
    "int-bases": 2,
    "multiline": 4,
    "array": 5,
    "inline-table": 6,
    "array-of-tables": 7,
    "datetime": 9,
    "date-time-local": 10,
    "toml-1.1": 11,
}

DIR = {
    "invalid/float": "float",
    "invalid/array": "array",
    "invalid/inline-table": "inline-table",
    "invalid/datetime": "datetime",
    "invalid/local-datetime": "datetime",
    "invalid/local-date": "date-time-local",
    "invalid/local-time": "date-time-local",
    "valid/float": "float",
    "valid/array": "array",
    "valid/inline-table": "inline-table",
    "valid/datetime": "datetime",
}

# Read by hand: the content rules below would file these wrongly. Each is a
# number, a date or a brace inside a key or a string, and the document uses
# only the subset.
OVERRIDES = {c: "subset" for c in [
    "valid/datetime/invalid-date-in-string.toml",
    "valid/key/alphanum.toml",
    "valid/key/numeric-02.toml",
    "valid/key/numeric-04.toml",
    "valid/key/numeric-05.toml",
    "valid/key/special-word.toml",
    "valid/spec-1.0.0/inline-table-3.toml",
    "valid/spec-1.0.0/keys-7.toml",
    "valid/table/keyword.toml",
    "valid/table/keyword-with-values.toml",
]}

NUM = r"(?:^|[=\[,{]\s*)"


def content_features(text, want=None):
    f = set()
    if '"""' in text or "'''" in text:
        f.add("multiline")
    if re.search(r"(?m)^\s*\[\[", text):
        f.add("array-of-tables")
    if re.search(r"=\s*\[", text):
        f.add("array")
    if re.search(r"[=\[,]\s*\{", text):
        f.add("inline-table")
    if re.search(NUM + r"[+-]?(?:\d[\d_]*\.\d|\d[\d_]*[eE][+-]?\d|inf\b|nan\b)", text, re.M):
        f.add("float")
    if re.search(NUM + r"[+-]?0[xob]", text, re.M) or re.search(NUM + r"[+-]?\d+_\d", text, re.M):
        f.add("int-bases")
    if re.search(r"\d{4}-\d\d-\d\d[Tt ]\d\d:\d\d", text):
        f.add("datetime")
    elif re.search(r"\d{4}-\d\d-\d\d", text) or re.search(NUM + r"\d\d:\d\d", text, re.M):
        f.add("date-time-local")
    if want is not None:
        walk(want, f)
    return f


def walk(v, f):
    if isinstance(v, dict) and set(v) == {"type", "value"}:
        f.add({"float": "float", "datetime": "datetime", "datetime-local": "datetime",
               "date-local": "date-time-local", "time-local": "date-time-local"}.get(v["type"], "subset"))
    elif isinstance(v, dict):
        for e in v.values():
            walk(e, f)
    elif isinstance(v, list):
        if v and all(isinstance(e, dict) and not (set(e) == {"type", "value"}) for e in v):
            f.add("array-of-tables")  # or an array of inline tables: the later of the two
        else:
            f.add("array")
        for e in v:
            walk(e, f)


def latest(fs):
    fs = set(fs) - {"subset"}
    return max(fs, key=FEATURES.get) if fs else "subset"


def classify(rel):
    text = (T.TESTS / rel).read_bytes().decode("utf-8", "replace")
    d = "/".join(rel.split("/")[:2])
    want = None
    if rel.startswith("valid/"):
        want = json.loads((T.TESTS / (rel[:-5] + ".json")).read_text())
    fs = content_features(text, want)
    if d in DIR:
        fs.add(DIR[d])
    # Underscores and the prefixes of the other bases, however misspelt.
    if d == "invalid/integer" and re.search(r"us|hex|oct|bin", rel.split("/")[-1]):
        fs.add("int-bases")
    return latest(fs)


def needs_11(rel):
    import tomllib
    want = json.loads((T.TESTS / (rel[:-5] + ".json")).read_text())
    try:
        got = T.tag_python(tomllib.loads((T.TESTS / rel).read_bytes().decode()))
    except Exception:
        return True
    return T.compare(want, got) is not None


def main():
    v10 = T.case_list("1.0.0")
    v11 = [c for c in T.case_list("1.1.0") if c not in set(v10)]
    rows = []
    for c in v10:
        rows.append((c, "1.0", OVERRIDES.get(c, classify(c))))
    for c in v11:
        # The 1.1 spec's examples mostly restate 1.0's: keep a valid one only
        # when a TOML 1.0 decoder (Python's tomllib) cannot read it as
        # expected, and drop the invalid ones, which 1.0 refuses as well.
        if "/spec-1.1.0/" in c and (c.startswith("invalid/") or not needs_11(c)):
            continue
        rows.append((c, "1.1", "toml-1.1"))
    out = Path(__file__).resolve().parent / "features.txt"
    with out.open("w") as f:
        f.write("# toml-test case, the TOML version the decoder runs it in, the feature it checks.\n"
                "# Written by mapgen.py; the feature's task is in mapgen.FEATURES and score.\n")
        for r in rows:
            f.write(" ".join(r) + "\n")
    from collections import Counter
    print(Counter(r[2] for r in rows))


if __name__ == "__main__":
    main()
