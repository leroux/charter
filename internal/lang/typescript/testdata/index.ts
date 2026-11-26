// Code generated from claims by charter gen. DO NOT EDIT.

import { appendFileSync } from "node:fs";

type Finished = (fn: (ctx: { task: { name: string; file?: { filepath?: string }; result?: { state?: string } } }) => void) => void;

function record(id: string, handles: string[]): void {
  const path = (globalThis as { process?: { env?: Record<string, string | undefined> } }).process?.env?.["CHARTER_EVENTS"];
  const onTestFinished = (globalThis as { onTestFinished?: Finished }).onTestFinished;
  if (!path || typeof onTestFinished !== "function") return;
  onTestFinished((ctx) => {
    const state = ctx.task.result?.state;
    const outcome = state === "pass" ? "passed" : state === "skip" ? "skipped" : "failed";
    for (const handle of handles.length ? handles : [""]) {
      const line: Record<string, string> = { id, test: ctx.task.name, file: ctx.task.file?.filepath ?? "", outcome };
      if (handle) line.handle = handle;
      appendFileSync(path, JSON.stringify(line) + "\n");
    }
  });
}

/** Opens a test that verifies an entry, naming by handle the scenarios it exercises. */
export function verifies(id: string, ...handles: string[]): void {
  if (typeof id !== "string" || id === "") throw new TypeError("claims: verifies needs a spec");
  record(id, handles);
}

/** Opens a probe that checks an observation against the live system. */
export function observes(name: string): void {
  if (typeof name !== "string" || name === "") throw new TypeError("claims: observes needs an observation");
  record(name, []);
}

/** The record's entries and observations, each symbol mapped to its id or name. */
export const claims = {
  /** The bucket refills at the rate *\/ */
  SW_010: "SW-010",
  /** A retired rule */
  /** @deprecated SW-012 is retired; see its Changes in claims/refill.md. */
  SW_012: "SW-012",
} as const;
