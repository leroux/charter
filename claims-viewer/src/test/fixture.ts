// A small charter in the shape claimsmap writes: two needs, an area of
// decisions in every state coverage distinguishes, a waiver, an
// observation, and a decision that reaches its need through another.

type Raw = Record<string, unknown>;

export function spec(id: string, file: string, line: number, fields: Raw = {}): Raw {
  return {
    name: id.replace(/-/g, "_"),
    id,
    level: "spec",
    status: "active",
    need: file.endsWith("needs.md"),
    lens: id.slice(0, id.indexOf("-")),
    kind: "",
    text: `${id} title`,
    file,
    pos: `${file}:${line}`,
    serves: [],
    see: [],
    realizes: [],
    concerns: [],
    related: [],
    accounts: [],
    description: "",
    rationale: "",
    guidance: "",
    verification: "",
    scenarios: [],
    changes: [],
    anchors: [],
    refs: [],
    supports: [],
    subject: "",
    source: "",
    observed: "",
    form: "area",
    ...fields,
  };
}

const realized = [{ kind: "realizes", file: "go/a/a.go", line: 3, decl: "func A()" }];
const scenario = (handle: string, line: number) => ({ handle, text: `when ${handle}`, line });
const test = (func: string, ...handles: string[]) => ({ file: "go/a/a_test.go", func, line: 10, handles });

export function fixture(): Raw {
  return {
    charter: "product",
    source: "claims @ main abc1234",
    root: "/src/product",
    problems: [],
    areas: [
      { file: "claims/needs.md", title: "Needs", code: [], needs: true },
      { file: "claims/refill.md", title: "Bucket refill", code: ["pkg/limit"], needs: false },
    ],
    waivers: [{ pattern: "SW-003/one", reason: "verified by review", file: "claims/waivers.txt", line: 3 }],
    entries: [
      spec("STK-001", "claims/needs.md", 5, {
        scenarios: [scenario("edit-takes-effect", 12)],
        refs: [{ file: "web/refill.spec.ts", func: "edit takes effect", line: 4, handles: ["edit-takes-effect"] }],
      }),
      spec("STK-002", "claims/needs.md", 20, { status: "deprecated" }),
      // Covered: realized, both scenarios tested.
      spec("SW-001", "claims/refill.md", 30, {
        kind: "invariant",
        serves: ["STK-001"],
        see: ["SW-002"],
        description: "The bucket refills at the rate.\nSee `SW-002`, and SW-002.",
        scenarios: [scenario("inserted-runs-next", 40), scenario("burst-caps-tokens", 42)],
        anchors: realized,
        refs: [test("TestInserted", "inserted-runs-next"), test("TestBranch", "burst-caps-tokens")],
      }),
      // A gap: no code, one scenario untested; serves through SW-001.
      spec("SW-002", "claims/refill.md", 10, {
        serves: ["SW-001"],
        scenarios: [scenario("tested", 14), scenario("untested", 15)],
        refs: [test("TestTested", "tested")],
      }),
      // Waived: realized, its scenario accepted by a waiver; serves nothing.
      spec("SW-003", "claims/refill.md", 50, { scenarios: [scenario("one", 55)], anchors: realized, see: ["STK-002"] }),
      // Planned and cited already; deprecated and still anchored.
      spec("SW-004", "claims/refill.md", 60, { status: "planned", serves: ["STK-001"], refs: [test("TestEarly")] }),
      spec("SW-005", "claims/refill.md", 70, { status: "deprecated", anchors: realized }),
      {
        name: "LockIsAdvisory",
        id: "LockIsAdvisory",
        level: "observation",
        status: "unsure",
        text: "flock is advisory",
        file: "claims/linux.go",
        pos: "claims/linux.go:9",
        supports: ["SW-001"],
        subject: "linux",
        source: "https://man7.org",
        refs: [],
      },
    ],
  };
}

// report is what claimsgen report writes over the fixture: SW-001's
// branch test failed.
export function report(): Raw {
  return {
    Entries: [
      {
        ID: "SW-001",
        Title: "SW-001 title",
        Status: "active",
        Need: false,
        Scenarios: [
          { Handle: "inserted-runs-next", Tests: [{ Package: "a", Test: "TestInserted", Outcome: "pass" }], Outcome: "pass" },
          { Handle: "burst-caps-tokens", Tests: [{ Package: "a", Test: "TestBranch", Outcome: "fail" }], Outcome: "fail" },
        ],
        Tests: null,
        Outcome: "fail",
      },
    ],
    Unverified: ["SW-001/burst-caps-tokens"],
  };
}
