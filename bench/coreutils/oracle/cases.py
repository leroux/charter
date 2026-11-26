"""The generated corpus: every case the oracle runs, GNU against the clone.

A case is a dict: id, group ("spec" or a task number), tool, args, stdin
(bytes), files ({name: bytes}, written to the case's directory, the
working directory of both runs) and stderr (whether stderr must match too).
The corpus comes from one fixed seed, so every run sees the same cases.
"""

import random

SEED = 20261126

WORDS = ["apple", "Apple", "APPLE", "banana", "Banana", "cherry", "date", "fig",
         "kiwi", "Zebra", "zebra", "a", "b", "B", "abc", "ab", "_x", "x-y", "10", "9"]
NUMS = ["0", "1", "2", "9", "10", "11", "100", "-1", "-10", "-0", "007", "2.5", "2.50",
        ".5", "-.5", "1e3", "+5", "  42", "\t7", "3x", "x3", "", "1.2.3", "99999999999999999999",
        "-99999999999999999999", "0.0001", "abc", "1,000"]
HUMAN = ["1K", "2K", "512", "1M", "1.5K", "3G", "0", "10", "2M", "999", "1T", "1k", "-1K", "100M"]


def lines_words(r, n):
    out = []
    for _ in range(n):
        k = r.choice([1, 1, 1, 2, 3])
        out.append(" ".join(r.choice(WORDS) for _ in range(k)))
    return out


def lines_dups(r, n, fold=False):
    """Runs of repeated lines, some differing only in case."""
    out = []
    while len(out) < n:
        w = r.choice(WORDS[:12])
        for _ in range(r.choice([1, 1, 2, 3])):
            out.append(w.upper() if fold and r.random() < 0.3 else w)
    return out


def lines_nums(r, n):
    out = []
    for _ in range(n):
        v = r.choice(NUMS)
        if r.random() < 0.3:
            v += " " + r.choice(WORDS)
        out.append(v)
    return out


def lines_records(r, n, d, nf=None):
    out = []
    for _ in range(n):
        if r.random() < 0.12:
            out.append(r.choice(["no delimiter here", "plain", "", "zz"]))
            continue
        k = nf or r.choice([2, 3, 3, 4, 5])
        fs = []
        for _ in range(k):
            fs.append(r.choice(WORDS + NUMS[:12] + [""]))
        out.append(d.join(fs))
    return out


def lines_text(r, n):
    """Prose-ish lines with tabs, runs of blanks, control and high bytes."""
    pieces = WORDS + ["  ", "\t", " \t ", "\x01", "\x7f", "\xff", "\v", "\r", "\f", "é", "x\x00y"]
    out = []
    for _ in range(n):
        out.append("".join(r.choice(pieces) + r.choice(["", " ", " "]) for _ in range(r.randint(0, 6))))
    return out


def blob(lines, r, trailing=None):
    if not lines:
        return b""
    s = "\n".join(lines)
    if trailing is None:
        trailing = r.random() < 0.8
    if trailing:
        s += "\n"
    return s.encode("latin-1")


def case(cases, group, tool, args, stdin=b"", files=None, stderr=False):
    cases.append({"id": f"{group}-{tool}-{len(cases):03d}", "group": group, "tool": tool,
                  "args": list(args), "stdin": stdin, "files": files or {}, "stderr": stderr})


