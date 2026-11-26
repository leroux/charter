import { describe, expect, it } from "vitest";
import { parseRegistry, servingTree } from "./model";
import { fixture, spec } from "./test/fixture";
import { areaNumber, areaOrder, arrange, firstSentence, needsServedBy, neighbours } from "./tour";

// Two more areas beside refill, in the order the export gives them, which
// is the charter's reading order: waits, then refill, then extra.
function charter() {
  const raw = fixture();
  const areas = raw.areas as Record<string, unknown>[];
  raw.areas = [
    areas[0],
    { file: "claims/waits.md", title: "Spawn, queue and wait", code: [], needs: false, summary: "How a caller waits.\n\nMore." },
    areas[1],
    { file: "claims/extra.md", title: "Extra", code: [], needs: false },
  ];
  (raw.entries as unknown[]).push(
    spec("SW-010", "claims/waits.md", 5, { serves: ["STK-001"] }),
    spec("SW-011", "claims/waits.md", 9, { serves: ["STK-001"] }),
    spec("SW-012", "claims/extra.md", 5, { serves: ["SW-010"] }),
  );
  return parseRegistry(raw);
}

describe("arrange", () => {
  it("walks the areas in the order the charter gives them, needs.md apart", () => {
    const tour = arrange(charter());
    expect(tour.path.map((a) => a.slug)).toEqual(["waits", "refill", "extra"]);
    expect(areaNumber(tour, "refill")).toBe("2");
    expect(areaNumber(tour, "needs")).toBe("");
  });

  it("homes a need in the area whose active decisions most often reach it", () => {
    const reg = charter();
    reg.byId.get("SW-011")!.status = "planned";
    // STK-001 is reached by two active decisions in refill (SW-001, SW-002
    // through it), one in waits, the planned SW-011 not counting, and one
    // in extra.
    const tour = arrange(reg);
    expect(tour.homeOf.get("STK-001")).toBe("refill");
    expect(tour.needsAt.get("refill")!.map((n) => n.id)).toEqual(["STK-001"]);
    expect(tour.unserved.map((n) => n.id)).toEqual(["STK-002"]);
  });

  it("breaks a tie for the area earlier in the reading order", () => {
    // Two active decisions each in waits and refill: waits is read first.
    expect(arrange(charter()).homeOf.get("STK-001")).toBe("waits");
  });

  it("gives the areas before and after one", () => {
    const tour = arrange(charter());
    expect(neighbours(tour, "waits")).toMatchObject({ prev: null, next: { slug: "refill" } });
    expect(neighbours(tour, "extra")).toMatchObject({ prev: { slug: "refill" }, next: null });
  });

  it("orders the serving tree by the reading order, then as the document has them", () => {
    const reg = charter();
    const tree = servingTree(reg, reg.byId.get("STK-001")!, () => true, areaOrder(arrange(reg)));
    expect(tree.map((n) => n.entry.id)).toEqual(["SW-010", "SW-011", "SW-001", "SW-004"]);
  });
});

describe("needsServedBy", () => {
  it("counts the area's live decisions that reach each need, through others too", () => {
    const reg = charter();
    expect(needsServedBy(reg, "refill").map((x) => [x.need.id, x.count])).toEqual([["STK-001", 3]]);
    expect(needsServedBy(reg, "extra").map((x) => [x.need.id, x.count])).toEqual([["STK-001", 1]]);
  });
});

describe("firstSentence", () => {
  it("is the summary's opening sentence", () => {
    expect(firstSentence("How a caller waits. And more.\n\nMore.")).toBe("How a caller waits.");
    expect(firstSentence("Who decides: the gate, the policy. Then more.")).toBe("Who decides: the gate, the policy.");
    expect(firstSentence("No full stop")).toBe("No full stop");
  });
});
