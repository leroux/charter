#!/usr/bin/env python3
"""Run a benchmark's conditions and runs in parallel, then tabulate.

  matrix.py --bench coreutils --conditions A,B,C --runs 3 -j 3
  matrix.py --bench coreutils --table        # tabulate results already written
"""

import argparse
import json
import subprocess
import sys
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path

BENCH = Path(__file__).resolve().parent.parent
RUN = Path(__file__).resolve().parent / "run.py"


def one(bench, cond, run, extra):
    cmd = [sys.executable, str(RUN), "--bench", bench, "--condition", cond, "--run", str(run), *extra]
    p = subprocess.run(cmd, capture_output=True, text=True)
    return cond, run, p.returncode, (p.stdout + p.stderr)[-2000:]


def table(bench):
    rows = []
    for s in sorted((BENCH / "results" / bench).glob("*/*/summary.json")):
        d = json.loads(s.read_text())
        last = d["tasks"][-1]["oracle"] if d["tasks"] else {}
        spec, tasks = last.get("spec", {}), last.get("tasks", {})
        cost = sum(t.get("cost") or 0 for t in d["tasks"])
        qa = d.get("qa", {})
        if d.get("invalid"):
            rows.append((d["condition"], d["run"], "invalid", d["invalid"][:60], "", ""))
            continue
        rows.append((d["condition"], d["run"], f"{spec.get('passed')}/{spec.get('total')}",
                     f"{tasks.get('passed')}/{tasks.get('total')}", f"{qa.get('correct')}/{qa.get('total')}", f"{cost:.2f}"))
    print("| Condition | Run | Spec at end | Tasks at end | Questions | Cost USD |")
    print("|---|---|---|---|---|---|")
    for r in rows:
        print("| " + " | ".join(str(x) for x in r) + " |")


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--bench", required=True)
    ap.add_argument("--conditions", default="A,B,C,D,E")
    ap.add_argument("--runs", type=int, default=3)
    ap.add_argument("-j", type=int, default=3)
    ap.add_argument("--table", action="store_true")
    args, extra = ap.parse_known_args()
    if args.table:
        table(args.bench)
        return
    jobs = [(c, r) for c in args.conditions.split(",") for r in range(1, args.runs + 1)]
    with ThreadPoolExecutor(args.j) as ex:
        for cond, run, code, tail in ex.map(lambda j: one(args.bench, j[0], j[1], extra), jobs):
            print(f"== {cond} run {run}: exit {code}\n{tail}", flush=True)
    table(args.bench)


if __name__ == "__main__":
    main()
