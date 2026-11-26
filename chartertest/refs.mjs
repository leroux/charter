// The charter's TypeScript extractor. Arguments: the typescript package
// directory and the module specifier tests import. It loads the compiler
// once, writes a line "ready", then serves requests until stdin ends: each
// is a line of JSON, a root and a list of files relative to it, answered
// by a line of JSON, {"records": [...]} or {"error": "..."}. A record is
// one per file, in order: parse diagnostics; every import of the specifier with its
// form (named, namespace, default, side-effect or type), its local binding
// and, for a named import, the exported name it takes; and every property
// access on a binding of `claims` with its line, the title of the enclosing
// test, whether any test call encloses it, whether a test or suite marked
// skip, fixme, fail or todo encloses it, and, when it is the first argument
// of a call of the bare name `verifies`, the call's remaining arguments as
// scenario handles, each a string or null when it is not a string literal,
// and whether the call is a leading statement of its test's body.
// It parses with the compiler API, syntax only, and reports facts. The
// rules are in Go.
import { createRequire } from "node:module";
import { readFileSync } from "node:fs";
import path from "node:path";
import { createInterface } from "node:readline";

const [compilerDir, specifier] = process.argv.slice(2);
const require = createRequire(import.meta.url);
const ts = require(path.resolve(compilerDir));

// Callee names that are not tests: suites, steps and hooks.
const notTests = new Set(["describe", "step", "beforeEach", "afterEach", "beforeAll", "afterAll"]);

// Modifiers under which a test's body never runs as a test.
const disabledTests = new Set(["skip", "fixme", "fail", "todo"]);

function extract(root, files) {
  const out = [];
  for (const file of files) {
    const text = readFileSync(path.resolve(root, file), "utf8");
    const kind = file.endsWith(".tsx") ? ts.ScriptKind.TSX : ts.ScriptKind.TS;
    const sf = ts.createSourceFile(file, text, ts.ScriptTarget.Latest, true, kind);
    const lineOf = (pos) => sf.getLineAndCharacterOfPosition(pos).line + 1;

    if (!sf.parseDiagnostics) throw new Error("this TypeScript version exposes no parseDiagnostics; the extractor cannot report syntax errors");
    const diagnostics = sf.parseDiagnostics.map((d) => ({
      line: lineOf(d.start ?? 0),
      message: ts.flattenDiagnosticMessageText(d.messageText, " "),
    }));

    const imports = [];
    const aliases = new Set();
    for (const st of sf.statements) {
      if (!ts.isImportDeclaration(st) || !ts.isStringLiteral(st.moduleSpecifier) || st.moduleSpecifier.text !== specifier) continue;
      const line = lineOf(st.getStart(sf));
      const clause = st.importClause;
      if (!clause) {
        imports.push({ line, form: "side-effect", local: "", imported: "" });
        continue;
      }
      if (clause.isTypeOnly) {
        imports.push({ line, form: "type", local: "", imported: "" });
        continue;
      }
      if (clause.name) imports.push({ line, form: "default", local: clause.name.text, imported: "" });
      const b = clause.namedBindings;
      if (b && ts.isNamespaceImport(b)) imports.push({ line, form: "namespace", local: b.name.text, imported: "" });
      if (b && ts.isNamedImports(b)) {
        const values = b.elements.filter((e) => !e.isTypeOnly);
        if (values.length === 0 && !clause.name) imports.push({ line, form: "type", local: "", imported: "" });
        for (const e of values) {
          const imported = (e.propertyName ?? e.name).text;
          imports.push({ line, form: "named", local: e.name.text, imported });
          if (imported === "claims") aliases.add(e.name.text);
        }
      }
    }

    // Names the file declares at its top level: a verifies among them is
    // the file's own, not the suite's helper.
    const localNames = new Set();
    for (const st of sf.statements) {
      if (ts.isFunctionDeclaration(st) && st.name) localNames.add(st.name.text);
      if (ts.isVariableStatement(st)) {
        for (const d of st.declarationList.declarations) if (ts.isIdentifier(d.name)) localNames.add(d.name.text);
      }
    }

    const refs = [];
    const visit = (node) => {
      if (ts.isPropertyAccessExpression(node) && ts.isIdentifier(node.expression) && aliases.has(node.expression.text)) {
        const ref = { line: lineOf(node.getStart(sf)), name: node.name.text, test: enclosingTest(node, sf), inTest: insideTest(node), call: false, handles: [], disabled: insideDisabledTest(node), localCallee: false, leading: false };
        const call = verifiesCall(node);
        if (call) {
          ref.call = true;
          ref.handles = call.arguments.slice(1).map((a) => (ts.isStringLiteralLike(a) ? a.text : null));
          ref.localCallee = localNames.has(call.expression.text);
          ref.leading = isLeading(call);
        }
        refs.push(ref);
      }
      ts.forEachChild(node, visit);
    };
    visit(sf);
    out.push({ file, diagnostics, imports, refs });
  }
  return out;
}

