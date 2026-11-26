// A test run joined to the charter: the JSON `claimsgen report` writes, per
// spec and scenario the citing Go tests and whether they passed. Optional:
// without one the viewer says which scenarios a test cites, not whether it
// passed. Pure.

export type Outcome = "pass" | "fail" | "skip" | "none" | "waived";
const OUTCOMES: Outcome[] = ["pass", "fail", "skip", "none", "waived"];

export interface TestOutcome {
  pkg: string;
  test: string;
  outcome: Outcome;
}

export interface Result {
  outcome: Outcome;
  tests: TestOutcome[];
}

export interface Report {
  // Keyed "ID/handle" for a scenario, "ID" for a spec's own citations.
  results: Map<string, Result>;
}

function isRecord(v: unknown): v is Record<string, unknown> {
  return typeof v === "object" && v !== null && !Array.isArray(v);
}

function outcome(v: unknown, what: string): Outcome {
  if (!OUTCOMES.includes(v as Outcome)) throw new Error(`${what}: unknown outcome ${JSON.stringify(v)}`);
  return v as Outcome;
}

function tests(v: unknown, what: string): TestOutcome[] {
  if (v === undefined || v === null) return [];
  if (!Array.isArray(v)) throw new Error(`${what} must be an array`);
  return v.map((t, i) => {
    if (!isRecord(t)) throw new Error(`${what} ${i} is not an object`);
    return { pkg: String(t.Package ?? ""), test: String(t.Test ?? ""), outcome: outcome(t.Outcome, `${what} ${i}`) };
  });
}

// isReport tells a report from a registry, so one drop target takes both.
export function isReport(raw: unknown): boolean {
  return isRecord(raw) && Array.isArray(raw.Entries);
}

export function parseReport(raw: unknown): Report {
  if (!isReport(raw)) throw new Error("a report is a JSON object with Entries");
  const results = new Map<string, Result>();
  for (const [i, e] of ((raw as Record<string, unknown>).Entries as unknown[]).entries()) {
    if (!isRecord(e) || typeof e.ID !== "string") throw new Error(`report entry ${i} has no ID`);
    const id = e.ID;
    const scenarios = Array.isArray(e.Scenarios) ? e.Scenarios : [];
    for (const s of scenarios) {
      if (!isRecord(s) || typeof s.Handle !== "string") throw new Error(`${id}: a scenario has no Handle`);
      results.set(`${id}/${s.Handle}`, { outcome: outcome(s.Outcome, `${id}/${s.Handle}`), tests: tests(s.Tests, `${id}/${s.Handle} tests`) });
    }
    const own = tests(e.Tests, `${id} tests`);
    if (scenarios.length === 0) results.set(id, { outcome: outcome(e.Outcome, id), tests: own });
    else if (own.length > 0) results.set(id, { outcome: "none", tests: own });
  }
  return { results };
}
