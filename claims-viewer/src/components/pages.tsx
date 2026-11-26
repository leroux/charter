import { Fragment, useMemo, useState, type MouseEvent, type ReactElement } from "react";
import { coverage, tally, type Tally } from "../covered";
import { gateCount, gaps } from "../gaps";
import { isDecision, needsOf, servingTree, type Area, type Entry, type TreeNode } from "../model";
import { search } from "../search";
import { areaHref, entryHref, href } from "../route";
import { areaNumber, areaOrder, firstSentence, needsServedBy, neighbours, servers } from "../tour";
import { AreaLink, CoverageBar, EntryLine, FileLink, InlineText, Prose, StatusBadge, useView, VerdictBadge } from "./bits";
import { EntryCard } from "./EntryCard";

// useAreaNumber gives an area's place in the reading order, from 1, or "".
function useAreaNumber(): (slug: string) => string {
  const { tour } = useView();
  return (slug) => areaNumber(tour, slug);
}

function liveDecisions(entries: Entry[]): number {
  return entries.filter((e) => isDecision(e) && e.status !== "deprecated").length;
}

// The walkthrough is the home page: what a charter is and how to read an
// entry, then the areas in the charter's reading order, each with the
// sentence that opens its document and the needs it does most to deliver.
export function Walkthrough() {
  const { reg, report, tour } = useView();
  const number = useAreaNumber();
  const needsArea = reg.areas.find((a) => a.kind === "needs");
  const observations = reg.areas.filter((a) => a.kind === "observations");
  const first = tour.path[0];
  const unserved = tour.unserved.filter((n) => n.status === "active");
  return (
    <div className="page walkthrough">
      <h1>{needsArea?.title ?? "The charter"}</h1>
      <p className="muted small">{reg.source}</p>
      <p className="lede">
        What people can rely on, and the decisions that deliver it, each tied to the code that realizes it and the tests
        that verify it. Read the areas in the order the charter gives them; each delivers the needs listed under it.
      </p>
      {needsArea?.summary && <Prose text={needsArea.summary} />}

      <div className="howto">
        <div>
          <span className="howto-n">1</span>
          <strong>A need</strong> is what a person can rely on, verified from their side.
        </div>
        <div>
          <span className="howto-n">2</span>
          <strong>A decision</strong> says how the system behaves at one boundary, and serves needs.
        </div>
        <div>
          <span className="howto-n">3</span>
          <strong>Scenarios</strong> make a decision checkable: each is cited by a test, or waived with a reason.
        </div>
        <div>
          <span className="howto-n">4</span>
          <strong>Code</strong> names the decision it realizes, so a decision leads to its code and back.
        </div>
      </div>

      {first && (
        <p className="start">
          <a className="button primary" href={areaHref(first.slug)}>
            Start reading: {number(first.slug)}. {first.title} →
          </a>{" "}
          <a className="button" href={href({ page: "needs" })}>
            Or begin with what people need
          </a>
        </p>
      )}

      <h2>The areas, in reading order</h2>
      <ol className="chapter-areas">
        {tour.path.map((a) => (
          <AreaItem key={a.slug} a={a} />
        ))}
      </ol>

      {unserved.length > 0 && (
        <>
          <h2>Needs no decision delivers yet</h2>
          <ul className="plain-list">
            {unserved.map((n) => (
              <li key={n.id}>
                <EntryLine e={n} />
              </li>
            ))}
          </ul>
        </>
      )}

      {observations.length > 0 && (
        <>
          <h2>What it rests on</h2>
          <p className="muted small">Facts about external systems the decisions rely on, each with its source, checked by probes.</p>
          <ul className="chapter-areas">
            {observations.map((a) => (
              <li key={a.slug}>
                <a href={areaHref(a.slug)}>{a.title}</a> <span className="muted small">{reg.inArea.get(a.slug)!.length} observations</span>
              </li>
            ))}
          </ul>
        </>
      )}
      {report === null && <p className="muted small">No test report is loaded, so a scenario shows as tested when a test cites it, not whether it passed.</p>}
    </div>
  );
}

