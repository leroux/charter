// The walkthrough: the area documents in the charter's own reading order,
// which the export gives them in (reading-order.txt, checked by the gate),
// each with the needs it does most to deliver. Nothing here names a
// project's areas; a charter without a reading order reads by file name.
// Pure.

import { isDecision, type Area, type Entry, type Registry } from "./model";

export interface Tour {
  // The decision areas in reading order, which the numbers, the sidebar
  // and prev and next follow.
  path: Area[];
  // Each area's place in path, from 0.
  index: Map<string, number>;
  // Each need's home: the area whose active decisions most often reach
  // it, the earlier in the reading order on a tie, as the digest places
  // it; "" when none does.
  homeOf: Map<string, string>;
  // The needs at home in each area, active before planned, then as
  // needs.md has them.
  needsAt: Map<string, Entry[]>;
  // Needs no active decision reaches.
  unserved: Entry[];
}

export function arrange(reg: Registry): Tour {
  const path = reg.areas.filter((a) => a.kind === "decisions");
  const index = new Map(path.map((a, i) => [a.slug, i]));
  const homeOf = new Map<string, string>();
  const needsAt = new Map<string, Entry[]>(path.map((a) => [a.slug, []]));
  const unserved: Entry[] = [];
  for (const need of reg.entries.filter((e) => e.need)) {
    const counts = new Map<string, number>();
    for (const d of servers(reg, need)) if (d.status === "active") counts.set(d.area, (counts.get(d.area) ?? 0) + 1);
    let home = "";
    for (const [slug, n] of counts) {
      if (!index.has(slug)) continue;
      const best = counts.get(home) ?? 0;
      if (n > best || (n === best && index.get(slug)! < (index.get(home) ?? Infinity))) home = slug;
    }
    homeOf.set(need.id, home);
    if (home) needsAt.get(home)!.push(need);
    else unserved.push(need);
  }
  for (const ns of needsAt.values()) ns.sort((a, b) => byStatus(a, b) || a.line - b.line);
  return { path, index, homeOf, needsAt, unserved };
}

// servers are the live decisions that serve an entry, transitively.
export function servers(reg: Registry, e: Entry): Entry[] {
  const seen = new Set<string>([e.id]);
  const out: Entry[] = [];
  const walk = (x: Entry) => {
    for (const d of reg.servedBy.get(x.id) ?? []) {
      if (seen.has(d.id) || d.status === "deprecated" || !isDecision(d)) continue;
      seen.add(d.id);
      out.push(d);
      walk(d);
    }
  };
  walk(e);
  return out;
}

// needsServedBy are the needs an area's live decisions reach, each with how
// many of the area's decisions reach it, most first.
export function needsServedBy(reg: Registry, slug: string): { need: Entry; count: number }[] {
  const counts = new Map<string, number>();
  for (const d of reg.inArea.get(slug) ?? []) {
    if (!isDecision(d) || d.status === "deprecated") continue;
    const seen = new Set<string>();
    const up = (x: Entry) => {
      for (const id of x.serves) {
        if (seen.has(id)) continue;
        seen.add(id);
        const p = reg.byId.get(id)!;
        if (p.need) counts.set(id, (counts.get(id) ?? 0) + 1);
        else up(p);
      }
    };
    up(d);
  }
  return [...counts]
    .map(([id, count]) => ({ need: reg.byId.get(id)!, count }))
    .filter((x) => x.need.status !== "deprecated")
    .sort((a, b) => b.count - a.count || a.need.line - b.need.line);
}

const STATUS_ORDER: Record<string, number> = { active: 0, planned: 1, deprecated: 2 };

function byStatus(a: Entry, b: Entry): number {
  return (STATUS_ORDER[a.status] ?? 0) - (STATUS_ORDER[b.status] ?? 0);
}

// areaNumber is an area's place in the walkthrough, from 1, or "".
export function areaNumber(tour: Tour, slug: string): string {
  const i = tour.index.get(slug);
  return i === undefined ? "" : String(i + 1);
}

// areaOrder compares two areas by the walkthrough: needs first, then the
// path, then the rest by name.
export function areaOrder(tour: Tour): (a: string, b: string) => number {
  const rank = (s: string) => tour.index.get(s) ?? Infinity;
  return (a, b) => rank(a) - rank(b) || (a < b ? -1 : a > b ? 1 : 0);
}

// neighbours are the areas before and after one on the path.
export function neighbours(tour: Tour, slug: string): { prev: Area | null; next: Area | null } {
  const i = tour.index.get(slug);
  if (i === undefined) return { prev: null, next: tour.path[0] ?? null };
  return { prev: tour.path[i - 1] ?? null, next: tour.path[i + 1] ?? null };
}

// firstSentence is the opening of an area's summary, for a line under its
// title.
export function firstSentence(text: string): string {
  const para = (text.split(/\n\s*\n/)[0] ?? "").replace(/\s+/g, " ").trim();
  const m = para.match(/^.+?[.!?](?=\s|$)/);
  return m ? m[0] : para;
}
