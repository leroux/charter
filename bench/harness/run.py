#!/usr/bin/env python3
"""Run one benchmark run: one condition, one run number, a stream of tasks.

The harness copies the benchmark's seed for the condition into a fresh
workspace outside the repository, gives each task to one claude -p session
resumed task after task, and after each task runs the hidden oracle on the
workspace. See bench/README.md for the contract a benchmark follows.
"""

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
import time
import uuid
from pathlib import Path

BENCH = Path(__file__).resolve().parent.parent
CHARTER = BENCH.parent
WORK = Path(os.environ.get("CHARTER_BENCH_WORK", "/tmp/charter-bench"))
CONDITIONS = "ABCDE"
SKILLS = ["claims-implement", "claims-review"]

# The oracles and the hidden parts of the tasks live under bench/. The
# agent's file tools are denied it through --settings, which the workspace
# never shows, and every tool call in the transcript is audited for it.
HIDDEN = str(BENCH)
SETTINGS = json.dumps({"permissions": {"deny": [
    f"{tool}(/{HIDDEN}/**)" for tool in ("Read", "Edit", "Write", "Grep", "Glob")]}})

# What CLAUDE.md says in each condition. A says nothing of a charter; each
# later condition adds to the one before.
CLAUDE_MD = {
    "A": """# Project

Make the changes asked of you. Keep the tests passing: `go test ./...`.
""",
    "B": """# Project

Make the changes asked of you. Keep the tests passing: `go test ./...`.

## The charter

`claims/` records what this project does: needs in `claims/needs.md`,
decisions in the other documents, each with a Why and numbered scenarios.
Read the entries a change touches before changing the code. Keep the
charter true: when a change alters what a decision says, change the
decision and add a dated line under its `Changes:`; when it adds behaviour,
add a decision or a scenario. If a request conflicts with a decision, say
so in your reply and do not silently break the decision.
""",
}
CLAUDE_MD["C"] = CLAUDE_MD["B"] + """
The charter is enforced. `go test ./...` runs the gate in
`claims/gate_test.go`: every active decision has code that carries
`//charter:realizes ID` in its doc comment, and every scenario has a test
whose first statement is `chartertest.Verifies(t, claims.ID, "handle")`.
After editing a document in `claims/`, run `go run ./cmd/claimsgen` to
regenerate `claims/specs.go`, and `go run ./cmd/claimsgen next SW` for a
new id; never write an id by hand. `go run ./cmd/claimsgen lint` says what
the gate does not fail on.
"""
CLAUDE_MD["D"] = CLAUDE_MD["C"] + """
Before you finish a change, run `go run ./cmd/claimsgen report -md` (every
active scenario must pass) and `go run ./cmd/claimsgen diff HEAD` (every
change to the record has its dated Changes line). The skills
claims-implement and claims-review in `.claude/skills/` describe the loop
and the review; follow them.
"""
CLAUDE_MD["E"] = CLAUDE_MD["D"] + """
For the decisions a change touches, run `go run ./cmd/claimsgen mutate -md
ID` and strengthen a test when a scenario's tests miss a mutant that is not
equivalent.
"""


def read_tasks(bench: Path):
    """The tasks in order: number, slug, kind, the text the agent sees."""
    out = []
    for p in sorted((bench / "tasks").glob("*.md")):
        m = re.match(r"(\d+)-(.+)\.md$", p.name)
        if not m:
            continue
        text = p.read_text()
        kind = "neutral"
        first, _, rest = text.partition("\n")
        if first.startswith("kind:"):
            kind, text = first.split(":", 1)[1].strip(), rest
        visible = text.split("---hidden---", 1)[0].strip()
        out.append({"n": int(m.group(1)), "slug": m.group(2), "kind": kind, "prompt": visible})
    return out


def pick(tasks, spec):
    if not spec:
        return tasks
    lo, _, hi = spec.partition("-")
    lo, hi = int(lo), int(hi or lo)
    return [t for t in tasks if lo <= t["n"] <= hi]