function AreaItem({ a }: { a: Area }) {
  const { reg, report, tour } = useView();
  const number = useAreaNumber();
  const es = reg.inArea.get(a.slug)!;
  const needs = tour.needsAt.get(a.slug)!.filter((n) => n.status !== "deprecated");
  const shown = needs.slice(0, 3);
  return (
    <li>
      <div className="area-row">
        <span className="area-n">{number(a.slug)}</span>
        <a className="area-title" href={areaHref(a.slug)}>
          {a.title}
        </a>
        <span className="muted small">{liveDecisions(es)} decisions</span>
        <span className="spacer" />
        <CoverageBar t={tally(reg, es, report)} />
      </div>
      {a.summary && <p className="area-intro">{firstSentence(a.summary)}</p>}
      {needs.length > 0 && (
        <ul className="area-needs">
          {shown.map((n) => (
            <li key={n.id}>
              <EntryLine e={n} />
            </li>
          ))}
          {needs.length > shown.length && (
            <li className="muted small">
              <a href={href({ page: "needs" })}>and {needs.length - shown.length} more</a>
            </li>
          )}
        </ul>
      )}
    </li>
  );
}

// NeedsPage lists the needs under the area that does most to deliver each,
// in the reading order, so a reader can start from the why.
export function NeedsPage() {
  const { reg, tour } = useView();
  const number = useAreaNumber();
  const [deprecated, setDeprecated] = useState(false);
  const needsArea = reg.areas.find((a) => a.kind === "needs");
  const retired = reg.entries.filter((e) => e.need && e.status === "deprecated");
  const keep = (n: Entry) => deprecated || n.status !== "deprecated";
  const unserved = tour.unserved.filter(keep);
  return (
    <div className="page">
      <p className="crumbs">
        <a href={href({ page: "overview" })}>Walkthrough</a>
      </p>
      <h1>{needsArea?.title ?? "Needs"}</h1>
      {needsArea?.summary && <Prose text={needsArea.summary} />}
      <p className="muted small">
        Each need is under the area whose decisions most often reach it, the earlier in the reading order on a tie.{" "}
        {needsArea && <FileLink file={needsArea.file} />}
      </p>
      {retired.length > 0 && (
        <label className="small">
          <input type="checkbox" checked={deprecated} onChange={(ev) => setDeprecated(ev.target.checked)} /> show {retired.length} retired needs
        </label>
      )}
      {tour.path.map((a) => {
        const ns = tour.needsAt.get(a.slug)!.filter(keep);
        if (ns.length === 0) return null;
        return (
          <section key={a.slug} className="chapter">
            <h2>
              <span className="chap-n">{number(a.slug)}</span> <a href={areaHref(a.slug)}>{a.title}</a>
            </h2>
            {ns.map((n) => (
              <NeedSummary key={n.id} need={n} />
            ))}
          </section>
        );
      })}
      {unserved.length > 0 && (
        <section className="chapter">
          <h2>
            <span className="chap-n">·</span> Delivered by no live decision
          </h2>
          {unserved.map((n) => (
            <NeedSummary key={n.id} need={n} />
          ))}
        </section>
      )}
    </div>
  );
}

function NeedSummary({ need }: { need: Entry }) {
  const { reg, report, tour } = useView();
  const number = useAreaNumber();
  const byArea = new Map<string, number>();
  for (const d of servers(reg, need)) byArea.set(d.area, (byArea.get(d.area) ?? 0) + 1);
  const order = areaOrder(tour);
  const areas = [...byArea].sort((a, b) => b[1] - a[1] || order(a[0], b[0]));
  return (
    <div className={`need-summary status-${need.status}`}>
      <EntryLine e={need}>
        <VerdictBadge verdict={coverage(reg, need, report).verdict} />
      </EntryLine>
      {need.text && <p className="need-text">{firstSentence(need.text)}</p>}
      {areas.length > 0 && (
        <div className="small need-areas">
          <span className="muted">delivered in </span>
          {areas.map(([slug, n], i) => (
            <Fragment key={slug}>
              {i > 0 && ", "}
              <a href={areaHref(slug)}>
                {number(slug)}. {reg.areaBySlug.get(slug)?.title}
              </a>{" "}
              <span className="muted">({n})</span>
            </Fragment>
          ))}
        </div>
      )}
    </div>
  );
}

