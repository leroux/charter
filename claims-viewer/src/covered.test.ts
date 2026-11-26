import { describe, expect, it } from "vitest";
import { coverage, tally } from "./covered";
import { gateCount, gaps } from "./gaps";
import { parseRegistry } from "./model";
import { parseReport } from "./report";
import { fixture, report } from "./test/fixture";

const reg = () => parseRegistry(fixture());

describe("coverage", () => {
  it("calls a realized decision with every scenario tested covered", () => {
    const r = reg();
    const c = coverage(r, r.byId.get("SW-001")!, null);
    expect(c.realized).toBe(true);
    expect(c.checks.map((x) => [x.handle, x.mark])).toEqual([
      ["inserted-runs-next", "cited"],
      ["burst-caps-tokens", "cited"],
    ]);
    expect(c.verdict).toBe("ok");
  });

  it("calls an unrealized decision with an untested scenario a gap", () => {
    const r = reg();
    const c = coverage(r, r.byId.get("SW-002")!, null);
    expect(c.realized).toBe(false);
    expect(c.checks.map((x) => x.mark)).toEqual(["cited", "missing"]);
    expect(c.verdict).toBe("gap");
  });

  it("marks a scenario a waiver names as waived, with the reason", () => {
    const r = reg();
    const c = coverage(r, r.byId.get("SW-003")!, null);
    expect(c.checks[0]).toMatchObject({ mark: "waived", waiver: { reason: "verified by review" } });
    expect(c.verdict).toBe("waived");
  });

  it("asks a need for tests only, and nothing of a planned or deprecated spec", () => {
    const r = reg();
    expect(coverage(r, r.byId.get("STK-001")!, null)).toMatchObject({ realized: null, verdict: "ok" });
    expect(coverage(r, r.byId.get("SW-004")!, null)).toMatchObject({ verdict: "none", checks: [{ mark: "cited" }] });
    expect(coverage(r, r.byId.get("SW-005")!, null)).toMatchObject({ verdict: "none", checks: [{ mark: "open" }] });
  });

  it("takes pass and fail from a report", () => {
    const r = reg();
    const c = coverage(r, r.byId.get("SW-001")!, parseReport(report()));
    expect(c.checks.map((x) => x.mark)).toEqual(["pass", "fail"]);
    expect(c.checks[1].tests).toEqual([{ pkg: "a", test: "TestBranch", outcome: "fail" }]);
    expect(c.verdict).toBe("fail");
  });

  it("tallies an area's active specs", () => {
    const r = reg();
    const t = tally(r, r.inArea.get("refill")!, null);
    expect([t.active, t.planned, t.deprecated, t.decisions, t.realized]).toEqual([3, 1, 1, 3, 2]);
    expect(t.marks).toMatchObject({ cited: 3, waived: 1, missing: 1 });
  });
});

describe("gaps", () => {
  it("lists what the gate fails on, and what it lets through", () => {
    const sections = gaps(reg(), null);
    const ids = (key: string) => sections.find((s) => s.key === key)!.items.map((i) => (i.handle ? `${i.entry.id}/${i.handle}` : i.entry.id));
    expect(ids("unrealized")).toEqual(["SW-002"]);
    expect(ids("untested")).toEqual(["SW-002/untested"]);
    expect(ids("waived")).toEqual(["SW-003/one"]);
    expect(ids("unserved")).toEqual(["SW-003"]);
    expect(ids("stale")).toEqual(["SW-003"]);
    expect(ids("deprecated-cited")).toEqual(["SW-005"]);
    expect(ids("planned")).toEqual(["SW-004"]);
    expect(ids("observations")).toEqual(["LockIsAdvisory"]);
    expect(gateCount(sections)).toBe(2);
  });

  it("lists a failing test from a report", () => {
    const sections = gaps(reg(), parseReport(report()));
    expect(sections.find((s) => s.key === "failing")!.items.map((i) => i.handle)).toEqual(["burst-caps-tokens"]);
  });
});
