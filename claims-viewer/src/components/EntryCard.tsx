import { useState } from "react";
import { coverage, type Check } from "../covered";
import type { Entry } from "../model";
import { entryHref } from "../route";
import { AreaLink, FileLink, IdLink, InlineText, KindTag, MarkBadge, Prose, StatusBadge, useView, VerdictBadge } from "./bits";

interface Props {
  e: Entry;
  // Show the area the entry is in: off on the area's own page.
  showArea?: boolean;
  open?: boolean;
}

// EntryCard is an entry as its document writes it, in the same order, with
// what the gate knows beside each part: the code that realizes it, and for
// each scenario the tests that cite it or the waiver that accepts none.
export function EntryCard({ e, showArea = true, open: initiallyOpen = true }: Props) {
  const { reg, report } = useView();
  const [open, setOpen] = useState(initiallyOpen);
  const cov = coverage(reg, e, report);
  const servedBy = reg.servedBy.get(e.id) ?? [];
  const seenBy = (reg.seenBy.get(e.id) ?? []).filter((x) => !e.see.includes(x.id));
  const supportedBy = reg.supportedBy.get(e.id) ?? [];
  return (
    <article className={`card status-${e.status}`} id={e.id}>
      <header className="card-head" onClick={() => setOpen(!open)}>
        <span className="caret">{open ? "▾" : "▸"}</span>
        <a className="card-id" href={entryHref(e.id)} onClick={(ev) => ev.stopPropagation()}>
          {e.id}
        </a>
        <span className="card-title">{e.title}</span>
        <span className="card-tags">
          {e.need && <span className="kind kind-need">need</span>}
          {e.level === "observation" && <span className="kind">observation</span>}
          <KindTag kind={e.kind} />
          <StatusBadge status={e.status} />
          <VerdictBadge verdict={cov.verdict} />
        </span>
      </header>
      {open && (
        <div className="card-body">
          <dl className="fields">
            {showArea && (
              <>
                <dt>Area</dt>
                <dd>
                  <AreaLink slug={e.area} /> <FileLink file={e.file} line={e.line} />
                </dd>
              </>
            )}
            {!showArea && (
              <>
                <dt>At</dt>
                <dd>
                  <FileLink file={e.file} line={e.line} />
                </dd>
              </>
            )}
            {e.serves.length > 0 && <Ids label="Serves" ids={e.serves} />}
            {e.see.length > 0 && <Ids label="See" ids={e.see} />}
            {e.supports.length > 0 && <Ids label="Supports" ids={e.supports} />}
            {e.accounts.length > 0 && (
              <>
                <dt>Accounts</dt>
                <dd>
                  {e.accounts.map((a) => (
                    <code key={a} className="surface">
                      {a}
                    </code>
                  ))}
                </dd>
              </>
            )}
            {e.level === "observation" && (
              <>
                <dt>Subject</dt>
                <dd>
                  <code>{e.subject}</code> {e.observed ? `observed ${e.observed}` : <span className="muted">never observed</span>}
                </dd>
                {e.source && (
                  <>
                    <dt>Source</dt>
                    <dd>
                      <a href={e.source} target="_blank" rel="noreferrer">
                        {e.source}
                      </a>
                    </dd>
                  </>
                )}
              </>
            )}
          </dl>
          {e.text && <Prose text={e.text} />}
          {e.why && <Prose text={e.why} label="Why" />}
          {e.verification && <Prose text={e.verification} label="Verification" />}
          {e.level === "spec" && <Scenarios e={e} checks={cov.checks} />}
          {cov.bare.length > 0 && (
            <section>
              <h4>Tests citing the entry, no scenario</h4>
              <Tests refs={cov.bare} />
            </section>
          )}
          {e.guidance && <Prose text={e.guidance} label="Guidance" />}
          {e.changes.length > 0 && (
            <section>
              <h4>Changes</h4>
              <ul className="changes">
                {e.changes.map((c) => (
                  <li key={c.line}>
                    <span className="date">{c.date}</span> <InlineText text={c.text} />
                  </li>
                ))}
              </ul>
            </section>
          )}
          {e.level === "spec" && (
            <section>
              <h4>
                Code{" "}
                {cov.realized === false && e.status === "active" && <span className="x">no code realizes it</span>}
                {cov.realized === null && <span className="muted small">a need is realized by the decisions that serve it</span>}
              </h4>
              {e.anchors.length > 0 && (
                <ul className="sites">
                  {e.anchors.map((a) => (
                    <li key={`${a.file}:${a.line}`}>
                      <span className={`anchor-kind ${a.kind}`}>{a.kind}</span> <code className="decl">{a.decl}</code>{" "}
                      <FileLink file={a.file} line={a.line} />
                    </li>
                  ))}
                </ul>
              )}
            </section>
          )}
          {(servedBy.length > 0 || seenBy.length > 0 || supportedBy.length > 0) && (
            <section className="backlinks">
              {servedBy.length > 0 && <Back label="Served by" entries={servedBy} />}
              {seenBy.length > 0 && <Back label="Seen by" entries={seenBy} />}
              {supportedBy.length > 0 && <Back label="Observations" entries={supportedBy} />}
            </section>
          )}
        </div>
      )}
    </article>
  );
}