function Stat({ n, label, tone }: { n: number | string; label: string; tone?: string }) {
  return (
    <div className={`stat ${tone ? `tone-${tone}` : ""}`}>
      <div className="stat-n">{n}</div>
      <div className="stat-label">{label}</div>
    </div>
  );
}

// CoveragePage is the charter in numbers: the totals, then each area and
// each need with its coverage, in the reading order.
export function CoveragePage() {
  const { reg, report, tour } = useView();
  const all = useMemo(() => tally(reg, reg.entries, report), [reg, report]);
  const sections = useMemo(() => gaps(reg, report), [reg, report]);
  const m = all.marks;
  const checks = m.pass + m.cited + m.waived + m.skip + m.missing + m.fail;
  const failing = gateCount(sections) + reg.problems.length;
  const number = useAreaNumber();
  const needs = [...tour.path.flatMap((a) => tour.needsAt.get(a.slug)!), ...tour.unserved].filter((n) => n.status !== "deprecated");
  return (
    <div className="page">
      <h1>Coverage</h1>
      <p className="muted">
        {reg.source}
        {report ? " · with a test report" : " · no test report: a scenario shows as tested when a test cites it, not whether it passed"}
      </p>
      <div className="stats">
        <Stat n={reg.entries.filter((e) => e.need && e.status === "active").length} label="active needs" />
        <Stat n={reg.entries.filter((e) => isDecision(e) && e.status === "active").length} label="active decisions" />
        <Stat n={all.planned} label="planned" />
        <Stat n={`${all.realized}/${all.decisions}`} label="decisions realized in code" tone={all.realized < all.decisions ? "critical" : "good"} />
        <Stat n={`${m.pass + m.cited}/${checks}`} label={report ? `scenarios tested (${m.pass} passed)` : "scenarios tested"} tone={m.missing + m.fail ? "critical" : "good"} />
        <Stat n={m.waived} label="waived" tone={m.waived ? "warning" : undefined} />
        <a className="stat-link" href={href({ page: "gaps" })}>
          <Stat n={failing} label="gate failures" tone={failing ? "critical" : "good"} />
        </a>
      </div>

      <h2>Areas</h2>
      <table className="grid">
        <thead>
          <tr>
            <th>Area</th>
            <th className="num">Active</th>
            <th className="num">Planned</th>
            <th className="num">Deprecated</th>
            <th>Coverage</th>
          </tr>
        </thead>
        <tbody>
          {tour.path.map((a) => (
            <AreaRow key={a.slug} a={a} n={number(a.slug)} />
          ))}
        </tbody>
      </table>

      <h2>Needs</h2>
      <table className="grid">
        <thead>
          <tr>
            <th>Need</th>
            <th>Delivered mostly in</th>
            <th className="num">Decisions</th>
            <th>Own scenarios</th>
            <th>Decisions' coverage</th>
          </tr>
        </thead>
        <tbody>
          {needs.map((n) => (
            <NeedRow key={n.id} need={n} />
          ))}
        </tbody>
      </table>
    </div>
  );
}

function AreaRow({ a, n }: { a: Area; n: string }) {
  const { reg, report } = useView();
  const t = tally(reg, reg.inArea.get(a.slug)!, report);
  return (
    <tr>
      <td>
        <span className="area-n">{n}</span> <a href={areaHref(a.slug)}>{a.title}</a>
      </td>
      <td className="num">{t.active}</td>
      <td className="num">{t.planned || ""}</td>
      <td className="num muted">{t.deprecated || ""}</td>
      <td>
        <CoverageBar t={t} />
      </td>
    </tr>
  );
}