def spec_cases(r, cases):
    """The seed's base behaviour; must hold after every task."""
    # wc: no flags, -l -w -c and their combinations, stdin or one file.
    flagsets = [[], ["-l"], ["-w"], ["-c"], ["-lw"], ["-wc"], ["-lc"], ["-lwc"], ["-c", "-l"],
                ["-w", "-l"], ["--lines"], ["--words"], ["--bytes"]]
    gens = [lambda: lines_words(r, r.randint(0, 12)), lambda: lines_text(r, r.randint(1, 10)),
            lambda: lines_nums(r, r.randint(1, 8)), lambda: [""] * r.randint(1, 4)]
    for i in range(70):
        data = blob(r.choice(gens)(), r)
        if i % 9 == 0:
            data = r.choice([b"", b"x", b"\n", b"  \n\t", b"a b", b"\x01\x02\n", b"one\x00two three"])
        fl = r.choice(flagsets)
        if r.random() < 0.5:
            case(cases, "spec", "wc", fl, stdin=data)
        else:
            case(cases, "spec", "wc", fl + ["in.txt"], files={"in.txt": data})
    # cut -f with -d, lines with and without the delimiter, stdin or files.
    lists = ["1", "2", "3", "1,3", "3,1", "2-", "-2", "1-2,4", "2-3", "5", "1-", "2,2", "1,2-3", "4-,1"]
    for i in range(80):
        d = r.choice([":", ",", "\t", " ", ";"])
        data = blob(lines_records(r, r.randint(0, 10), d), r)
        lst = r.choice(lists)
        if d == "\t" and r.random() < 0.5:
            dargs = []
        else:
            dargs = r.choice([["-d" + d], ["-d", d], ["--delimiter=" + d]])
        fargs = r.choice([["-f" + lst], ["-f", lst], ["--fields=" + lst]])
        args = dargs + fargs if r.random() < 0.6 else fargs + dargs
        k = r.random()
        if k < 0.5:
            case(cases, "spec", "cut", args, stdin=data)
        elif k < 0.8:
            case(cases, "spec", "cut", args + ["in.txt"], files={"in.txt": data})
        else:
            data2 = blob(lines_records(r, r.randint(1, 5), d), r)
            case(cases, "spec", "cut", args + ["a.txt", "b.txt"], files={"a.txt": data, "b.txt": data2})
    # sort plain, -n, -r, stdin or files.
    sflags = [[], ["-r"], ["-n"], ["-rn"], ["-nr"], ["-n", "-r"], ["--numeric-sort"], ["--reverse"]]
    for i in range(70):
        g = r.choice([lines_words, lines_nums, lines_nums, lines_text, lines_dups])
        data = blob(g(r, r.randint(0, 14)), r)
        fl = r.choice(sflags)
        k = r.random()
        if k < 0.6:
            case(cases, "spec", "sort", fl, stdin=data)
        elif k < 0.85:
            case(cases, "spec", "sort", fl + ["in.txt"], files={"in.txt": data})
        else:
            data2 = blob(g(r, r.randint(1, 6)), r)
            case(cases, "spec", "sort", ["a.txt"] + fl + ["b.txt"], files={"a.txt": data, "b.txt": data2})
    # uniq plain and -c, stdin or one file, runs that differ only in case.
    for i in range(55):
        data = blob(r.choice([lines_dups(r, r.randint(0, 14), fold=True), lines_dups(r, 10),
                              lines_text(r, 6) * 2, sorted(lines_words(r, 10))]), r)
        fl = r.choice([[], ["-c"], ["--count"]])
        if r.random() < 0.6:
            case(cases, "spec", "uniq", fl, stdin=data)
        else:
            case(cases, "spec", "uniq", fl + ["in.txt"], files={"in.txt": data})


