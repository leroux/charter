// The charter as the viewer reads it: the JSON claimsmap writes, parsed,
// validated and indexed. Pure: no React, no DOM.
//
// Two kinds of spec: needs, in needs.md, and decisions, in the area
// documents; and observations, facts about external systems in Go files. A
// decision serves needs or other decisions and sees related entries; an
// observation supports specs. The earlier form's Realizes counts as Serves,
// and its Concerns and Related as See, as the gate reads them.

export type Level = "spec" | "observation";
export type SpecStatus = "active" | "planned" | "deprecated";
export type ObservationStatus = "observed" | "documented" | "unsure" | "broken";
export type Status = SpecStatus | ObservationStatus;
export const SPEC_STATUSES: SpecStatus[] = ["active", "planned", "deprecated"];
export const OBSERVATION_STATUSES: ObservationStatus[] = ["observed", "documented", "unsure", "broken"];

export type Kind = "invariant" | "bound" | "structure" | "";
export const KINDS: Kind[] = ["invariant", "bound", "structure"];

export interface Scenario {
  handle: string;
  text: string;
  line: number;
}

export interface Change {
  date: string;
  text: string;
  line: number;
}

export interface Anchor {
  kind: "realizes" | "concerns";
  file: string;
  line: number;
  decl: string;
}

export interface Ref {
  file: string;
  func: string;
  line: number;
  // The scenario handles the citation names; none when it cites the spec.
  handles: string[];
}

export interface Waiver {
  pattern: string;
  reason: string;
  file: string;
  line: number;
}

export type AreaKind = "needs" | "decisions" | "observations";

// An area is a document of the charter: needs.md, an area document, or the
// Go file an observation is declared in. Its slug, the file's stem, names it
// in a link.
export interface Area {
  slug: string;
  kind: AreaKind;
  file: string;
  title: string;
  code: string[];
  // The prose that opens the document: what the area holds.
  summary: string;
}

export interface Entry {
  name: string;
  id: string;
  level: Level;
  status: Status;
  need: boolean;
  lens: string;
  kind: Kind;
  title: string;
  file: string;
  pos: string;
  line: number;
  area: string;
  // Ids, each resolved to an entry: what it serves, what it sees.
  serves: string[];
  see: string[];
  accounts: string[];
  text: string;
  why: string;
  guidance: string;
  verification: string;
  scenarios: Scenario[];
  changes: Change[];
  anchors: Anchor[];
  refs: Ref[];
  // Observations only.
  subject: string;
  source: string;
  observed: string;
  supports: string[];
}

export interface Registry {
  charter: string;
  source: string;
  root: string;
  problems: string[];
  areas: Area[];
  areaBySlug: Map<string, Area>;
  waivers: Waiver[];
  entries: Entry[];
  byId: Map<string, Entry>;
  // Entries in each area, in document order.
  inArea: Map<string, Entry[]>;
  servedBy: Map<string, Entry[]>;
  seenBy: Map<string, Entry[]>;
  supportedBy: Map<string, Entry[]>;
}

export function isDecision(e: Entry): boolean {
  return e.level === "spec" && !e.need;
}

function isRecord(v: unknown): v is Record<string, unknown> {
  return typeof v === "object" && v !== null && !Array.isArray(v);
}

function str(v: unknown): string {
  return typeof v === "string" ? v : "";
}

function num(v: unknown): number {
  return typeof v === "number" ? v : 0;
}

function strs(v: unknown, what: string): string[] {
  if (v === undefined || v === null) return [];
  if (!Array.isArray(v) || v.some((x) => typeof x !== "string")) throw new Error(`${what} must be an array of strings`);
  return (v as string[]).slice();
}

function records(v: unknown, what: string): Record<string, unknown>[] {
  if (v === undefined || v === null) return [];
  if (!Array.isArray(v) || v.some((x) => !isRecord(x))) throw new Error(`${what} must be an array of objects`);
  return v as Record<string, unknown>[];
}

// stem is a file's name without its directory and extension.
export function stem(file: string): string {
  const base = file.slice(file.lastIndexOf("/") + 1);
  const dot = base.lastIndexOf(".");
  return dot > 0 ? base.slice(0, dot) : base;
}

function lineOf(pos: string): number {
  const n = Number(pos.slice(pos.lastIndexOf(":") + 1));
  return Number.isFinite(n) ? n : 0;
}