def sh(cmd, cwd, timeout=None, env=None):
    return subprocess.run(cmd, cwd=cwd, capture_output=True, text=True, timeout=timeout, env=env)


def prepare(bench: Path, cond: str, ws: Path):
    """Copy the condition's seed into a fresh workspace and commit it."""
    if ws.exists():
        shutil.rmtree(ws)
    seed = bench / "seed" / ("plain" if cond == "A" else "chartered")
    shutil.copytree(seed, ws, symlinks=True)
    if cond == "B":
        # The markdown alone: no gate, no tooling, no anchors or citations.
        strip_tooling(ws)
    gomod = ws / "go.mod"
    if gomod.exists():
        gomod.write_text(gomod.read_text().replace("@CHARTER@", str(CHARTER)))
    (ws / "CLAUDE.md").write_text(CLAUDE_MD[cond])
    if cond in "DE":
        for s in SKILLS:
            shutil.copytree(CHARTER / ".claude" / "skills" / s, ws / ".claude" / "skills" / s)
    sh(["git", "init", "-q", "-b", "main"], ws)
    sh(["git", "add", "-A"], ws)
    sh(["git", "-c", "user.name=bench", "-c", "user.email=bench@localhost", "commit", "-qm", "seed"], ws)


def strip_tooling(ws: Path):
    """Turn a chartered tree into condition B's: the area documents stay,
    the gate, the generated files, the commands, the anchors and the
    citations go, and so does the charter module requirement."""
    claims = ws / "claims"
    for p in list(claims.rglob("*")):
        generated = p.name == "DIGEST.md"  # generated by the tooling, which B lacks
        if p.is_file() and (generated or (p.suffix != ".md" and p.name != "reading-order.txt")):
            p.unlink()
    for d in sorted([p for p in claims.rglob("*") if p.is_dir()], reverse=True):
        if not any(d.iterdir()):
            d.rmdir()
    shutil.rmtree(ws / "cmd" / "claimsgen", ignore_errors=True)
    for p in ws.rglob("*.go"):
        src = p.read_text()
        new = re.sub(r"^//\s*\n//charter:(realizes|concerns)[^\n]*\n", "", src, flags=re.M)
        new = re.sub(r"^//charter:(realizes|concerns)[^\n]*\n", "", new, flags=re.M)
        new = re.sub(r"^\s*chartertest\.(Verifies|Observes)\([^\n]*\)\n", "", new, flags=re.M)
        if new != src:
            new = re.sub(r'^\s*"[^"\n]*/claims"\n', "", new, flags=re.M)
            new = re.sub(r'^\s*"github\.com/leroux/charter/chartertest"\n', "", new, flags=re.M)
            p.write_text(new)
    gomod = ws / "go.mod"
    if gomod.exists():
        lines = [l for l in gomod.read_text().splitlines() if "leroux/charter" not in l]
        gomod.write_text("\n".join(lines) + "\n")
    sh(["gofmt", "-w", "."], ws)
    if shutil.which("goimports"):
        sh(["goimports", "-w", "."], ws)


def agent(prompt, ws: Path, session: str, first: bool, model: str, dry: bool, max_turns: int):
    """One task to the agent. Returns the parsed JSON result."""
    if dry:
        return {"dry_run": True, "total_cost_usd": 0, "num_turns": 0, "duration_ms": 0, "result": ""}
    cmd = ["claude", "-p", prompt, "--output-format", "json",
           "--permission-mode", "bypassPermissions", "--settings", SETTINGS,
           "--setting-sources", "project,local",
           "--max-turns", str(max_turns)]
    cmd += ["--session-id", session] if first else ["--resume", session]
    if model:
        cmd += ["--model", model]
    p = subprocess.run(cmd, cwd=ws, capture_output=True, text=True, timeout=3600, stdin=subprocess.DEVNULL)
    try:
        return json.loads(p.stdout)
    except json.JSONDecodeError:
        return {"error": "no JSON from claude", "stdout": p.stdout[-4000:], "stderr": p.stderr[-4000:]}