function NeedRow({ need }: { need: Entry }) {
  const { reg, report, tour } = useView();
  const number = useAreaNumber();
  const serving = useMemo(() => servers(reg, need), [reg, need]);
  const own = tally(reg, [need], report);
  const theirs = tally(reg, serving, report);
  const home = tour.homeOf.get(need.id) ?? "";
  return (
    <tr className={`status-${need.status}`}>
      <td>
        <EntryLine e={need} />
      </td>
      <td className="small">
        {home ? (
          <a href={areaHref(home)}>
            {number(home)}. {reg.areaBySlug.get(home)?.title}
          </a>
        ) : (
          <span className="x">none</span>
        )}
      </td>
      <td className="num">{serving.length || <span className="x">0</span>}</td>
      <td>{need.status === "active" ? <CoverageBar t={own} /> : <StatusBadge status={need.status} />}</td>
      <td>
        <CoverageBar t={theirs} />
      </td>
    </tr>
  );
}

export function AreaPage({ slug }: { slug: string }) {
  const { reg, report, tour } = useView();
  const area = reg.areaBySlug.get(slug);
  const [deprecated, setDeprecated] = useState(false);
  const [open, setOpen] = useState(true);
  const [version, setVersion] = useState(0);
  const number = useAreaNumber();
  if (!area) return <NotFound what={`area ${slug}`} />;
  const entries = reg.inArea.get(slug)!;
  const shown = entries.filter((e) => deprecated || e.status !== "deprecated");
  const retired = entries.filter((e) => e.status === "deprecated").length;
  const t = tally(reg, entries, report);
  const placed = tour.index.has(slug);
  const served = area.kind === "decisions" ? needsServedBy(reg, slug) : [];
  return (
    <div className="page">
      <p className="crumbs">
        <a href={href({ page: "overview" })}>Walkthrough</a>
        {placed && (
          <>
            <span>›</span>
            <span>
              area {number(slug)} of {tour.path.length}
            </span>
          </>
        )}
        <span className="spacer" />
        {placed && <PrevNext slug={slug} compact />}
      </p>
      <h1>
        {placed && <span className="area-n big">{number(slug)}</span>} {area.title}
      </h1>
      {area.summary && (
        <div className="intro">
          <Prose text={area.summary} />
        </div>
      )}
      <p className="muted small">
        <FileLink file={area.file} />
        {area.code.length > 0 && (
          <>
            {" "}
            · code{" "}
            {area.code.map((c) => (
              <code key={c} className="surface">
                {c}
              </code>
            ))}
          </>
        )}
      </p>
      {served.length > 0 && (
        <div className="serves-box">
          <span className="muted small">Its decisions deliver</span>
          <ul>
            {served.slice(0, 6).map(({ need, count }) => (
              <li key={need.id}>
                <EntryLine e={need}>
                  <span className="muted">{count} {count === 1 ? "decision" : "decisions"}</span>
                </EntryLine>
              </li>
            ))}
          </ul>
          {served.length > 6 && <span className="muted small">and {served.length - 6} more needs, in part</span>}
        </div>
      )}
      <div className="toolbar">
        <CoverageBar t={t} wide />
        <span className="spacer" />
        {retired > 0 && (
          <label>
            <input type="checkbox" checked={deprecated} onChange={(ev) => setDeprecated(ev.target.checked)} /> show {retired} deprecated
          </label>
        )}
        <button
          onClick={() => {
            setOpen(!open);
            setVersion(version + 1);
          }}
        >
          {open ? "collapse all" : "expand all"}
        </button>
      </div>
      <nav className="toc">
        {shown.map((e) => (
          <a key={e.id} href={`#${e.id}`} onClick={(ev) => scrollTo(ev, e.id)} className={`toc-item status-${e.status}`} title={e.title}>
            <VerdictDot e={e} /> {e.id}
          </a>
        ))}
      </nav>
      {shown.map((e) => (
        <EntryCard key={`${e.id}-${version}`} e={e} showArea={false} open={open} />
      ))}
      {placed && <PrevNext slug={slug} />}
    </div>
  );
}

