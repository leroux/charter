import os
import subprocess
import sys
import unittest
from pathlib import Path

import claims
from slugs import env
from slugs.slug import separator, slug

SVC = Path(__file__).resolve().parents[2] / "svc"


class SlugTest(unittest.TestCase):
    def test_lower_and_joined(self):
        claims.verifies(self, claims.SW_261126_4NQQ, "lower-and-joined", "no-edge-separator")
        self.assertEqual(slug("Hello, World!"), "hello-world")
        self.assertEqual(slug("  --Go--  "), "go")

    def test_separator_from_env(self):
        claims.verifies(self, claims.SW_261126_X6HB, "separator-from-env")
        os.environ[env.SEPARATOR] = "_"
        try:
            self.assertEqual(slug("Hello World", separator()), "hello_world")
        finally:
            del os.environ[env.SEPARATOR]

    def test_same_slug_both_ways(self):
        claims.verifies(self, claims.STK_261126_78RN, "same-slug-both-ways")
        title = "It's a Brand-New Day, 2026!"
        prog = 'package main\nimport ("fmt"; "os"; "example.com/slugs/svc/slug")\nfunc main() { fmt.Print(slug.Of(os.Args[1])) }\n'
        d = SVC / "cmd" / "slugtmp"
        d.mkdir(parents=True, exist_ok=True)
        try:
            (d / "main.go").write_text(prog)
            out = subprocess.run(["go", "run", "./cmd/slugtmp", title], cwd=SVC, capture_output=True, text=True, check=True)
        finally:
            (d / "main.go").unlink(missing_ok=True)
            d.rmdir()
        self.assertEqual(out.stdout, slug(title))


if __name__ == "__main__":
    unittest.main()
