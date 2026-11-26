// The page as a URL hash, so a link can be shared and back works:
//   #/              the walkthrough
//   #/needs         the needs, by chapter
//   #/coverage      the totals and the coverage of every need and area
//   #/area/turn     an area document
//   #/e/SW-118      an entry, with the tree of what serves it
//   #/gaps          the work list
//   #/search/words  search results
// Pure.

export type Route =
  | { page: "overview" }
  | { page: "needs" }
  | { page: "coverage" }
  | { page: "area"; slug: string }
  | { page: "entry"; id: string }
  | { page: "gaps" }
  | { page: "search"; q: string };

export function parseRoute(hash: string): Route {
  const [page, ...rest] = hash.replace(/^#?\/?/, "").split("/");
  const arg = decodeURIComponent(rest.join("/"));
  switch (page) {
    case "area":
      return arg ? { page: "area", slug: arg } : { page: "overview" };
    case "e":
      return arg ? { page: "entry", id: arg } : { page: "overview" };
    case "gaps":
      return { page: "gaps" };
    case "needs":
      return { page: "needs" };
    case "coverage":
      return { page: "coverage" };
    case "search":
      return { page: "search", q: arg };
    default:
      return { page: "overview" };
  }
}

export function href(r: Route): string {
  switch (r.page) {
    case "overview":
      return "#/";
    case "area":
      return `#/area/${encodeURIComponent(r.slug)}`;
    case "entry":
      return `#/e/${encodeURIComponent(r.id)}`;
    case "gaps":
      return "#/gaps";
    case "needs":
      return "#/needs";
    case "coverage":
      return "#/coverage";
    case "search":
      return `#/search/${encodeURIComponent(r.q)}`;
  }
}

export const entryHref = (id: string) => href({ page: "entry", id });
export const areaHref = (slug: string) => href({ page: "area", slug });
