"""Running toml-test cases against a decoder command, without the network.

The cases are the files of the toml-test submodule (bench/vendor/toml-test);
the comparison follows toml-test's own (json.go): tables and arrays
structurally, floats numerically (NaN by its unsigned spelling), date-times
as instants to the millisecond (all toml-test requires) after toml-test's
replacements, everything else as strings.
"""

import datetime as dt
import json
import math
import os
import re
import shutil
import subprocess
import tempfile
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path

ORACLE = Path(__file__).resolve().parent
BENCH = ORACLE.parent
TOMLTEST = BENCH.parent / "vendor" / "toml-test"
TESTS = TOMLTEST / "tests"


def case_list(version):
    """The .toml cases toml-test lists for a TOML version: valid/... and invalid/..."""
    lines = (TESTS / f"files-toml-{version}").read_text().split()
    return [l for l in lines if l.endswith(".toml")]


def build(ws, name="toml-decoder"):
    """Build the workspace's cmd/<name> into a temp dir. Returns (binary, error)."""
    out = Path(tempfile.mkdtemp(prefix="toml-oracle-")) / name
    env = dict(os.environ, GOPROXY="off", GOFLAGS="-mod=mod", GOWORK="off")
    try:
        p = subprocess.run(["go", "build", "-o", str(out), f"./cmd/{name}"], cwd=ws,
                           capture_output=True, text=True, timeout=300, env=env)
    except subprocess.TimeoutExpired:
        return None, "go build timed out"
    if p.returncode != 0:
        return None, (p.stdout + p.stderr)[-2000:]
    return out, None


def cleanup(binary):
    if binary:
        shutil.rmtree(Path(binary).parent, ignore_errors=True)


def run_decoder(binary, data, args=(), timeout=10):
    try:
        p = subprocess.run([str(binary), *args], input=data, capture_output=True, timeout=timeout)
    except subprocess.TimeoutExpired:
        return None, b"", b"timeout"
    return p.returncode, p.stdout, p.stderr


_REPL = str.maketrans({" ": "T", "t": "T", "z": "Z"})


def _instant(kind, s):
    s = s.translate(_REPL)
    if kind == "datetime":
        m = re.fullmatch(r"(\d{4})-(\d\d)-(\d\d)T(\d\d):(\d\d):(\d\d)(?:\.(\d+))?(Z|[+-]\d\d:\d\d)", s)
        if not m:
            raise ValueError(s)
        y, mo, d, h, mi, se, frac, off = m.groups()
        tz = dt.timezone.utc if off == "Z" else dt.timezone(
            (1 if off[0] == "+" else -1) * dt.timedelta(hours=int(off[1:3]), minutes=int(off[4:6])))
        ns = int((frac or "0").ljust(3, "0")[:3])
        t = dt.datetime(int(y), int(mo), int(d), int(h), int(mi), int(se), tzinfo=tz)
        return (t.astimezone(dt.timezone.utc).replace(tzinfo=None), ns)
    if kind == "datetime-local":
        m = re.fullmatch(r"(\d{4})-(\d\d)-(\d\d)T(\d\d):(\d\d):(\d\d)(?:\.(\d+))?", s)
        if not m:
            raise ValueError(s)
        y, mo, d, h, mi, se, frac = m.groups()
        return (dt.datetime(int(y), int(mo), int(d), int(h), int(mi), int(se)), int((frac or "0").ljust(3, "0")[:3]))
    if kind == "date-local":
        if not re.fullmatch(r"\d{4}-\d\d-\d\d", s):
            raise ValueError(s)
        return dt.date.fromisoformat(s)
    if kind == "time-local":
        # toml-test parses with layout 15:04:05, which accepts a fraction.
        m = re.fullmatch(r"(\d\d):(\d\d):(\d\d)(?:\.(\d+))?", s)
        if not m:
            raise ValueError(s)
        h, mi, se, frac = m.groups()
        return (int(h), int(mi), int(se), int((frac or "0").ljust(3, "0")[:3]))
    raise ValueError(kind)


