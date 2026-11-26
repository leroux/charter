import { createContext, useContext, type ReactNode } from "react";
import { MARK_LABEL, VERDICT_LABEL, type Mark, type Tally, type Verdict } from "../covered";
import { statusGlyph, statusTone, type Entry, type Kind, type Registry, type Status } from "../model";
import { blocks, inline, type Inline } from "../prose";
import type { Report } from "../report";
import type { Tour } from "../tour";
import { areaHref, entryHref } from "../route";

export interface View {
  reg: Registry;
  report: Report | null;
  tour: Tour;
}

export const ViewContext = createContext<View | null>(null);

export function useView(): View {
  return useContext(ViewContext)!;
}

export function StatusBadge({ status }: { status: Status }) {
  return (
    <span className={`badge tone-${statusTone(status)}`} title={status}>
      {statusGlyph(status)} {status}
    </span>
  );
}

export function KindTag({ kind }: { kind: Kind }) {
  return kind ? <span className={`kind kind-${kind}`}>{kind}</span> : null;
}

export function MarkBadge({ mark }: { mark: Mark }) {
  return <span className={`mark mark-${mark}`}>{MARK_LABEL[mark]}</span>;
}

export function VerdictBadge({ verdict }: { verdict: Verdict }) {
  return verdict === "none" ? null : <span className={`verdict verdict-${verdict}`}>{VERDICT_LABEL[verdict]}</span>;
}

// IdLink names an entry by id and shows its title on hover; a deprecated
// entry is struck through.
export function IdLink({ id, title = false }: { id: string; title?: boolean }) {
  const { reg } = useView();
  const e = reg.byId.get(id);
  if (!e) return <code>{id}</code>;
  return (
    <a className={`idlink status-${e.status}`} href={entryHref(id)} title={`${e.title} (${e.status})`}>
      <span className="idlink-id">{id}</span>
      {title && <span className="idlink-title"> {e.title}</span>}
    </a>
  );
}

export function AreaLink({ slug }: { slug: string }) {
  const { reg } = useView();
  const a = reg.areaBySlug.get(slug);
  return <a href={areaHref(slug)}>{a?.title ?? slug}</a>;
}

// FileLink opens file:line in the editor when the export named the
// checkout's root.
export function FileLink({ file, line, children }: { file: string; line?: number; children?: ReactNode }) {
  const { reg } = useView();
  const at = line ? `${file}:${line}` : file;
  const label = children ?? <code>{at}</code>;
  if (!reg.root) return <span className="file">{label}</span>;
  return (
    <a className="file" href={`vscode://file${reg.root}/${at}`} title={`open ${at} in the editor`}>
      {label}
    </a>
  );
}

function Inlines({ parts }: { parts: Inline[] }) {
  return (
    <>
      {parts.map((p, i) =>
        p.t === "code" ? <code key={i}>{p.v}</code> : p.t === "id" ? <IdLink key={i} id={p.v} /> : <span key={i}>{p.v}</span>,
      )}
    </>
  );
}

export function Prose({ text, label }: { text: string; label?: string }) {
  const { reg } = useView();
  const isId = (id: string) => reg.byId.has(id);
  const bs = blocks(text, isId);
  return (
    <div className="prose">
      {bs.map((b, i) =>
        b.t === "p" ? (
          <p key={i}>
            {i === 0 && label && <strong>{label}: </strong>}
            <Inlines parts={b.v} />
          </p>
        ) : (
          <ul key={i}>
            {b.v.map((item, j) => (
              <li key={j}>
                <Inlines parts={item} />
              </li>
            ))}
          </ul>
        ),
      )}
    </div>
  );
}

export function InlineText({ text }: { text: string }) {
  const { reg } = useView();
  return <Inlines parts={inline(text.replace(/\s+/g, " "), (id) => reg.byId.has(id))} />;
}

// CoverageBar is a tally's active checks as one stacked bar: verified,
// waived, skipped, untested and failing.
export function CoverageBar({ t, wide = false }: { t: Tally; wide?: boolean }) {
  const m = t.marks;
  const parts: [string, number, string][] = [
    ["ok", m.pass + m.cited, `${m.pass + m.cited} tested${m.pass ? ` (${m.pass} passed)` : ""}`],
    ["waived", m.waived, `${m.waived} waived`],
    ["skip", m.skip, `${m.skip} skipped`],
    ["missing", m.missing, `${m.missing} untested`],
    ["fail", m.fail, `${m.fail} failed`],
  ];
  const total = parts.reduce((n, p) => n + p[1], 0);
  const unrealized = t.decisions - t.realized;
  const title = [...parts.filter((p) => p[1] > 0).map((p) => p[2]), unrealized ? `${unrealized} decisions unrealized` : ""]
    .filter(Boolean)
    .join(", ");
  if (total === 0) return <span className="muted small">{t.active ? "no scenarios" : "nothing active"}</span>;
  return (
    <span className={`covbar ${wide ? "wide" : ""}`} title={title}>
      <span className="covbar-track">
        {parts.map(([k, n]) => (n > 0 ? <span key={k} className={`seg-${k}`} style={{ flexGrow: n }} /> : null))}
      </span>
      <span className="covbar-label">
        {m.pass + m.cited}/{total}
        {m.waived > 0 && <span className="w"> · {m.waived}w</span>}
        {m.missing + m.fail > 0 && <span className="x"> · {m.missing + m.fail}✕</span>}
        {unrealized > 0 && <span className="x"> · {unrealized} unrealized</span>}
      </span>
    </span>
  );
}

export function EntryLine({ e, children }: { e: Entry; children?: ReactNode }) {
  return (
    <span className="entry-line">
      <span className={`glyph tone-${statusTone(e.status)}`} title={e.status}>
        {statusGlyph(e.status)}
      </span>
      <IdLink id={e.id} />
      <span className={`entry-line-title status-${e.status}`}>{e.title}</span>
      {children && <span className="entry-meta">{children}</span>}
    </span>
  );
}
