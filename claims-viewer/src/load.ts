// Where the charter comes from: the file ?data= names, claims.json beside
// the page by default, or a file the person drops; the report beside it,
// if one was written; and charters.json, which lists the charters the page
// can switch between. The DOM parts stay in App.tsx.

import { parseRegistry, type Registry } from "./model";
import { isReport, parseReport, type Report } from "./report";

export function dataURL(search: string): string {
  return new URLSearchParams(search).get("data") || "./claims.json";
}

// reportURL is the report beside a registry: claims.json's is
// claims-report.json.
export function reportURL(search: string): string {
  const report = new URLSearchParams(search).get("report");
  if (report) return report;
  return dataURL(search).replace(/(\.json)?$/, "-report.json");
}

// A charter the page can show: its export, and how the switch names it.
export interface CharterLink {
  data: string;
  label: string;
}

export const INDEX_URL = "./charters.json";

// parseIndex reads charters.json: [{"data": "claims.json", "label": "…"}].
export function parseIndex(raw: unknown): CharterLink[] {
  if (!Array.isArray(raw)) throw new Error("charters.json must be an array");
  return raw.map((c, i) => {
    if (typeof c !== "object" || c === null || typeof (c as CharterLink).data !== "string") {
      throw new Error(`charters.json entry ${i} has no data`);
    }
    const { data, label } = c as CharterLink;
    return { data, label: typeof label === "string" && label ? label : data.replace(/\.json$/, "") };
  });
}

// sameData says whether two data URLs name one file, ./claims.json being
// claims.json.
export function sameData(a: string, b: string): boolean {
  return a.replace(/^\.\//, "") === b.replace(/^\.\//, "");
}

function parseJSON(text: string, name: string): unknown {
  try {
    return JSON.parse(text);
  } catch (err) {
    throw new Error(`${name}: not JSON: ${(err as Error).message}`);
  }
}

export type Loaded = { kind: "registry"; reg: Registry } | { kind: "report"; report: Report };

// loadFromText reads a registry or a report, whichever the text holds.
export function loadFromText(text: string, name: string): Loaded {
  const raw = parseJSON(text, name);
  try {
    return isReport(raw) ? { kind: "report", report: parseReport(raw) } : { kind: "registry", reg: parseRegistry(raw) };
  } catch (err) {
    throw new Error(`${name}: ${(err as Error).message}`);
  }
}

export async function loadFromURL(url: string, fetchFn: typeof fetch = fetch): Promise<Registry> {
  const res = await fetchFn(url);
  if (!res.ok) throw new Error(`${url}: ${res.status} ${res.statusText}`.trim());
  const got = loadFromText(await res.text(), url);
  if (got.kind !== "registry") throw new Error(`${url}: a report, not a registry`);
  return got.reg;
}

// optionalJSON is the JSON at url, or null when there is none. A dev server
// answers a missing file with the page, which is none too.
async function optionalJSON(url: string, fetchFn: typeof fetch): Promise<unknown> {
  const res = await fetchFn(url).catch(() => null);
  if (!res?.ok) return null;
  const text = await res.text();
  if (!/^\s*[[{]/.test(text)) return null;
  return parseJSON(text, url);
}

// loadReport is the report at url, or null when there is none.
export async function loadReport(url: string, fetchFn: typeof fetch = fetch): Promise<Report | null> {
  const raw = await optionalJSON(url, fetchFn);
  if (raw === null) return null;
  if (!isReport(raw)) throw new Error(`${url}: a registry, not a report`);
  return parseReport(raw);
}

// loadIndex is the charters the page can switch between, or none.
export async function loadIndex(url: string = INDEX_URL, fetchFn: typeof fetch = fetch): Promise<CharterLink[]> {
  const raw = await optionalJSON(url, fetchFn);
  return raw === null ? [] : parseIndex(raw);
}