// parseRegistry validates what the viewer relies on and indexes the rest.
// A shape it cannot use throws with the first problem; a relation to an id
// that is no entry is dropped and listed in problems.
export function parseRegistry(raw: unknown): Registry {
  if (!isRecord(raw)) throw new Error("the registry must be a JSON object");
  if (!Array.isArray(raw.entries)) throw new Error("entries must be an array");
  const problems = strs(raw.problems, "problems");

  const areas: Area[] = [];
  const areaBySlug = new Map<string, Area>();
  const areaByFile = new Map<string, Area>();
  const addArea = (file: string, kind: AreaKind, title: string, code: string[], summary = "") => {
    let slug = stem(file) || "unnamed";
    if (areaBySlug.has(slug)) slug = `${slug}-${file.slice(file.lastIndexOf(".") + 1)}`;
    const a: Area = { slug, kind, file, title: title || stem(file), code, summary };
    areas.push(a);
    areaBySlug.set(slug, a);
    areaByFile.set(file, a);
    return a;
  };
  for (const [i, a] of records(raw.areas, "areas").entries()) {
    const file = str(a.file);
    if (!file) throw new Error(`area ${i} has no file`);
    addArea(file, a.needs === true ? "needs" : "decisions", str(a.title), strs(a.code, `area ${file} code`), str(a.summary));
  }

  const entries: Entry[] = [];
  const byId = new Map<string, Entry>();
  for (const [i, e] of raw.entries.entries()) {
    if (!isRecord(e)) throw new Error(`entry ${i} is not an object`);
    const name = str(e.name);
    if (!name) throw new Error(`entry ${i} has no name`);
    const level = str(e.level) as Level;
    if (level !== "spec" && level !== "observation") throw new Error(`${name}: unknown level ${JSON.stringify(e.level)}`);
    const status = str(e.status) as Status;
    const allowed: Status[] = level === "spec" ? SPEC_STATUSES : OBSERVATION_STATUSES;
    if (!allowed.includes(status)) throw new Error(`${name}: unknown status ${JSON.stringify(e.status)} for a ${level}`);
    if (level === "spec" && !str(e.id)) throw new Error(`${name}: a spec needs an id`);
    const id = str(e.id) || name;
    if (byId.has(id)) throw new Error(`${id} is declared twice`);
    const kind = str(e.kind) as Kind;
    const file = str(e.file);
    const pos = str(e.pos);
    const anchors = records(e.anchors, `${id}.anchors`).map((a) => {
      if (a.kind !== "realizes" && a.kind !== "concerns") throw new Error(`${id}: unknown anchor kind ${JSON.stringify(a.kind)}`);
      return { kind: a.kind, file: str(a.file), line: num(a.line), decl: str(a.decl) } as Anchor;
    });
    const entry: Entry = {
      name,
      id,
      level,
      status,
      need: e.need === true,
      lens: str(e.lens),
      kind: KINDS.includes(kind) ? kind : "",
      title: str(e.text),
      file,
      pos,
      line: lineOf(pos),
      area: "",
      serves: [...strs(e.serves, `${id}.serves`), ...strs(e.realizes, `${id}.realizes`)],
      see: [...strs(e.see, `${id}.see`), ...strs(e.concerns, `${id}.concerns`), ...strs(e.related, `${id}.related`)],
      accounts: strs(e.accounts, `${id}.accounts`),
      text: str(e.description),
      why: str(e.rationale),
      guidance: str(e.guidance),
      verification: str(e.verification),
      scenarios: records(e.scenarios, `${id}.scenarios`).map((s) => ({ handle: str(s.handle), text: str(s.text), line: num(s.line) })),
      changes: records(e.changes, `${id}.changes`).map((c) => ({ date: str(c.date), text: str(c.text), line: num(c.line) })),
      anchors,
      refs: records(e.refs, `${id}.refs`).map((r) => ({ file: str(r.file), func: str(r.func), line: num(r.line), handles: strs(r.handles, `${id}.refs.handles`) })),
      subject: str(e.subject),
      source: str(e.source),
      observed: str(e.observed),
      supports: strs(e.supports, `${id}.supports`),
    };
    // An entry whose file names no area document comes from a spec file
    // of the earlier form, gathered in one area, or is an observation,
    // whose Go file is its area.
    let area = areaByFile.get(file);
    if (!area) {
      if (level === "observation") area = addArea(file, "observations", `Observations in ${file.slice(file.lastIndexOf("/") + 1)}`, []);
      else {
        const dir = file.slice(0, file.lastIndexOf("/") + 1) || "specs/";
        area = areaByFile.get(dir) ?? addArea(dir, entry.need ? "needs" : "decisions", "Spec files", []);
        areaByFile.set(dir, area);
      }
    }
    entry.area = area.slug;
    entries.push(entry);
    byId.set(id, entry);
  }

  const known = (e: Entry, rel: string, ids: string[]) =>
    [...new Set(ids)].filter((id) => {
      if (byId.has(id)) return true;
      problems.push(`${e.id} ${rel} ${id}, which is not an entry`);
      return false;
    });
  const servedBy = new Map<string, Entry[]>();
  const seenBy = new Map<string, Entry[]>();
  const supportedBy = new Map<string, Entry[]>();
  const add = (m: Map<string, Entry[]>, k: string, e: Entry) => m.set(k, [...(m.get(k) ?? []), e]);
  for (const e of entries) {
    e.serves = known(e, "serves", e.serves);
    e.see = known(e, "sees", e.see);
    e.supports = known(e, "supports", e.supports);
    for (const id of e.serves) add(servedBy, id, e);
    for (const id of e.see) add(seenBy, id, e);
    for (const id of e.supports) add(supportedBy, id, e);
  }

  const inArea = new Map<string, Entry[]>(areas.map((a) => [a.slug, []]));
  for (const e of entries) inArea.get(e.area)!.push(e);
  for (const list of inArea.values()) list.sort((a, b) => (a.file === b.file ? a.line - b.line : a.file < b.file ? -1 : 1));

  const order: Record<AreaKind, number> = { needs: 0, decisions: 1, observations: 2 };
  // Stable: the decision areas keep the order the export gives them, the
  // charter's reading order.
  areas.sort((a, b) => order[a.kind] - order[b.kind]);

  return {
    charter: str(raw.charter),
    source: str(raw.source),
    root: str(raw.root),
    problems,
    areas,
    areaBySlug,
    waivers: records(raw.waivers, "waivers").map((w) => ({ pattern: str(w.pattern), reason: str(w.reason), file: str(w.file), line: num(w.line) })),
    entries,
    byId,
    inArea,
    servedBy,
    seenBy,
    supportedBy,
  };
}

