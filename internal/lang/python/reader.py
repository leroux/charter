"""The charter's Python reader. It reads a JSON request on stdin, a root and
the files to read with what it needs to know of the tree, and writes the
facts it finds as one JSON object on stdout: anchors, citations,
references, environment constants, reads and literals, and refused
directives and calls. It parses with the standard ast and tokenize modules
alone and decides nothing: the rules are the gate's."""

import ast
import io
import json
import os
import re
import sys
import tokenize

DIRECTIVE = re.compile(r"^#\s*charter:(realizes|concerns)\b(.*)$")
DYNAMIC = re.compile(r"#\s*charter:dynamic\b")
ENV_NAME = re.compile(r"^[A-Z][A-Z0-9]*(?:_[A-Z0-9]+)+$")
HELPERS = {"verifies", "observes", "Spec", "Observation"}


def is_test_file(rel):
    base = os.path.basename(rel)
    return base.startswith("test_") or base.endswith("_test.py") or base == "conftest.py"


def comments(src):
    """Each comment's line and text."""
    out = {}
    try:
        for tok in tokenize.generate_tokens(io.StringIO(src).readline):
            if tok.type == tokenize.COMMENT:
                out[tok.start[0]] = tok.string
    except (tokenize.TokenError, IndentationError, SyntaxError):
        pass
    return out


def split_ids(rest):
    return [p.strip() for p in rest.split(",") if p.strip()]