def _float(s):
    s = s.lower()
    if s.lstrip("+-") == "nan":
        return "nan"
    return float(s)


def compare(want, have, key=""):
    """None when have matches want, else a short reason."""
    if isinstance(want, dict) and set(want) == {"type", "value"} and isinstance(want.get("type"), str):
        if not isinstance(have, dict) or not isinstance(have.get("type"), str):
            return f"{key}: want a {want['type']} value, have {json.dumps(have)[:80]}"
        wt, ht = want["type"], have["type"]
        if wt != ht:
            return f"{key}: want type {wt}, have {ht}"
        wv, hv = want["value"], have.get("value")
        if not isinstance(hv, str):
            return f"{key}: value is not a string"
        try:
            if wt == "float":
                a, b = _float(wv), _float(hv)
                ok = a == b
            elif wt in ("datetime", "datetime-local", "date-local", "time-local"):
                ok = _instant(wt, wv) == _instant(wt, hv)
            elif wt == "bool":
                ok = wv.lower() == hv.lower()
            else:
                ok = wv == hv
        except (ValueError, OverflowError):
            return f"{key}: malformed {wt} {hv!r}"
        return None if ok else f"{key}: want {wv!r}, have {hv!r}"
    if isinstance(want, dict):
        if not isinstance(have, dict):
            return f"{key}: want a table, have {type(have).__name__}"
        if set(want) != set(have):
            missing, extra = sorted(set(want) - set(have)), sorted(set(have) - set(want))
            return f"{key}: keys differ: missing {missing[:5]}, extra {extra[:5]}"
        for k in want:
            r = compare(want[k], have[k], f"{key}.{k}" if key else k)
            if r:
                return r
        return None
    if isinstance(want, list):
        if not isinstance(have, list) or len(want) != len(have):
            return f"{key}: want an array of {len(want)}, have {json.dumps(have)[:80]}"
        for i, (a, b) in enumerate(zip(want, have)):
            r = compare(a, b, f"{key}[{i}]")
            if r:
                return r
        return None
    return f"{key}: unexpected expected value"


def check_case(binary, rel, args=()):
    """Run one toml-test case. Returns (passed, reason)."""
    data = (TESTS / rel).read_bytes()
    code, out, err = run_decoder(binary, data, args)
    if code is None:
        return False, "decoder timed out"
    if rel.startswith("invalid/"):
        if code != 0:
            return True, ""
        return False, "accepted an invalid document"
    if code != 0:
        return False, "refused: " + err.decode("utf-8", "replace").strip()[:200]
    try:
        have = json.loads(out)
    except ValueError:
        return False, "stdout is not JSON"
    want = json.loads((TESTS / (rel[:-5] + ".json")).read_text())
    r = compare(want, have)
    return (r is None), (r or "")


def check_cases(binary, rels, args=(), jobs=8):
    with ThreadPoolExecutor(jobs) as ex:
        return dict(zip(rels, ex.map(lambda r: check_case(binary, r, args), rels)))


def tag_python(v):
    """A value from Python's tomllib as toml-test's tagged JSON."""
    import datetime as _dt
    if isinstance(v, bool):
        return {"type": "bool", "value": "true" if v else "false"}
    if isinstance(v, int):
        return {"type": "integer", "value": str(v)}
    if isinstance(v, float):
        if math.isnan(v):
            return {"type": "float", "value": "nan"}
        if math.isinf(v):
            return {"type": "float", "value": "inf" if v > 0 else "-inf"}
        return {"type": "float", "value": repr(v)}
    if isinstance(v, str):
        return {"type": "string", "value": v}
    if isinstance(v, _dt.datetime):
        return {"type": "datetime" if v.tzinfo else "datetime-local", "value": v.isoformat()}
    if isinstance(v, _dt.date):
        return {"type": "date-local", "value": v.isoformat()}
    if isinstance(v, _dt.time):
        return {"type": "time-local", "value": v.isoformat()}
    if isinstance(v, list):
        return [tag_python(e) for e in v]
    if isinstance(v, dict):
        return {k: tag_python(e) for k, e in v.items()}
    raise TypeError(type(v))
