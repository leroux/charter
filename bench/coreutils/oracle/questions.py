"""The understandability questions and how each is answered.

Each question names a command line and its input; the right answer is
what GNU coreutils 9.4 prints for it with LC_ALL=C (or a measure of that:
its byte count). Run `python3 questions.py` to rewrite ../questions.jsonl,
the part the agent sees: only the id and the question.
"""

import json
from pathlib import Path

HOW = (" Assume LC_ALL=C. Answer with exactly what it prints on standard output, "
       "spaces exact, each newline written as \\n.")

# id, tool, args, stdin, files, measure ("stdout" or "bytes")
Q = [
    ("q01", "sort", [], "banana\nApple\napple\n_x\nB\n", {}, "stdout"),
    ("q02", "sort", ["-n"], "10\n9\n-3\n2.5\nabc\n", {}, "stdout"),
    ("q03", "sort", ["-n"], "2.50\n2.5\n02.5\n", {}, "stdout"),
    ("q04", "sort", ["-rn"], "b 1\na 1\nc 2\n", {}, "stdout"),
    ("q05", "sort", ["-n"], "+5\n-0\n0\n-1\n", {}, "stdout"),
    ("q06", "uniq", [], "a\na\nA\na\nb\nb\n", {}, "stdout"),
    ("q07", "uniq", ["-c"], "x\nx\ny\nx\n", {}, "stdout"),
    ("q08", "uniq", ["-c"], "a\na", {}, "stdout"),
    ("q09", "wc", [], "one two\nthree", {}, "stdout"),
    ("q10", "wc", ["-w"], "a\x01b \x01 c\n", {}, "stdout"),
    ("q11", "wc", ["-l"], "a\nb\nc", {}, "stdout"),
    ("q12", "wc", ["f.txt"], "", {"f.txt": "hello world\n"}, "stdout"),
    ("q13", "cut", ["-d:", "-f3,1"], "a:b:c\nplain\n", {}, "stdout"),
    ("q14", "cut", ["-d:", "-f5"], "a:b\n", {}, "bytes"),
    ("q15", "cut", ["-f1"], "a b\tc\n", {}, "stdout"),
    ("q16", "cut", ["-d,", "-f2-"], "x,y,z,\n", {}, "stdout"),
    ("q17", "sort", ["a.txt", "b.txt"], "", {"a.txt": "b", "b.txt": "a\n"}, "stdout"),
    ("q18", "sort", ["-t:", "-k2,2n"], "x:3\ny:10\nz:2\n", {}, "stdout"),
    ("q19", "wc", ["-l", "a.txt", "b.txt"], "", {"a.txt": "1\n2\n", "b.txt": "3\n"}, "stdout"),
    ("q20", "uniq", ["-d"], "a\na\nb\nc\nc\n", {}, "stdout"),
]


def text(tool, args, stdin, files, measure):
    cmd = " ".join([tool, *args])
    if files:
        desc = "; ".join(f"`{n}` holds the bytes {json.dumps(d)}" for n, d in files.items())
        s = f"In a directory where {desc}, what does `{cmd}` print?"
    else:
        s = f"What does `{cmd}` print when its standard input is the bytes {json.dumps(stdin)}?"
    if measure == "bytes":
        return s + " Assume LC_ALL=C. Answer with the number of bytes it prints on standard output, a bare number."
    return s + HOW


def questions():
    return [{"id": i, "tool": t, "args": a, "stdin": s.encode("latin-1"),
             "files": {k: v.encode("latin-1") for k, v in f.items()}, "measure": m,
             "question": text(t, a, s, f, m)} for i, t, a, s, f, m in Q]


if __name__ == "__main__":
    out = Path(__file__).resolve().parent.parent / "questions.jsonl"
    out.write_text("".join(json.dumps({"id": q["id"], "question": q["question"]}) + "\n" for q in questions()))