class File:
    def __init__(self, req, rel):
        self.req = req
        self.rel = rel
        with open(os.path.join(req["root"], rel), encoding="utf-8") as f:
            self.src = f.read()
        self.tree = ast.parse(self.src, filename=rel)
        self.lines = self.src.split("\n")
        self.comments = comments(self.src)
        self.facts = req["facts"]

    def at(self, line):
        return {"File": self.rel, "Line": line}

    def refuse(self, line, rule, message):
        self.facts["Refused"].append(dict(self.at(line), Rule=rule, Message=message))

    # Anchors ---------------------------------------------------------

    def declarations(self, node=None, prefix=""):
        """Each def and class with the line its decorators start at."""
        node = node or self.tree
        for child in ast.iter_child_nodes(node):
            if isinstance(child, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
                start = min([d.lineno for d in child.decorator_list] + [child.lineno])
                kind = "class" if isinstance(child, ast.ClassDef) else "def"
                name = prefix + child.name
                yield start, kind, name
                yield from self.declarations(child, name + ".")
            else:
                yield from self.declarations(child, prefix)

    def decls(self, node=None, prefix=""):
        """Each def and class, nested ones by their dotted path, from its
        first decorator to its last line, and each top-level assignment."""
        package = os.path.splitext(os.path.basename(self.rel))[0]
        node = node or self.tree
        for child in ast.iter_child_nodes(node):
            if isinstance(child, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
                start = min([d.lineno for d in child.decorator_list] + [child.lineno])
                name = prefix + child.name
                self.facts["Decls"].append(dict(self.at(start), End=child.end_lineno, Name=name, Package=package))
                self.decls(child, name + ".")
            elif node is self.tree and isinstance(child, (ast.Assign, ast.AnnAssign)):
                targets = child.targets if isinstance(child, ast.Assign) else [child.target]
                for t in targets:
                    if isinstance(t, ast.Name):
                        self.facts["Decls"].append(dict(self.at(child.lineno), End=child.end_lineno, Name=t.id, Package=package))

    def anchors(self):
        starts = {start: (kind, name) for start, kind, name in self.declarations()}
        for line, text in sorted(self.comments.items()):
            m = DIRECTIVE.match(text.strip())
            if not m:
                continue
            kind, rest = m.group(1), m.group(2)
            # The directive is on the comment lines directly above a def or
            # a class, or its decorators: every line between is a comment.
            n = line + 1
            while n <= len(self.lines) and self.lines[n - 1].strip().startswith("#"):
                n += 1
            owner = starts.get(n)
            code_before = self.lines[line - 1].strip() and not self.lines[line - 1].strip().startswith("#")
            if owner is None or code_before:
                self.refuse(line, "rule 24", "# charter:%s is written on the comment lines directly above the def or class that carries the decision; this comment is above none" % kind)
                continue
            ids = split_ids(rest)
            if not ids:
                self.refuse(line, "rule 24", "# charter:%s names at least one spec" % kind)
                continue
            self.facts["Anchors"].append(dict(self.at(line), Kind=kind, IDs=ids, Decl=owner[0] + " " + owner[1], DeclKind=owner[0]))

    # The environment -------------------------------------------------

    def env_consts(self):
        base = os.path.splitext(os.path.basename(self.rel))[0]
        for st in self.tree.body:
            if isinstance(st, ast.Assign) and len(st.targets) == 1 and isinstance(st.targets[0], ast.Name):
                name = st.targets[0].id
                if not name.startswith("_") and isinstance(st.value, ast.Constant) and isinstance(st.value.value, str):
                    self.facts["EnvConsts"].append(dict(self.at(st.lineno), Name=base + "." + name, Value=st.value.value))

    def env_aliases(self):
        """The names this file binds to the constants module, and the
        constants it imports by name."""
        module = self.req.get("env_module", "")
        base = module.rsplit(".", 1)[-1] if module else ""
        modules, consts = set(), {}
        if not module:
            return modules, consts
        for st in ast.walk(self.tree):
            if isinstance(st, ast.Import):
                for a in st.names:
                    if a.name == module:
                        modules.add(a.asname or a.name)
            elif isinstance(st, ast.ImportFrom):
                src = (st.module or "")
                for a in st.names:
                    if (src + "." + a.name).strip(".") == module or (st.level and a.name == base):
                        modules.add(a.asname or a.name)
                    elif src == module or (st.level and src == base):
                        consts[a.asname or a.name] = base + "." + a.name
        return modules, consts

    def dotted(self, node):
        if isinstance(node, ast.Name):
            return node.id
        if isinstance(node, ast.Attribute):
            inner = self.dotted(node.value)
            return inner + "." + node.attr if inner else ""
        return ""

    def env_reads(self):
        modules, consts = self.env_aliases()
        base = self.req.get("env_module", "").rsplit(".", 1)[-1]
        dynamic = {line for line, text in self.comments.items() if DYNAMIC.search(text)}

        def read(node, arg):
            r = dict(self.at(node.lineno), Const="", Literal="", Dynamic=node.lineno in dynamic)
            if isinstance(arg, ast.Attribute) and self.dotted(arg.value) in modules:
                r["Const"] = base + "." + arg.attr
            elif isinstance(arg, ast.Name) and arg.id in consts:
                r["Const"] = consts[arg.id]
            elif isinstance(arg, ast.Constant) and isinstance(arg.value, str):
                r["Literal"] = arg.value
            self.facts["EnvReads"].append(r)

        for node in ast.walk(self.tree):
            if isinstance(node, ast.Call) and node.args:
                f = self.dotted(node.func)
                if f in ("os.getenv", "os.environ.get", "getenv", "environ.get"):
                    read(node, node.args[0])
            elif isinstance(node, ast.Subscript) and self.dotted(node.value) in ("os.environ", "environ"):
                read(node, node.slice)
            if isinstance(node, ast.Constant) and isinstance(node.value, str) and ENV_NAME.match(node.value):
                self.facts["Literals"].append(dict(self.at(node.lineno), Value=node.value))

    # Citations -------------------------------------------------------

    def claims_names(self):
        """The names bound to the claims module, and the helpers imported
        from it by name."""
        mod = self.req["claims_module"]
        aliases, helpers = set(), {}
        for st in ast.walk(self.tree):
            if isinstance(st, ast.Import):
                for a in st.names:
                    if a.name == mod or a.name.endswith("." + mod):
                        aliases.add(a.asname or a.name)
            elif isinstance(st, ast.ImportFrom):
                src = st.module or ""
                if src == mod or src.endswith("." + mod):
                    for a in st.names:
                        if a.name == "*":
                            self.refuse(st.lineno, "rule 10", "the claims module must not be star-imported")
                        elif a.name in ("verifies", "observes"):
                            helpers[a.asname or a.name] = a.name
                else:
                    for a in st.names:
                        if a.name == mod:
                            aliases.add(a.asname or a.name)
        return aliases, helpers

    def citations(self):
        aliases, helpers = self.claims_names()
        if not aliases and not helpers:
            return

        def marker(call):
            if not isinstance(call, ast.Call):
                return ""
            f = call.func
            if isinstance(f, ast.Name) and f.id in helpers:
                return helpers[f.id]
            if isinstance(f, ast.Attribute) and f.attr in ("verifies", "observes") and self.dotted(f.value) in aliases:
                return f.attr
            return ""

        leading, collected = {}, {}

        def mark(fn, style):
            body = fn.body
            i = 0
            if body and isinstance(body[0], ast.Expr) and isinstance(body[0].value, ast.Constant) and isinstance(body[0].value.value, str):
                i = 1  # a docstring runs nothing
            while i < len(body) and isinstance(body[i], ast.Expr) and marker(body[i].value):
                leading[id(body[i].value)] = (fn.name, style)
                i += 1
            collected[id(fn)] = (fn.name, style)

        def is_testcase(cls):
            return any(self.dotted(b).endswith("TestCase") for b in cls.bases)

        for st in self.tree.body:
            if isinstance(st, (ast.FunctionDef, ast.AsyncFunctionDef)) and st.name.startswith("test"):
                mark(st, "pytest")
            elif isinstance(st, ast.ClassDef):
                for m in st.body:
                    if isinstance(m, (ast.FunctionDef, ast.AsyncFunctionDef)) and m.name.startswith("test"):
                        if is_testcase(st):
                            mark(m, "unittest")
                        elif st.name.startswith("Test"):
                            mark(m, "pytest-class")

        enclosing = {}

        def walk(node, fn):
            for child in ast.iter_child_nodes(node):
                if isinstance(child, (ast.FunctionDef, ast.AsyncFunctionDef, ast.Lambda)):
                    walk(child, child)
                else:
                    enclosing[id(child)] = fn
                    walk(child, fn)

        walk(self.tree, None)
        for node in ast.walk(self.tree):
            if isinstance(node, (ast.Assign, ast.AugAssign, ast.AnnAssign)):
                targets = node.targets if isinstance(node, ast.Assign) else [node.target]
                for t in targets:
                    base = t
                    while isinstance(base, (ast.Attribute, ast.Subscript)) and not (isinstance(base, ast.Attribute) and self.dotted(base.value) in aliases):
                        base = base.value
                    if isinstance(base, ast.Attribute) and self.dotted(base.value) in aliases and base.attr not in HELPERS:
                        self.refuse(node.lineno, "rule 1", "%s is assigned outside its declaration" % base.attr)
            name = marker(node)
            if name:
                if len(node.args) > 1:
                    node.args[1]._charter_cited = True
                fn = enclosing.get(id(node))
                call = "verifies" if name == "verifies" else "observes"
                if fn is None or id(fn) not in collected:
                    if fn is None or isinstance(fn, ast.Lambda) or not fn.name.startswith("test"):
                        self.refuse(node.lineno, "rule 24", "%s sits outside a test; it is the first statement of the test that verifies the spec" % call)
                    else:
                        self.refuse(node.lineno, "rule 24", "%s sits in %s, which the runner never collects; a test is a test_ method of a unittest.TestCase or a test_ function" % (call, fn.name))
                    continue
                if id(node) not in leading:
                    self.refuse(node.lineno, "rule 24", "%s is the leading statement of its test, before anything runs, skips or branches; here something precedes or encloses it" % call)
                    continue
                test, style = leading[id(node)]
                first = node.args[0] if node.args else None
                if style == "unittest" and not (isinstance(first, ast.Name) and first.id == "self"):
                    self.refuse(node.lineno, "rule 24", "a unittest test passes self as the first argument of %s" % call)
                    continue
                if style != "unittest" and not (isinstance(first, ast.Constant) and first.value is None):
                    self.refuse(node.lineno, "rule 24", "a pytest test passes None as the first argument of %s" % call)
                    continue
                spec = node.args[1] if len(node.args) > 1 else None
                if not (isinstance(spec, ast.Attribute) and self.dotted(spec.value) in aliases and spec.attr not in HELPERS):
                    self.refuse(node.lineno, "rule 24", "%s names a spec by its symbol, such as claims.SW_118" % call)
                    continue
                handles = []
                for a in node.args[2:]:
                    if isinstance(a, ast.Constant) and isinstance(a.value, str):
                        handles.append(a.value)
                    else:
                        self.refuse(node.lineno, "rule 24", "a scenario handle is a string literal")
                self.facts["Citations"].append(dict(self.at(node.lineno), Symbol=spec.attr, Handles=handles, Test=test, Observes=name == "observes"))
                spec._charter_cited = True
        for node in ast.walk(self.tree):
            if isinstance(node, ast.Attribute) and self.dotted(node.value) in aliases and node.attr not in HELPERS and not getattr(node, "_charter_cited", False):
                if isinstance(node.ctx, ast.Load):
                    self.facts["References"].append(dict(self.at(node.lineno), Symbol=node.attr))


def main():
    req = json.load(sys.stdin)
    facts = {k: [] for k in ("Anchors", "Citations", "References", "EnvConsts", "EnvReads", "Literals", "Refused", "Decls", "TestFiles")}
    req["facts"] = facts
    env_file = req.get("env_file", "")
    for rel in req["files"]:
        try:
            f = File(req, rel)
        except SyntaxError as e:
            print(json.dumps({"error": "parse %s:%s: %s" % (rel, e.lineno, e.msg)}))
            return
        if is_test_file(rel):
            facts["TestFiles"].append(rel)
            f.citations()
            continue
        f.anchors()
        f.decls()
        if rel == env_file:
            f.env_consts()
            continue
        f.env_reads()
    for k, v in facts.items():
        if k != "TestFiles":
            v.sort(key=lambda f: (f["File"], f["Line"]))
    print(json.dumps({"facts": facts}))


main()
