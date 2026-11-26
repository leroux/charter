import { useCallback, useEffect, useMemo, useState } from "react";
import { ViewContext } from "./components/bits";
import { AreaPage, CoveragePage, EntryPage, GapsPage, NeedsPage, SearchPage, Sidebar, Walkthrough } from "./components/pages";
import { dataURL, loadFromText, loadFromURL, loadIndex, loadReport, reportURL, sameData, type CharterLink } from "./load";
import type { Registry } from "./model";
import type { Report } from "./report";
import { href, parseRoute, type Route } from "./route";
import { arrange } from "./tour";

export function App() {
  const [reg, setReg] = useState<Registry | null>(null);
  const [report, setReport] = useState<Report | null>(null);
  const [charters, setCharters] = useState<CharterLink[]>([]);
  const [route, setRoute] = useState<Route>(() => parseRoute(window.location.hash));
  const [error, setError] = useState("");
  const tour = useMemo(() => (reg ? arrange(reg) : null), [reg]);

  useEffect(() => {
    // A hash naming an element on the page, #SW-118 from a table of
    // contents, is a scroll, not a route.
    const onHash = () => {
      if (window.location.hash.startsWith("#/") || window.location.hash === "") setRoute(parseRoute(window.location.hash));
    };
    window.addEventListener("hashchange", onHash);
    return () => window.removeEventListener("hashchange", onHash);
  }, []);
  // Scrolling returns a promise in newer browsers; an effect returns only a
  // clean-up, so the call is a statement.
  useEffect(() => {
    window.scrollTo(0, 0);
  }, [route]);

  const take = useCallback((text: string, name: string) => {
    const got = loadFromText(text, name);
    if (got.kind === "registry") setReg(got.reg);
    else setReport(got.report);
    setError("");
  }, []);

  useEffect(() => {
    const search = window.location.search;
    const fail = (err: Error) => setError(err.message);
    loadFromURL(dataURL(search)).then(setReg).catch(fail);
    loadReport(reportURL(search)).then(setReport).catch(fail);
    loadIndex().then(setCharters).catch(fail);
  }, []);

  useEffect(() => {
    const onDrop = (ev: DragEvent) => {
      ev.preventDefault();
      const file = ev.dataTransfer?.files?.[0];
      if (!file) return;
      file
        .text()
        .then((text) => take(text, file.name))
        .catch((err: Error) => setError(err.message));
    };
    const onDrag = (ev: DragEvent) => ev.preventDefault();
    window.addEventListener("drop", onDrop);
    window.addEventListener("dragover", onDrag);
    return () => {
      window.removeEventListener("drop", onDrop);
      window.removeEventListener("dragover", onDrag);
    };
  }, [take]);

  return (
    <div className="app">
      <TopBar reg={reg} report={report} charters={charters} route={route} onClearReport={() => setReport(null)} />
      {error && (
        <div className="error" role="alert">
          {error}
          {reg ? " — showing the last charter that loaded" : ""}
        </div>
      )}
      {!reg || !tour ? (
        <div className="empty">
          {error ? "" : "loading… "}
          Drop a charter's JSON here, as claimsmap writes it, or open the page with ?data=&lt;url&gt;.
        </div>
      ) : (
        <ViewContext.Provider value={{ reg, report, tour }}>
          <div className="body">
            <Sidebar
              page={route.page}
              current={route.page === "area" ? route.slug : route.page === "entry" ? (reg.byId.get(route.id)?.area ?? "") : ""}
            />
            <main className="main">
              {route.page === "overview" && <Walkthrough />}
              {route.page === "needs" && <NeedsPage />}
              {route.page === "coverage" && <CoveragePage />}
              {route.page === "area" && <AreaPage key={route.slug} slug={route.slug} />}
              {route.page === "entry" && <EntryPage key={route.id} id={route.id} />}
              {route.page === "gaps" && <GapsPage />}
              {route.page === "search" && <SearchPage q={route.q} />}
            </main>
          </div>
        </ViewContext.Provider>
      )}
    </div>
  );
}

interface TopBarProps {
  reg: Registry | null;
  report: Report | null;
  charters: CharterLink[];
  route: Route;
  onClearReport: () => void;
}

function TopBar({ reg, report, charters, route, onClearReport }: TopBarProps) {
  const [q, setQ] = useState(route.page === "search" ? route.q : "");
  useEffect(() => {
    if (route.page !== "search") setQ("");
  }, [route]);
  const go = (value: string) => {
    setQ(value);
    // An exact id goes straight to the entry.
    const exact = reg?.byId.get(value.trim().toUpperCase()) ?? reg?.byId.get(value.trim());
    window.location.hash = exact ? href({ page: "entry", id: exact.id }) : value.trim() ? href({ page: "search", q: value }) : href({ page: "overview" });
  };
  const current = dataURL(window.location.search);
  const switchTo = (data: string) => {
    const p = new URLSearchParams(window.location.search);
    p.delete("report");
    p.set("data", data);
    window.location.href = `${window.location.pathname}?${p}#/`;
  };
  return (
    <header className="top">
      <a className="brand" href={href({ page: "overview" })}>
        Charter
      </a>
      {charters.length > 1 && (
        <span className="seg">
          {charters.map((c) => (
            <button key={c.data} className={sameData(c.data, current) ? "on" : ""} onClick={() => switchTo(c.data)} title={c.data}>
              {c.label}
            </button>
          ))}
        </span>
      )}
      <nav className="topnav">
        <a className={route.page === "overview" ? "on" : ""} href={href({ page: "overview" })}>
          Walkthrough
        </a>
        <a className={route.page === "needs" ? "on" : ""} href={href({ page: "needs" })}>
          Needs
        </a>
        <a className={route.page === "coverage" ? "on" : ""} href={href({ page: "coverage" })}>
          Coverage
        </a>
        <a className={route.page === "gaps" ? "on" : ""} href={href({ page: "gaps" })}>
          Gaps
        </a>
      </nav>
      <form
        className="searchform"
        onSubmit={(ev) => {
          ev.preventDefault();
          go(q);
        }}
      >
        <input type="search" placeholder="search ids, titles, scenarios, code, tests…" value={q} onChange={(ev) => setQ(ev.target.value)} />
      </form>
      <span className="spacer" />
      {report ? (
        <span className="small">
          report loaded{" "}
          <button className="link" onClick={onClearReport}>
            clear
          </button>
        </span>
      ) : (
        <span className="small muted" title="drop a claimsgen report JSON on the page, or write one beside the charter's">
          no report
        </span>
      )}
      <span className="source small muted">{reg?.source}</span>
    </header>
  );
}