function Ids({ label, ids }: { label: string; ids: string[] }) {
  return (
    <>
      <dt>{label}</dt>
      <dd className="ids">
        {ids.map((id) => (
          <IdLink key={id} id={id} />
        ))}
      </dd>
    </>
  );
}

function Back({ label, entries }: { label: string; entries: Entry[] }) {
  return (
    <div>
      <span className="muted">{label}</span>{" "}
      <span className="ids">
        {entries.map((x) => (
          <IdLink key={x.id} id={x.id} />
        ))}
      </span>
    </div>
  );
}

function Scenarios({ e, checks }: { e: Entry; checks: Check[] }) {
  if (e.scenarios.length === 0) {
    const c = checks[0];
    return (
      <section>
        <h4>
          Tests <MarkBadge mark={c.mark} />
        </h4>
        <Tests refs={c.refs} check={c} />
        {c.waiver && <WaiverNote check={c} />}
      </section>
    );
  }
  return (
    <section>
      <h4>Scenarios</h4>
      <ol className="scenarios">
        {e.scenarios.map((s, i) => {
          const c = checks[i];
          return (
            <li key={s.handle} id={`${e.id}/${s.handle}`} className={`scenario mark-row-${c.mark}`}>
              <div className="scenario-head">
                <code className="handle">{s.handle}</code>
                <MarkBadge mark={c.mark} />
              </div>
              <div className="scenario-text">
                <InlineText text={s.text} />
              </div>
              <Tests refs={c.refs} check={c} />
              {c.waiver && <WaiverNote check={c} />}
            </li>
          );
        })}
      </ol>
    </section>
  );
}

function WaiverNote({ check }: { check: Check }) {
  const w = check.waiver!;
  return (
    <div className="waiver">
      <span className="muted">waived by</span> <code>{w.pattern}</code> <FileLink file={w.file} line={w.line} />: {w.reason}
    </div>
  );
}

function Tests({ refs, check }: { refs: { file: string; func: string; line: number }[]; check?: Check }) {
  if (refs.length === 0) return null;
  const outcome = (func: string) => check?.tests.find((t) => t.test === func || t.test.startsWith(func + "/"))?.outcome;
  return (
    <ul className="sites tests">
      {refs.map((r) => {
        const o = outcome(r.func);
        return (
          <li key={`${r.file}:${r.line}:${r.func}`}>
            {o && <span className={`outcome outcome-${o}`}>{o}</span>}
            <code className="decl">{r.func || "(test)"}</code> <FileLink file={r.file} line={r.line} />
          </li>
        );
      })}
    </ul>
  );
}
