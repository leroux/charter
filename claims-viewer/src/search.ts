// Search over every part of an entry a person would look for it by: id,
// title, text, scenarios, surfaces, the code that anchors it and the tests
// that cite it. Every word of the query must match somewhere. Pure.

import type { Entry, Registry } from "./model";

export interface Hit {
  entry: Entry;
  // Where the best match is: "id", "title", or the part of the entry.
  where: string;
  snippet: string;
  score: number;
}

function fields(e: Entry): [string, string][] {
  return [
    ["id", e.id],
    ["title", e.title],
    ["text", e.text],
    ["why", e.why],
    ["guidance", e.guidance],
    ...e.scenarios.map((s): [string, string] => [`scenario ${s.handle}`, `${s.handle} ${s.text}`]),
    ...e.changes.map((c): [string, string] => [`change ${c.date}`, c.text]),
    ["accounts", e.accounts.join(" ")],
    ...e.anchors.map((a): [string, string] => ["code", `${a.file} ${a.decl}`]),
    ...e.refs.map((r): [string, string] => ["test", `${r.file} ${r.func}`]),
    ["source", `${e.subject} ${e.source}`],
  ];
}

const WEIGHT: Record<string, number> = { id: 100, title: 20 };

export function search(reg: Registry, query: string, limit = 200): Hit[] {
  const words = query.toLowerCase().split(/\s+/).filter(Boolean);
  if (words.length === 0) return [];
  const hits: Hit[] = [];
  for (const e of reg.entries) {
    const fs = fields(e).map(([k, v]) => [k, v, v.toLowerCase()] as const);
    if (!words.every((w) => fs.some(([, , low]) => low.includes(w)))) continue;
    let score = 0;
    let best: (typeof fs)[number] | null = null;
    let bestWeight = -1;
    for (const f of fs) {
      const n = words.filter((w) => f[2].includes(w)).length;
      if (n === 0) continue;
      const weight = (WEIGHT[f[0].split(" ")[0]] ?? 1) * n;
      score += weight;
      if (weight > bestWeight) {
        best = f;
        bestWeight = weight;
      }
    }
    if (e.id.toLowerCase() === words.join(" ")) score += 1000;
    if (e.status === "deprecated") score /= 4;
    hits.push({ entry: e, where: best![0], snippet: snippet(best![1], best![2], words[0]), score });
  }
  return hits.sort((a, b) => b.score - a.score || (a.entry.id < b.entry.id ? -1 : 1)).slice(0, limit);
}

// snippet is about a hundred characters of text around the first match.
function snippet(text: string, low: string, word: string): string {
  const flat = text.replace(/\s+/g, " ");
  if (flat.length <= 120) return flat;
  const at = Math.max(0, low.replace(/\s+/g, " ").indexOf(word) - 40);
  return (at > 0 ? "…" : "") + flat.slice(at, at + 120) + (at + 120 < flat.length ? "…" : "");
}
