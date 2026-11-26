// Gaps: the work list, what a person improving the charter's coverage or
// reviewing it would look at, each item an entry and why. The gate fails
// the first sections, an unrealized decision and an untested scenario; the
// rest are what it does not fail on. Pure.

import { coverage, type Check } from "./covered";
import { isDecision, needsOf, type Entry, type Registry } from "./model";
import type { Report } from "./report";

export interface GapItem {
  entry: Entry;
  // The scenario's handle when the item is one scenario.
  handle: string;
  note: string;
}

export interface GapSection {
  key: string;
  title: string;
  hint: string;
  // Whether the gate fails the charter for an item of this section.
  gate: boolean;
  items: GapItem[];
}

export function gaps(reg: Registry, report: Report | null): GapSection[] {
  const sections: Record<string, GapSection> = {};
  const section = (key: string, title: string, hint: string, gate = false): GapSection =>
    (sections[key] = { key, title, hint, gate, items: [] as GapItem[] });
  const failing = section("failing", "Failing tests", "a citing Go test failed in the loaded report", true);
  const unrealized = section("unrealized", "Active decisions no code realizes", "add //charter:realizes ID to the declaration that carries it", true);
  const untested = section("untested", "Active scenarios no test verifies", "cite the scenario from a test, or waive it with a reason", true);
  const waived = section("waived", "Waived", "accepted without a test; each is debt with a reason");
  const unserved = section("unserved", "Active decisions that serve no need", "a decision reaches a need through Serves, directly or through another decision");
  const stale = section("stale", "Active entries that point at deprecated ones", "Serves or See names an entry that has been retired");
  const deprecatedCited = section("deprecated-cited", "Deprecated entries still cited or anchored", "the tests or code may verify a decision that no longer holds");
  const planned = section("planned", "Planned", "ruled, not built; a test may cite one already");
  const observations = section("observations", "Observations in doubt", "a fact the specs rest on that is unsure or broken");

  const scenarioNote = (c: Check) => (c.waiver ? c.waiver.reason : "");
  for (const e of reg.entries) {
    if (e.level === "observation") {
      if (e.status === "unsure" || e.status === "broken") {
        observations.items.push({ entry: e, handle: "", note: `${e.status}; supports ${e.supports.join(", ") || "nothing"}` });
      }
      continue;
    }
    const cov = coverage(reg, e, report);
    if (e.status === "deprecated") {
      const n = e.refs.length + e.anchors.length;
      if (n > 0) deprecatedCited.items.push({ entry: e, handle: "", note: `${e.refs.length} citing tests, ${e.anchors.length} anchors` });
      continue;
    }
    if (e.status === "planned") {
      const cited = cov.checks.filter((c) => c.refs.length > 0).length;
      planned.items.push({ entry: e, handle: "", note: cited > 0 ? `${cited} of ${cov.checks.length} already tested` : "" });
      continue;
    }
    if (cov.realized === false) unrealized.items.push({ entry: e, handle: "", note: "" });
    for (const c of cov.checks) {
      const item = { entry: e, handle: c.handle, note: scenarioNote(c) };
      if (c.mark === "fail") failing.items.push(item);
      else if (c.mark === "missing") untested.items.push(item);
      else if (c.mark === "waived") waived.items.push(item);
    }
    if (isDecision(e) && needsOf(reg, e).length === 0) {
      unserved.items.push({ entry: e, handle: "", note: e.serves.length ? `serves ${e.serves.join(", ")}` : "serves nothing" });
    }
    const retired = [...e.serves, ...e.see].filter((id) => reg.byId.get(id)?.status === "deprecated");
    if (retired.length > 0) stale.items.push({ entry: e, handle: "", note: retired.join(", ") });
  }
  return Object.values(sections);
}

// gateCount is how many items the gate would fail on.
export function gateCount(sections: GapSection[]): number {
  return sections.filter((s) => s.gate).reduce((n, s) => n + s.items.length, 0);
}