// needsOf are the needs an entry serves, directly or through the decisions
// it serves, in id order.
export function needsOf(reg: Registry, e: Entry): Entry[] {
  const seen = new Set<string>();
  const needs: Entry[] = [];
  const walk = (x: Entry) => {
    for (const id of x.serves) {
      if (seen.has(id)) continue;
      seen.add(id);
      const up = reg.byId.get(id)!;
      if (up.need) needs.push(up);
      else walk(up);
    }
  };
  walk(e);
  return needs.sort((a, b) => (a.id < b.id ? -1 : 1));
}

export interface TreeNode {
  entry: Entry;
  children: TreeNode[];
  // Set when the entry was already shown higher in the tree.
  repeat: boolean;
}

// servingTree is what serves an entry, transitively: each decision under
// what it serves, by area in the order given, then as its document has
// them. An entry reached twice is shown again without its children, so a
// cycle or a diamond stays finite.
export function servingTree(
  reg: Registry,
  e: Entry,
  keep: (x: Entry) => boolean = () => true,
  areaOrder: (a: string, b: string) => number = (a, b) => (a < b ? -1 : a > b ? 1 : 0),
): TreeNode[] {
  const shown = new Set<string>([e.id]);
  const build = (x: Entry): TreeNode[] =>
    (reg.servedBy.get(x.id) ?? [])
      .filter(keep)
      .sort((a, b) => areaOrder(a.area, b.area) || a.line - b.line)
      .map((c) => {
        if (shown.has(c.id)) return { entry: c, children: [], repeat: true };
        shown.add(c.id);
        return { entry: c, children: build(c), repeat: false };
      });
  return build(e);
}

export type Tone = "good" | "warning" | "critical" | "muted";

export function statusTone(s: Status): Tone {
  switch (s) {
    case "active":
    case "observed":
      return "good";
    case "unsure":
      return "warning";
    case "broken":
      return "critical";
    default:
      return "muted";
  }
}

export function statusGlyph(s: Status): string {
  switch (s) {
    case "active":
    case "observed":
      return "●";
    case "planned":
      return "◌";
    case "deprecated":
      return "⊘";
    case "documented":
      return "≡";
    case "unsure":
      return "?";
    case "broken":
      return "✕";
  }
}