function PrevNext({ slug, compact = false }: { slug: string; compact?: boolean }) {
  const { tour } = useView();
  const number = useAreaNumber();
  const { prev, next } = neighbours(tour, slug);
  if (compact) {
    return (
      <span className="prevnext-compact">
        {prev ? (
          <a href={areaHref(prev.slug)} title={prev.title}>
            ← {number(prev.slug)}
          </a>
        ) : (
          <span className="muted">←</span>
        )}
        {next ? (
          <a href={areaHref(next.slug)} title={next.title}>
            {number(next.slug)} →
          </a>
        ) : (
          <span className="muted">→</span>
        )}
      </span>
    );
  }
  return (
    <nav className="prevnext">
      {prev ? (
        <a className="pn prev" href={areaHref(prev.slug)}>
          <span className="muted small">← previous</span>
          <span>
            {number(prev.slug)}. {prev.title}
          </span>
        </a>
      ) : (
        <a className="pn prev" href={href({ page: "overview" })}>
          <span className="muted small">← back to</span>
          <span>the walkthrough</span>
        </a>
      )}
      {next ? (
        <a className="pn next" href={areaHref(next.slug)}>
          <span className="muted small">next →</span>
          <span>
            {number(next.slug)}. {next.title}
          </span>
        </a>
      ) : (
        <a className="pn next" href={href({ page: "coverage" })}>
          <span className="muted small">the end →</span>
          <span>coverage in numbers</span>
        </a>
      )}
    </nav>
  );
}

// A link to an entry on the same page scrolls to it, leaving the route.
function scrollTo(ev: MouseEvent, id: string) {
  ev.preventDefault();
  document.getElementById(id)?.scrollIntoView({ behavior: "smooth", block: "start" });
}

function VerdictDot({ e }: { e: Entry }) {
  const { reg, report } = useView();
  const v = coverage(reg, e, report).verdict;
  return <span className={`dot verdict-${v} status-${e.status}`} title={v === "none" ? e.status : v} />;
}

export function EntryPage({ id }: { id: string }) {
  const { reg, tour } = useView();
  const number = useAreaNumber();
  const e = reg.byId.get(id);
  if (!e) return <NotFound what={id} />;
  const needs = e.need ? [] : needsOf(reg, e);
  const tree = servingTree(reg, e, () => true, areaOrder(tour));
  return (
    <div className="page">
      <p className="crumbs">
        <a href={href({ page: "overview" })}>Walkthrough</a>
        <span>›</span>
        {e.need ? (
          <a href={href({ page: "needs" })}>What people need</a>
        ) : (
          <>
            <a href={areaHref(e.area)}>
              {number(e.area)}. {reg.areaBySlug.get(e.area)?.title}
            </a>
          </>
        )}
      </p>
      {needs.length > 0 && (
        <div className="serves-box">
          <span className="muted small">Delivers</span>
          <ul>
            {needs.map((n) => (
              <li key={n.id}>
                <EntryLine e={n} />
              </li>
            ))}
          </ul>
        </div>
      )}
      {!e.need && e.serves.length > 0 && <UpChain e={e} />}
      <EntryCard e={e} />
      {tree.length > 0 && (
        <section>
          <h2>What serves it</h2>
          <p className="muted small">The decisions that deliver it, by area in the walkthrough's order; a decision under another serves that one.</p>
          <ServingTree nodes={tree} />
        </section>
      )}
    </div>
  );
}

// UpChain shows how a decision reaches its needs when it does so through
// other decisions.
function UpChain({ e }: { e: Entry }) {
  const { reg } = useView();
  const viaDecisions = e.serves.some((id) => !reg.byId.get(id)!.need);
  if (!viaDecisions) return null;
  const chain = (x: Entry, seen: Set<string>): ReactElement[] =>
    x.serves.map((id) => {
      const up = reg.byId.get(id)!;
      return (
        <li key={id}>
          <EntryLine e={up} />
          {!up.need && !seen.has(id) && <ul className="tree">{chain(up, new Set([...seen, id]))}</ul>}
        </li>
      );
    });
  return (
    <section className="upchain">
      <h4>Serves, up to the needs</h4>
      <ul className="tree">{chain(e, new Set([e.id]))}</ul>
    </section>
  );
}