def save_transcript(session: str, out: Path, name: str):
    """Copy the session's transcript, which claude writes under
    ~/.claude/projects, into the results, so the trail outlives the
    container."""
    for t in (Path.home() / ".claude" / "projects").glob(f"*/{session}.jsonl"):
        shutil.copyfile(t, out / name)
        return


def touched_hidden(out: Path, name: str):
    """The tool calls in a saved transcript that name the benchmark's
    hidden directory: a run whose agent looked there measures nothing."""
    t = out / name
    if not t.exists():
        return []
    hits = []
    for line in t.read_text().splitlines():
        try:
            ev = json.loads(line)
        except json.JSONDecodeError:
            continue
        content = (ev.get("message") or {}).get("content")
        if not isinstance(content, list):
            continue
        for c in content:
            if isinstance(c, dict) and c.get("type") == "tool_use":
                blob = json.dumps(c.get("input", {}))
                if HIDDEN in blob or "/charter/bench" in blob:
                    hits.append(blob[:300])
    return hits


def save_history(ws: Path, out: Path):
    """The workspace's git history, a commit per task, as a bundle:
    `git clone workspace.bundle` restores every state the agent left."""
    sh(["git", "bundle", "create", str(out / "workspace.bundle"), "--all"], ws)


def score(bench: Path, ws: Path, n: int):
    p = sh([str(bench / "oracle" / "score"), str(ws), str(n)], bench, timeout=600)
    try:
        return json.loads(p.stdout)
    except json.JSONDecodeError:
        return {"error": "oracle printed no JSON", "stdout": p.stdout[-2000:], "stderr": p.stderr[-2000:]}


def gate(ws: Path, cond: str):
    """Whether the workspace's tests, and with them the gate, pass."""
    p = sh(["go", "test", "-count=1", "./..."], ws, timeout=900)
    out = {"tests_pass": p.returncode == 0, "tail": (p.stdout + p.stderr)[-3000:]}
    if cond in "CDE":
        g = sh(["go", "test", "-count=1", "./claims"], ws, timeout=600)
        out["gate_pass"] = g.returncode == 0
    return out


def diffstat(ws: Path):
    p = sh(["git", "diff", "--shortstat", "HEAD"], ws)
    return p.stdout.strip()


def commit(ws: Path, msg: str):
    sh(["git", "add", "-A"], ws)
    sh(["git", "-c", "user.name=bench", "-c", "user.email=bench@localhost", "commit", "-qm", msg, "--allow-empty"], ws)


def questions(bench: Path, ws: Path, out: Path, model: str, dry: bool):
    """A fresh session answers every question; the oracle grades them."""
    qs = [json.loads(l) for l in (bench / "questions.jsonl").read_text().splitlines() if l.strip()]
    prompt = ("Answer each question about this codebase by reading it. Do not run, build or modify "
              "anything. Reply with JSON lines only, one per question: {\"id\": ..., \"answer\": ...}, "
              "each answer one line.\n\n" + "\n".join(json.dumps(q) for q in qs))
    qsession = str(uuid.uuid4())
    res = agent(prompt, ws, qsession, True, model, dry, 40) if not dry else {"result": "", "dry_run": True}
    if not dry:
        save_transcript(qsession, out, "qa-transcript.jsonl")
    answers = []
    for line in str(res.get("result", "")).splitlines():
        line = line.strip().strip("`")
        if line.startswith("{"):
            try:
                answers.append(json.loads(line))
            except json.JSONDecodeError:
                pass
    af = out / "answers.jsonl"
    af.write_text("\n".join(json.dumps(a) for a in answers) + "\n")
    p = sh([str(bench / "oracle" / "grade"), str(ws), str(af)], bench, timeout=600)
    try:
        grade = json.loads(p.stdout)
    except json.JSONDecodeError:
        grade = {"error": "grader printed no JSON", "stdout": p.stdout[-2000:], "stderr": p.stderr[-2000:]}
    hits = touched_hidden(out, "qa-transcript.jsonl") if not dry else []
    if hits:
        grade["invalid"] = "the agent reached for the hidden oracle: " + hits[0]
    return {"agent": usage(res), "grade": grade}