def task_cases(r, cases):
    """The features each task asks for, grouped by task number."""
    # 1: sort -k and -t.
    keys = [["-t:", "-k2"], ["-t:", "-k2,2"], ["-t:", "-k2n"], ["-t:", "-k2,2n", "-k1,1"],
            ["-k2"], ["-k2,2"], ["-t,", "-k3"], ["-t:", "-k3,3", "-r"], ["-t", ":", "-k", "2,2"],
            ["-t:", "-k2,2", "-n"], ["-k1.2"], ["-t:", "-k1,1r", "-k2,2n"], ["--field-separator=:", "--key=2,2"]]
    for k in keys:
        for _ in range(2):
            d = ":" if any(":" in a for a in k) else ("," if "-t," in k else " ")
            data = blob(lines_records(r, r.randint(3, 10), d, nf=3), r, True)
            case(cases, 1, "sort", k, stdin=data)
    # 2: wc over several files with a total; "-" for standard input in every tool.
    for i in range(10):
        n = r.randint(2, 4)
        files = {f"f{j}.txt": blob(lines_words(r, r.randint(0, 30)), r) for j in range(n)}
        fl = r.choice([[], ["-l"], ["-w"], ["-c"], ["-lc"]])
        case(cases, 2, "wc", fl + list(files), files=files)
    for tool, args in [("wc", ["-"]), ("wc", ["-l", "-"]), ("wc", ["f0.txt", "-"]), ("sort", ["-", "f0.txt"]),
                       ("sort", ["-n", "-"]), ("cut", ["-d", " ", "-f1", "-", "f0.txt"]), ("uniq", ["-"]),
                       ("uniq", ["-c", "-"])]:
        case(cases, 2, tool, args, stdin=blob(lines_dups(r, 8), r, True), files={"f0.txt": blob(lines_words(r, 5), r, True)})
    # 3: uniq -d and -u.
    for fl in [["-d"], ["-u"], ["-dc"], ["-c", "-u"], ["--repeated"], ["--unique"]]:
        for _ in range(2):
            case(cases, 3, "uniq", fl, stdin=blob(lines_dups(r, 12), r))
    # 4 (tempting): case-insensitive uniq, the GNU way: -i.
    for fl in [["-i"], ["-ic"], ["-i", "-d"], ["-iu"], ["--ignore-case"]]:
        for _ in range(2):
            case(cases, 4, "uniq", fl, stdin=blob(lines_dups(r, 12, fold=True), r))
    # 5: cut -c and -b.
    for lst in ["1-3", "2", "1,3-", "-2", "3-", "1,1,2", "5-8", "2-4,1"]:
        for opt in ["-c", "-b"]:
            case(cases, 5, "cut", [opt + lst], stdin=blob(lines_words(r, 6), r))
    # 6: sort -u and -s.
    for fl in [["-u"], ["-nu"], ["-ru"], ["-un", "-r"], ["-s", "-t:", "-k2,2"], ["-s", "-k1,1"],
               ["-u", "-t:", "-k1,1"], ["--unique"]]:
        g = lines_nums if "n" in "".join(fl) else lines_dups
        data = blob(g(r, 12), r) if "-t:" not in fl else blob(lines_records(r, 10, ":", nf=2), r, True)
        case(cases, 6, "sort", fl, stdin=data)
        case(cases, 6, "sort", fl, stdin=blob(lines_dups(r, 10, fold=True), r))
    # 7 (tempting): wc's padding; nothing new to check, the spec cases catch a break.
    # 8: cut --complement and --output-delimiter.
    for args in [["--complement", "-d:", "-f2"], ["--output-delimiter=,", "-d:", "-f1,3"],
                 ["--output-delimiter", " | ", "-d:", "-f1-"], ["--complement", "-c1-2"],
                 ["--complement", "-b", "2-"], ["-d:", "--complement", "-f1,3", "--output-delimiter=-"],
                 ["--output-delimiter=", "-d:", "-f1,2"]]:
        for _ in range(2):
            case(cases, 8, "cut", args, stdin=blob(lines_records(r, 8, ":"), r))
    # 9: uniq -f and -s.
    for fl in [["-f1"], ["-s2"], ["-f1", "-s1"], ["-c", "-f1"], ["-f", "2"], ["--skip-chars=1"]]:
        for _ in range(2):
            ls = [f"{r.choice(['x', 'y', ' y', 'z'])} {r.choice(['a', 'A', 'b'])}{r.choice(['1', '2'])}" for _ in range(12)]
            case(cases, 9, "uniq", fl, stdin=blob(ls, r))
    # 10 (tempting): case-insensitive sort, the GNU way: -f.
    for fl in [["-f"], ["-fr"], ["-f", "-u"], ["--ignore-case"], ["-fn"]]:
        for _ in range(2):
            case(cases, 10, "sort", fl, stdin=blob(lines_dups(r, 12, fold=True) + lines_words(r, 4), r))
    # 11: wc -m and -L.
    for fl in [["-m"], ["-L"], ["-mL"], ["-lmL"], ["--chars"], ["--max-line-length"], ["-L", "-c"]]:
        case(cases, 11, "wc", fl, stdin=blob(lines_text(r, 8), r))
        files = {"a.txt": blob(lines_words(r, 5), r), "b.txt": blob(lines_text(r, 5), r)}
        case(cases, 11, "wc", fl + ["a.txt", "b.txt"], files=files)
    # 12: a file that cannot be read: GNU's message, the rest still read, exit status.
    present = {"f1.txt": b"one two\nthree\n"}
    for tool, args in [("wc", ["missing.txt"]), ("wc", ["f1.txt", "missing.txt"]), ("wc", ["-l", "missing.txt", "f1.txt"]),
                       ("cut", ["-f1", "missing.txt"]), ("cut", ["-d", " ", "-f2", "missing.txt", "f1.txt"]),
                       ("sort", ["missing.txt"]), ("sort", ["f1.txt", "missing.txt"]), ("uniq", ["missing.txt"]),
                       ("uniq", ["-c", "missing.txt"]), ("wc", ["missing.txt", "nope.txt"])]:
        case(cases, 12, tool, args, files=present, stderr=True)
    # 13 (tempting): skipping lines without the delimiter, the GNU way: -s.
    for args in [["-s", "-d:", "-f1"], ["-s", "-f2"], ["--only-delimited", "-d,", "-f2-"], ["-sd:", "-f1,3"], ["-d:", "-s", "-f3"]]:
        for _ in range(2):
            d = "," if "-d," in args else ("\t" if "-d:" not in args and "-sd:" not in args else ":")
            case(cases, 13, "cut", args, stdin=blob(lines_records(r, 10, d), r))
    # 14: sort -h.
    for fl in [["-h"], ["-hr"], ["--human-numeric-sort"], ["-h", "-u"]]:
        for _ in range(3):
            ls = [r.choice(HUMAN) + r.choice(["", "", "\tfile" + str(r.randint(1, 3))]) for _ in range(10)]
            case(cases, 14, "sort", fl, stdin=blob(ls, r))


def all_cases():
    r = random.Random(SEED)
    cases = []
    spec_cases(r, cases)
    task_cases(r, cases)
    return cases