function ServingTree({ nodes }: { nodes: TreeNode[] }) {
  const { reg, report } = useView();
  const number = useAreaNumber();
  const [deprecated, setDeprecated] = useState(false);
  // A nested decision names its area when it is not its parent's.
  const render = (ns: TreeNode[], top: boolean, parentArea = ""): ReactElement => {
    const kept = ns.filter((n) => deprecated || n.entry.status !== "deprecated");
    return (
      <ul className="tree">
        {kept.map((n, i) => (
          <Fragment key={n.entry.id}>
            {top && (i === 0 || kept[i - 1].entry.area !== n.entry.area) && (
              <li className="tree-area">
                <a href={areaHref(n.entry.area)}>
                  {number(n.entry.area)}. {reg.areaBySlug.get(n.entry.area)?.title}
                </a>
              </li>
            )}
            <li>
              <EntryLine e={n.entry}>
                {!top && n.entry.area !== parentArea && <AreaLink slug={n.entry.area} />}
                <VerdictBadge verdict={coverage(reg, n.entry, report).verdict} />
                {n.repeat && <span className="muted">(above)</span>}
              </EntryLine>
              {n.children.length > 0 && render(n.children, false, n.entry.area)}
            </li>
          </Fragment>
        ))}
      </ul>
    );
  };
  const anyDeprecated = flattenAll(nodes).some((x) => x.status === "deprecated");
  return (
    <>
      {anyDeprecated && (
        <label className="small">
          <input type="checkbox" checked={deprecated} onChange={(ev) => setDeprecated(ev.target.checked)} /> show deprecated
        </label>
      )}
      {render(nodes, true)}
    </>
  );
}

function flattenAll(nodes: TreeNode[]): Entry[] {
  return nodes.flatMap((n) => [n.entry, ...flattenAll(n.children)]);
}

export function GapsPage() {
  const { reg, report } = useView();
  const sections = useMemo(() => gaps(reg, report), [reg, report]);
  return (
    <div className="page">
      <h1>Gaps</h1>
      <p className="muted">
        What the gate fails on first, then what it lets through that a reviewer would want to know, each by area in the walkthrough's order.
        {!report && " Load a report to see failing tests."}
      </p>
      {reg.problems.length > 0 && (
        <section className="gap-section">
          <h2>
            Problems loading <span className="count x">{reg.problems.length}</span> <span className="gate-tag">gate</span>
          </h2>
          <ul className="problems">
            {reg.problems.map((p) => (
              <li key={p}>
                <InlineText text={p} />
              </li>
            ))}
          </ul>
        </section>
      )}
      <nav className="toc">
        {sections.map((s) => (
          <a key={s.key} href={`#gap-${s.key}`} onClick={(ev) => scrollTo(ev, `gap-${s.key}`)} className="toc-item toc-text">
            {s.title} <span className={s.gate && s.items.length ? "x" : "muted"}>{s.items.length}</span>
          </a>
        ))}
      </nav>
      {sections.map((s) => (
        <section key={s.key} className="gap-section" id={`gap-${s.key}`}>
          <h2>
            {s.title} <span className={`count ${s.gate && s.items.length ? "x" : ""}`}>{s.items.length}</span>
            {s.gate && <span className="gate-tag">gate</span>}
          </h2>
          <p className="muted small">{s.hint}</p>
          {s.items.length === 0 ? <p className="muted">None.</p> : <GapList items={s.items} />}
        </section>
      ))}
    </div>
  );
}