LIMIT_MARKERS = ("hit your", "usage limit", "rate limit", "rate_limit", "overloaded")


def failed(res):
    """Why the agent's turn did not count, or None: an API error, a usage
    or rate limit, or no JSON. Such a task is no measurement."""
    if res.get("dry_run"):
        return None
    if "error" in res:
        return res["error"]
    reply = str(res.get("result", "")).lower()
    if res.get("is_error") or any(m in reply for m in LIMIT_MARKERS) and len(reply) < 300:
        return "agent error: " + str(res.get("result", ""))[:200]
    return None


def usage(res):
    keys = ["total_cost_usd", "num_turns", "duration_ms", "is_error", "dry_run", "error", "session_id"]
    u = {k: res[k] for k in keys if k in res}
    if "result" in res:
        u["reply"] = str(res["result"])[-3000:]
    return u


def run(args):
    bench = BENCH / args.bench
    cond = args.condition
    if cond not in CONDITIONS:
        sys.exit(f"condition is one of {CONDITIONS}")
    out = BENCH / "results" / args.bench / cond / str(args.run)
    out.mkdir(parents=True, exist_ok=True)
    ws = WORK / f"{args.bench}-{cond}-{args.run}"
    prepare(bench, cond, ws)
    session = str(uuid.uuid4())
    summary = {"bench": args.bench, "condition": cond, "run": args.run, "session": session,
               "model": args.model, "started": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()), "tasks": []}
    for i, t in enumerate(pick(read_tasks(bench), args.tasks)):
        t0 = time.time()
        res = agent(t["prompt"], ws, session, i == 0, args.model, args.dry_run, args.max_turns)
        why = failed(res)
        record = {"task": t["n"], "slug": t["slug"], "kind": t["kind"], "agent": usage(res), "invalid": why,
                  "seconds": round(time.time() - t0, 1), "diff": diffstat(ws),
                  "checks": gate(ws, cond), "oracle": score(bench, ws, t["n"])}
        commit(ws, f"task {t['n']}: {t['slug']}")
        if not args.dry_run:
            save_transcript(session, out, "transcript.jsonl")
            hits = touched_hidden(out, "transcript.jsonl")
            if hits and not why:
                why = record["invalid"] = "the agent reached for the hidden oracle: " + hits[0]
        (out / f"task-{t['n']:02d}.json").write_text(json.dumps(record, indent=2))
        summary["tasks"].append({k: record[k] for k in ("task", "kind", "oracle")} |
                                {"cost": record["agent"].get("total_cost_usd"), "tests_pass": record["checks"]["tests_pass"]})
        if why:
            summary["invalid"] = f"task {t['n']}: {why}"
            print(f"{args.bench} {cond} run {args.run} task {t['n']}: INVALID, stopping: {why}", flush=True)
            break
        print(f"{args.bench} {cond} run {args.run} task {t['n']}: oracle {json.dumps(record['oracle'].get('spec'))} "
              f"{json.dumps(record['oracle'].get('tasks'))} cost {record['agent'].get('total_cost_usd')}", flush=True)
    if not args.no_questions and "invalid" not in summary:
        qa = questions(bench, ws, out, args.model, args.dry_run)
        (out / "qa.json").write_text(json.dumps(qa, indent=2))
        summary["qa"] = qa["grade"]
    save_history(ws, out)
    (out / "summary.json").write_text(json.dumps(summary, indent=2))
    return summary


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--bench", required=True)
    ap.add_argument("--condition", required=True)
    ap.add_argument("--run", type=int, default=1)
    ap.add_argument("--tasks", help="a task number or range, such as 1-3; default every task")
    ap.add_argument("--model", default="")
    ap.add_argument("--max-turns", type=int, default=80)
    ap.add_argument("--dry-run", action="store_true", help="no agent: test the seed, the strip and the oracle")
    ap.add_argument("--no-questions", action="store_true")
    run(ap.parse_args())


if __name__ == "__main__":
    main()
