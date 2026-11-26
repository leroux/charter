import { describe, expect, it } from "vitest";
import { needsOf, parseRegistry, servingTree } from "./model";
import { fixture, spec } from "./test/fixture";

describe("parseRegistry", () => {
  it("reads the areas, needs first, and an observation's file as its own area", () => {
    const reg = parseRegistry(fixture());
    expect(reg.areas.map((a) => [a.slug, a.kind, a.title])).toEqual([
      ["needs", "needs", "Needs"],
      ["refill", "decisions", "Bucket refill"],
      ["linux", "observations", "Observations in linux.go"],
    ]);
    expect(reg.byId.get("LockIsAdvisory")!.area).toBe("linux");
  });

  it("lists an area's entries in document order", () => {
    const reg = parseRegistry(fixture());
    expect(reg.inArea.get("refill")!.map((e) => e.id)).toEqual(["SW-002", "SW-001", "SW-003", "SW-004", "SW-005"]);
  });

  it("indexes what serves, sees and supports an entry", () => {
    const reg = parseRegistry(fixture());
    expect(reg.servedBy.get("STK-001")!.map((e) => e.id)).toEqual(["SW-001", "SW-004"]);
    expect(reg.seenBy.get("SW-002")!.map((e) => e.id)).toEqual(["SW-001"]);
    expect(reg.supportedBy.get("SW-001")!.map((e) => e.id)).toEqual(["LockIsAdvisory"]);
  });

  it("reads a spec file's Realizes as Serves, and its Concerns and Related as See", () => {
    const raw = fixture();
    (raw.entries as unknown[]).push(
      spec("SW-009", "claims/specs/SW-009-x.md", 1, { form: "file", realizes: ["STK-001"], concerns: ["SW-001"], related: ["SW-002"] }),
    );
    const reg = parseRegistry(raw);
    const e = reg.byId.get("SW-009")!;
    expect(e.serves).toEqual(["STK-001"]);
    expect(e.see).toEqual(["SW-001", "SW-002"]);
    expect(reg.areaBySlug.get(e.area)!.title).toBe("Spec files");
  });

  it("drops a relation to an id that is no entry and lists it", () => {
    const raw = fixture();
    (raw.entries as Record<string, unknown>[])[2].serves = ["STK-001", "STK-404"];
    const reg = parseRegistry(raw);
    expect(reg.byId.get("SW-001")!.serves).toEqual(["STK-001"]);
    expect(reg.problems).toContain("SW-001 serves STK-404, which is not an entry");
  });

  it("refuses a shape it cannot use", () => {
    expect(() => parseRegistry([])).toThrow("JSON object");
    expect(() => parseRegistry({ entries: [spec("SW-1", "claims/a.md", 1, { status: "done" })] })).toThrow('unknown status "done"');
    expect(() => parseRegistry({ entries: [spec("SW-1", "claims/a.md", 1), spec("SW-1", "claims/a.md", 2)] })).toThrow("declared twice");
  });
});

describe("relations", () => {
  it("finds the needs a decision reaches through another decision", () => {
    const reg = parseRegistry(fixture());
    expect(needsOf(reg, reg.byId.get("SW-002")!).map((e) => e.id)).toEqual(["STK-001"]);
    expect(needsOf(reg, reg.byId.get("SW-003")!)).toEqual([]);
  });

  it("builds the tree of what serves a need, each decision under what it serves", () => {
    const reg = parseRegistry(fixture());
    const tree = servingTree(reg, reg.byId.get("STK-001")!);
    expect(tree.map((n) => [n.entry.id, n.children.map((c) => c.entry.id)])).toEqual([
      ["SW-001", ["SW-002"]],
      ["SW-004", []],
    ]);
  });

  it("shows an entry reached twice once in full, so a cycle stays finite", () => {
    const raw = fixture();
    (raw.entries as Record<string, unknown>[])[2].serves = ["STK-001", "SW-002"];
    const reg = parseRegistry(raw);
    const tree = servingTree(reg, reg.byId.get("STK-001")!);
    const sw001 = tree.find((n) => n.entry.id === "SW-001")!;
    expect(sw001.children[0].entry.id).toBe("SW-002");
    expect(sw001.children[0].children[0]).toMatchObject({ repeat: true });
  });
});