process.stdout.write("ready\n");
for await (const line of createInterface({ input: process.stdin })) {
  let reply;
  try {
    const { root, files } = JSON.parse(line);
    reply = { records: extract(root, files) };
  } catch (e) {
    reply = { error: String(e?.stack ?? e) };
  }
  process.stdout.write(JSON.stringify(reply) + "\n");
}

// isTestCallee reports whether a callee is `test` or `it`, possibly followed
// by property accesses and calls such as `.only`, `.fail` or `.each(...)`,
// none of them `describe`, `step` or a hook.
function isTestCallee(expr) {
  let e = expr;
  for (;;) {
    if (ts.isIdentifier(e)) return e.text === "test" || e.text === "it";
    if (ts.isPropertyAccessExpression(e)) {
      if (notTests.has(e.name.text)) return false;
      e = e.expression;
    } else if (ts.isCallExpression(e)) {
      e = e.expression;
    } else {
      return false;
    }
  }
}

// enclosingTest is the title of the innermost test call around node: the
// first argument's text when it is a string or a template literal, and ""
// when it is anything else or there is no such call.
function enclosingTest(node, sf) {
  for (let n = node.parent; n; n = n.parent) {
    if (!ts.isCallExpression(n) || !isTestCallee(n.expression)) continue;
    const arg = n.arguments[0];
    if (!arg) return "";
    if (ts.isStringLiteralLike(arg)) return arg.text;
    if (ts.isTemplateExpression(arg)) return arg.getText(sf).slice(1, -1);
    return "";
  }
  return "";
}

// insideTest reports whether any test call encloses node, whatever its title.
function insideTest(node) {
  for (let n = node.parent; n; n = n.parent) {
    if (ts.isCallExpression(n) && isTestCallee(n.expression)) return true;
  }
  return false;
}

// insideDisabledTest reports whether a test or suite call that encloses
// node, `test`, `it` or `describe` followed by property accesses and calls,
// carries skip, fixme, fail or todo anywhere in its callee chain: under
// describe.skip or test.describe.fixme no test body runs either.
function insideDisabledTest(node) {
  for (let n = node.parent; n; n = n.parent) {
    if (!ts.isCallExpression(n)) continue;
    let disabled = false;
    let e = n.expression;
    for (;;) {
      if (ts.isPropertyAccessExpression(e)) {
        if (disabledTests.has(e.name.text)) disabled = true;
        e = e.expression;
      } else if (ts.isCallExpression(e)) {
        e = e.expression;
      } else {
        break;
      }
    }
    if (disabled && ts.isIdentifier(e) && (e.text === "test" || e.text === "it" || e.text === "describe")) return true;
  }
  return false;
}

// verifiesCall returns the call `verifies(node, ...)` whose first argument
// is node and whose callee is the bare name verifies, or undefined. A
// verifies reached through a property is not the helper the suite imports,
// so node there is a reference like any other.
function verifiesCall(node) {
  const p = node.parent;
  if (!p || !ts.isCallExpression(p) || p.arguments[0] !== node) return undefined;
  if (isBareVerifies(p)) return p;
  return undefined;
}

// isBareVerifies reports whether call's callee is the identifier verifies.
function isBareVerifies(call) {
  return ts.isIdentifier(call.expression) && call.expression.text === "verifies";
}

// isLeading reports whether call opens the body of a test: an expression
// statement of the block a test call's function argument runs, with only
// such verifies statements before it, as Go's markLeading reads a test.
// One that a statement precedes, or a block, a branch or a closure
// encloses, is not leading.
function isLeading(call) {
  const st = call.parent;
  if (!st || !ts.isExpressionStatement(st) || st.expression !== call) return false;
  const block = st.parent;
  if (!block || !ts.isBlock(block)) return false;
  const fn = block.parent;
  if (!fn || !(ts.isArrowFunction(fn) || ts.isFunctionExpression(fn)) || fn.body !== block) return false;
  const test = fn.parent;
  if (!test || !ts.isCallExpression(test) || !test.arguments.includes(fn) || !isTestCallee(test.expression)) return false;
  for (const prev of block.statements) {
    if (prev === st) return true;
    if (!ts.isExpressionStatement(prev) || !ts.isCallExpression(prev.expression) || !isBareVerifies(prev.expression)) return false;
  }
  return false;
}
