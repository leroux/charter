# claims-viewer

A charter as documents: its needs, its area documents in its reading order,
its observations, each entry with what the gate knows of it. It reads the
JSON `go run ./cmd/claimsmap` writes from `chartertest.Registry`, and
names no project's areas, ids or charters: any charter in the form
`claims/README.md` describes reads the same way.

    make claims-viewer-dev      # export both charters, then the Vite dev server
    make claims-viewer          # export, typecheck, test, vite build -> claims-viewer/dist
    make claims-viewer-report   # run the Go tests, so the page shows what passed
    cd claims-viewer && pnpm test

The export runs the gate's reader, which needs node and
`web/node_modules`; the targets install them.

Pages, each a URL hash a link can share:

- **Walkthrough** (`#/`): how to read an entry, then the areas numbered
  in the charter's reading order (`reading-order.txt`, which the gate
  holds complete), each with the sentence that opens its document and the
  needs it does most to deliver.
- **Needs** (`#/needs`): every need under the area whose decisions most
  often reach it, the earlier in the reading order on a tie, as the
  digest places them.
- **Area** (`#/area/refill`): the document as written, its summary, the
  needs its decisions deliver, then entry by entry with coverage beside
  each part: the code that realizes it, and per scenario the tests that
  cite it or the waiver that accepts none. Previous and next follow the
  reading order, so the charter reads end to end. Deprecated entries are
  hidden until asked for.
- **Entry** (`#/e/SW-118`): one entry, how it reaches its needs, and the
  tree of decisions that serve it, by area in the reading order.
- **Coverage** (`#/coverage`): the totals, and every area and need with
  its coverage.
- **Gaps** (`#/gaps`): what the gate fails on (an unrealized decision, an
  untested scenario, a failing test), then what it lets through: waivers,
  decisions that reach no need, relations to retired entries, deprecated
  entries still cited, planned entries, observations in doubt.
- **Search** (`#/search/words`): ids, titles, text, scenarios, surfaces,
  code and tests; an exact id opens the entry.

Coverage is rule 25's, read from the export: a scenario is tested when a
test cites it, waived when `waivers.txt` accepts it, untested otherwise.
That a test cites a scenario is not that it passed: with
`claims-report.json` beside the page (`make claims-viewer-report`), or a
`claimsgen report` JSON dropped on it, the Go tests' outcomes show; the
TypeScript tests are not in a report and stay "tested".

The page loads `?data=<file>`, `claims.json` by default, and the report
beside it, `claims-report.json` for `claims.json`. `charters.json`, when
present, lists the charters the top bar switches between, as
`[{"data": "claims.json", "label": "product"}]`; this repository's
Makefile writes it for the product's charter and the charter's own. Drop
a registry or a report JSON on the page to look at another. A
`file:line` opens in VS Code when the export knows the checkout's root.
