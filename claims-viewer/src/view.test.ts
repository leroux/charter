import { describe, expect, it } from "vitest";
import { dataURL, loadFromText, loadFromURL, loadIndex, loadReport, parseIndex, reportURL, sameData } from "./load";
import { parseRegistry } from "./model";
import { blocks, inline } from "./prose";
import { href, parseRoute } from "./route";
import { search } from "./search";
import { fixture, report } from "./test/fixture";

describe("search", () => {
  const reg = parseRegistry(fixture());
  it("puts the entry an id names first", () => {
    expect(search(reg, "sw-002")[0].entry.id).toBe("SW-002");
  });
  it("finds an entry by a scenario handle, a test and a code site, and says where", () => {
    expect(search(reg, "untested").map((h) => [h.entry.id, h.where])).toContainEqual(["SW-002", "scenario untested"]);
    expect(search(reg, "TestBranch").map((h) => h.entry.id)).toEqual(["SW-001"]);
    expect(search(reg, "func A").map((h) => h.entry.id)).toContain("SW-003");
  });
  it("needs every word to match", () => {
    expect(search(reg, "bucket nothing")).toEqual([]);
  });
});

describe("routes", () => {
  it("round-trips every page", () => {
    for (const r of [
      { page: "overview" },
      { page: "area", slug: "bucket-log" },
      { page: "entry", id: "SW-261001-K7Q2" },
      { page: "gaps" },
      { page: "search", q: "a / b" },
    ] as const) {
      expect(parseRoute(href(r))).toEqual(r);
    }
    expect(parseRoute("")).toEqual({ page: "overview" });
    expect(parseRoute("#/nope")).toEqual({ page: "overview" });
  });
});

describe("prose", () => {
  const isId = (id: string) => id === "SW-002" || id === "SW-261001-K7Q2";
  it("joins a hard-wrapped paragraph and splits on blank lines and lists", () => {
    const bs = blocks("one\ntwo\n\n- a\n  more\n- b", isId);
    expect(bs).toEqual([
      { t: "p", v: [{ t: "text", v: "one two" }] },
      { t: "list", v: [[{ t: "text", v: "a more" }], [{ t: "text", v: "b" }]] },
    ]);
  });
  it("links known ids, numbered or dated, and leaves code alone", () => {
    expect(inline("see SW-002 and `SW-002`, SW-999, SW-261001-K7Q2", isId)).toEqual([
      { t: "text", v: "see " },
      { t: "id", v: "SW-002" },
      { t: "text", v: " and " },
      { t: "code", v: "SW-002" },
      { t: "text", v: ", SW-999, " },
      { t: "id", v: "SW-261001-K7Q2" },
    ]);
  });
});

describe("load", () => {
  it("picks the charter's file from the query, and its report beside it", () => {
    expect(dataURL("")).toBe("./claims.json");
    expect(dataURL("?data=other.json")).toBe("other.json");
    expect(reportURL("?data=other.json")).toBe("other-report.json");
    expect(reportURL("?data=a.json&report=r.json")).toBe("r.json");
    expect(sameData("./claims.json", "claims.json")).toBe(true);
  });

  it("reads the charters to switch between, naming one by its file when it has no label", async () => {
    expect(parseIndex([{ data: "claims.json", label: "product" }, { data: "tooling.json" }])).toEqual([
      { data: "claims.json", label: "product" },
      { data: "tooling.json", label: "tooling" },
    ]);
    expect(() => parseIndex({})).toThrow("must be an array");
    const missing = (async () => new Response("", { status: 404 })) as unknown as typeof fetch;
    expect(await loadIndex("./charters.json", missing)).toEqual([]);
  });

  it("tells a registry from a report", () => {
    expect(loadFromText(JSON.stringify(fixture()), "a").kind).toBe("registry");
    expect(loadFromText(JSON.stringify(report()), "b").kind).toBe("report");
    expect(() => loadFromText("{", "c")).toThrow("c: not JSON");
  });

  it("loads from a URL and says which URL failed", async () => {
    const ok = (async () => new Response(JSON.stringify(fixture()))) as unknown as typeof fetch;
    expect((await loadFromURL("./claims.json", ok)).entries.length).toBe(8);
    const missing = (async () => new Response("", { status: 404, statusText: "Not Found" })) as unknown as typeof fetch;
    await expect(loadFromURL("./claims.json", missing)).rejects.toThrow("./claims.json: 404 Not Found");
  });

  it("takes a missing report, or a dev server's page in its place, as none", async () => {
    const missing = (async () => new Response("", { status: 404 })) as unknown as typeof fetch;
    expect(await loadReport("./claims-report.json", missing)).toBeNull();
    const page = (async () => new Response("<!doctype html>")) as unknown as typeof fetch;
    expect(await loadReport("./claims-report.json", page)).toBeNull();
    const ok = (async () => new Response(JSON.stringify(report()))) as unknown as typeof fetch;
    expect((await loadReport("./claims-report.json", ok))!.results.get("SW-001/burst-caps-tokens")!.outcome).toBe("fail");
  });
});