function GapList({ items }: { items: { entry: Entry; handle: string; note: string }[] }) {
  const { tour } = useView();
  const number = useAreaNumber();
  // Grouped by area, as a person fixing them works one area at a time.
  const byArea = new Map<string, typeof items>();
  for (const i of items) byArea.set(i.entry.area, [...(byArea.get(i.entry.area) ?? []), i]);
  const order = areaOrder(tour);
  return (
    <dl className="gap-list">
      {[...byArea]
        .sort((a, b) => order(a[0], b[0]))
        .map(([slug, its]) => (
          <Fragment key={slug}>
            <dt>
              <span className="area-n">{number(slug)}</span> <AreaLink slug={slug} /> <span className="muted small">{its.length}</span>
            </dt>
            <dd>
              <ul>
                {its.map((i) => (
                  <li key={`${i.entry.id}/${i.handle}`}>
                    <EntryLine e={i.entry}>
                      {i.handle && (
                        <a className="handle" href={entryHref(i.entry.id)}>
                          <code>{i.handle}</code>
                        </a>
                      )}
                    </EntryLine>
                    {i.note && (
                      <div className="note">
                        <InlineText text={i.note} />
                      </div>
                    )}
                  </li>
                ))}
              </ul>
            </dd>
          </Fragment>
        ))}
    </dl>
  );
}

export function SearchPage({ q }: { q: string }) {
  const { reg, report } = useView();
  const hits = useMemo(() => search(reg, q), [reg, q]);
  return (
    <div className="page">
      <h1>
        {hits.length} {hits.length === 1 ? "match" : "matches"} for “{q}”
      </h1>
      <ul className="hits">
        {hits.map((h) => {
          const c = coverage(reg, h.entry, report);
          return (
            <li key={h.entry.id}>
              <EntryLine e={h.entry}>
                <AreaLink slug={h.entry.area} />
                <VerdictBadge verdict={c.verdict} />
              </EntryLine>
              {h.where !== "id" && h.where !== "title" && (
                <div className="snippet">
                  <span className="muted">{h.where}:</span> {h.snippet}
                </div>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
}

function NotFound({ what }: { what: string }) {
  return (
    <div className="page">
      <h1>Not in this charter</h1>
      <p>
        {what} is not here. <a href={href({ page: "overview" })}>Back to the walkthrough</a>, or search for it.
      </p>
    </div>
  );
}

// Sidebar is the walkthrough as a table of contents: the needs, then the
// areas numbered in the reading order, then the observations.
export function Sidebar({ current, page }: { current: string; page: string }) {
  const { reg, report, tour } = useView();
  const number = useAreaNumber();
  const tallies = useMemo(() => new Map(reg.areas.map((a) => [a.slug, tally(reg, reg.inArea.get(a.slug)!, report)])), [reg, report]);
  const observations = reg.areas.filter((a) => a.kind === "observations");
  return (
    <nav className="sidebar">
      <a href={href({ page: "overview" })} className={`side-item top ${page === "overview" ? "on" : ""}`}>
        <span className="side-title">Walkthrough</span>
      </a>
      <a href={href({ page: "needs" })} className={`side-item top ${page === "needs" ? "on" : ""}`}>
        <span className="side-title">What people need</span>
        <span className="side-n muted">{reg.entries.filter((e) => e.need && e.status === "active").length}</span>
      </a>
      <div className="side-head">Areas, in reading order</div>
      {tour.path.map((a) => (
        <a key={a.slug} href={areaHref(a.slug)} className={`side-item ${current === a.slug ? "on" : ""}`} title={a.title}>
          <span className="side-num">{number(a.slug)}</span>
          <span className="side-title">{a.title}</span>
          <SideMarks t={tallies.get(a.slug)!} />
        </a>
      ))}
      {observations.length > 0 && (
        <div>
          <div className="side-head">What it rests on</div>
          {observations.map((a) => (
            <a key={a.slug} href={areaHref(a.slug)} className={`side-item ${current === a.slug ? "on" : ""}`}>
              <span className="side-title">{a.title}</span>
            </a>
          ))}
        </div>
      )}
    </nav>
  );
}

function SideMarks({ t }: { t: Tally }) {
  const bad = t.marks.missing + t.marks.fail + (t.decisions - t.realized);
  if (bad > 0)
    return (
      <span className="side-n x" title={`${bad} gaps`}>
        {bad}
      </span>
    );
  if (t.marks.waived > 0)
    return (
      <span className="side-n w" title={`${t.marks.waived} waived`}>
        {t.marks.waived}w
      </span>
    );
  return null;
}
