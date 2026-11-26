// An entry's prose as the documents write it: hard-wrapped paragraphs,
// "- " lists, `code` and entry ids. Enough markdown to read a decision, not
// a renderer. Pure.

export type Inline = { t: "text"; v: string } | { t: "code"; v: string } | { t: "id"; v: string };
export type Block = { t: "p"; v: Inline[] } | { t: "list"; v: Inline[][] };

// An id: a numbered one, SW-118, or a dated one, SW-261001-K7Q2.
const ID = /\b[A-Z]{2,4}-(?:\d{6}-[A-Z0-9]{4}|\d+)\b/g;

export function inline(text: string, isId: (id: string) => boolean): Inline[] {
  const out: Inline[] = [];
  const push = (x: Inline) => {
    const last = out[out.length - 1];
    if (x.t === "text" && last?.t === "text") last.v += x.v;
    else if (x.v) out.push(x);
  };
  text.split(/(`[^`]*`)/).forEach((part) => {
    if (part.length > 1 && part.startsWith("`") && part.endsWith("`")) {
      push({ t: "code", v: part.slice(1, -1) });
      return;
    }
    let at = 0;
    for (const m of part.matchAll(ID)) {
      if (!isId(m[0])) continue;
      push({ t: "text", v: part.slice(at, m.index) });
      push({ t: "id", v: m[0] });
      at = m.index! + m[0].length;
    }
    push({ t: "text", v: part.slice(at) });
  });
  return out;
}

export function blocks(text: string, isId: (id: string) => boolean): Block[] {
  const out: Block[] = [];
  for (const para of text.split(/\n\s*\n/)) {
    const lines = para.split("\n");
    if (lines.length > 0 && /^\s*[-*] /.test(lines[0])) {
      const items: string[] = [];
      for (const l of lines) {
        if (/^\s*[-*] /.test(l)) items.push(l.replace(/^\s*[-*] /, ""));
        else if (items.length > 0) items[items.length - 1] += " " + l.trim();
      }
      out.push({ t: "list", v: items.map((i) => inline(i, isId)) });
    } else {
      const flat = lines.map((l) => l.trim()).join(" ").trim();
      if (flat) out.push({ t: "p", v: inline(flat, isId) });
    }
  }
  return out;
}
