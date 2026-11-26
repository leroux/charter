// Coverage as the gate's rule 25 reads it: an active decision is covered
// when code realizes it and a test verifies each scenario, or a waiver
// accepts the missing test; a need is realized by the decisions that serve
// it, so it asks for tests only; a spec without scenarios is verified by any
// citation. A report, when one is loaded, says whether the citing Go tests
// passed. Pure.

import { isDecision, type Entry, type Ref, type Registry, type Waiver } from "./model";
import type { Report, TestOutcome } from "./report";

// A mark is what is known of one scenario, or of a spec without scenarios.
//   pass, fail, skip: the report ran its Go tests
//   cited: a test names it, and no report says how it went
//   waived: no test names it, and a waiver line accepts that
//   missing: active, and nothing verifies it
//   open: not active, and no test names it yet
export type Mark = "pass" | "fail" | "skip" | "cited" | "waived" | "missing" | "open";

export interface Check {
  // ID/handle for a scenario, ID for a spec without scenarios.
  key: string;
  handle: string;
  mark: Mark;
  refs: Ref[];
  tests: TestOutcome[];
  waiver: Waiver | null;
}

// Verdict is an entry's coverage at a glance: ok, waived (covered with a
// waiver), gap (unrealized or a scenario unverified), fail (a test failed),
// or none, for an entry coverage does not apply to.
export type Verdict = "ok" | "waived" | "gap" | "fail" | "none";

export interface Coverage {
  // null for a need or an observation, which no anchor realizes.
  realized: boolean | null;
  checks: Check[];
  // Tests that cite a spec with scenarios without naming one.
  bare: Ref[];
  verdict: Verdict;
}

// waiverFor is the waiver line naming the spec or scenario key; the gate
// refuses a glob, so a waiver names its key exactly.
export function waiverFor(reg: Registry, key: string): Waiver | null {
  return reg.waivers.find((w) => w.pattern === key) ?? null;
}

function mark(active: boolean, refs: Ref[], ran: Mark | null, waiver: Waiver | null): Mark {
  if (ran === "pass" || ran === "fail") return ran;
  if (refs.length > 0) return ran === "skip" ? "skip" : "cited";
  if (waiver) return "waived";
  return active ? "missing" : "open";
}

export function coverage(reg: Registry, e: Entry, report: Report | null): Coverage {
  if (e.level !== "spec") return { realized: null, checks: [], bare: [], verdict: "none" };
  const active = e.status === "active";
  const ran = (key: string): { mark: Mark | null; tests: TestOutcome[] } => {
    const r = report?.results.get(key);
    if (!r) return { mark: null, tests: [] };
    const m = r.outcome === "pass" || r.outcome === "fail" || r.outcome === "skip" ? r.outcome : null;
    return { mark: m, tests: r.tests };
  };
  const checks: Check[] = [];
  let bare: Ref[] = [];
  if (e.scenarios.length === 0) {
    const r = ran(e.id);
    const waiver = waiverFor(reg, e.id);
    checks.push({ key: e.id, handle: "", mark: mark(active, e.refs, r.mark, waiver), refs: e.refs, tests: r.tests, waiver });
  } else {
    for (const s of e.scenarios) {
      const key = `${e.id}/${s.handle}`;
      const refs = e.refs.filter((x) => x.handles.includes(s.handle));
      const r = ran(key);
      const waiver = waiverFor(reg, key);
      checks.push({ key, handle: s.handle, mark: mark(active, refs, r.mark, waiver), refs, tests: r.tests, waiver });
    }
    bare = e.refs.filter((x) => x.handles.length === 0);
  }
  const realized = isDecision(e) ? e.anchors.some((a) => a.kind === "realizes") : null;
  let verdict: Verdict = "none";
  if (active) {
    const marks = new Set(checks.map((c) => c.mark));
    if (marks.has("fail")) verdict = "fail";
    else if (realized === false || marks.has("missing")) verdict = "gap";
    else if (marks.has("waived")) verdict = "waived";
    else verdict = "ok";
  }
  return { realized, checks, bare, verdict };
}

// Tally counts what a set of entries holds and how covered it is: entries
// by status, and over the active specs their decisions realized and their
// checks by mark.
export interface Tally {
  entries: number;
  active: number;
  planned: number;
  deprecated: number;
  decisions: number;
  realized: number;
  marks: Record<Mark, number>;
  verdicts: Record<Verdict, number>;
}

export function emptyTally(): Tally {
  return {
    entries: 0,
    active: 0,
    planned: 0,
    deprecated: 0,
    decisions: 0,
    realized: 0,
    marks: { pass: 0, fail: 0, skip: 0, cited: 0, waived: 0, missing: 0, open: 0 },
    verdicts: { ok: 0, waived: 0, gap: 0, fail: 0, none: 0 },
  };
}

export function tally(reg: Registry, entries: Entry[], report: Report | null): Tally {
  const t = emptyTally();
  for (const e of entries) {
    if (e.level !== "spec") continue;
    t.entries++;
    if (e.status === "active") t.active++;
    else if (e.status === "planned") t.planned++;
    else t.deprecated++;
    if (e.status !== "active") continue;
    const c = coverage(reg, e, report);
    t.verdicts[c.verdict]++;
    if (c.realized !== null) {
      t.decisions++;
      if (c.realized) t.realized++;
    }
    for (const k of c.checks) t.marks[k.mark]++;
  }
  return t;
}

export const MARK_LABEL: Record<Mark, string> = {
  pass: "passed",
  fail: "failed",
  skip: "skipped",
  cited: "tested",
  waived: "waived",
  missing: "untested",
  open: "not yet tested",
};

export const VERDICT_LABEL: Record<Verdict, string> = {
  ok: "covered",
  waived: "covered with waivers",
  gap: "gap",
  fail: "failing",
  none: "",
};
