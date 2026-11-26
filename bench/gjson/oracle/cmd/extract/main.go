// Command extract is the build-time tool that cuts the hidden test files
// out of gjson's upstream test files. It is not run at score time.
//
//	extract changed OLD NEW          the Test functions NEW adds or changes
//	extract keep SRC OUT Name...      SRC's named Test functions, the
//	                                  top-level declarations they reach and
//	                                  the imports those use, written to OUT
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	switch os.Args[1] {
	case "changed":
		changed(os.Args[2], os.Args[3])
	case "keep":
		keep(os.Args[2], os.Args[3], os.Args[4:])
	default:
		fmt.Fprintln(os.Stderr, "usage: extract changed OLD NEW | keep SRC OUT Name...")
		os.Exit(2)
	}
}

func parse(path string) (*token.FileSet, *ast.File, []byte) {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return fs, f, src
}

func testFuncs(path string) map[string]string {
	fs, f, src := parse(path)
	out := map[string]string{}
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || !strings.HasPrefix(fd.Name.Name, "Test") {
			continue
		}
		body := src[fs.Position(fd.Pos()).Offset:fs.Position(fd.End()).Offset]
		out[fd.Name.Name] = strings.Join(strings.Fields(string(body)), " ")
	}
	return out
}

func changed(oldPath, newPath string) {
	old, nw := testFuncs(oldPath), testFuncs(newPath)
	var names []string
	for n, b := range nw {
		if old[n] != b {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	fmt.Println(strings.Join(names, " "))
}

func declNames(d ast.Decl) []string {
	switch d := d.(type) {
	case *ast.FuncDecl:
		if d.Recv == nil {
			return []string{d.Name.Name}
		}
	case *ast.GenDecl:
		var out []string
		for _, s := range d.Specs {
			switch s := s.(type) {
			case *ast.ValueSpec:
				for _, n := range s.Names {
					out = append(out, n.Name)
				}
			case *ast.TypeSpec:
				out = append(out, s.Name.Name)
			}
		}
		return out
	}
	return nil
}

func keep(srcPath, outPath string, names []string) {
	label := os.Getenv("EXTRACT_LABEL")
	if label == "" {
		label = srcPath
	}
	fs, f, _ := parse(srcPath)
	byName := map[string]ast.Decl{}
	for _, d := range f.Decls {
		for _, n := range declNames(d) {
			byName[n] = d
		}
	}
	kept := map[ast.Decl]bool{}
	var queue []ast.Decl
	for _, n := range names {
		d, ok := byName[n]
		if !ok {
			fmt.Fprintf(os.Stderr, "extract: %s has no %s\n", srcPath, n)
			os.Exit(1)
		}
		if !kept[d] {
			kept[d] = true
			queue = append(queue, d)
		}
	}
	used := map[string]bool{} // package names used by selectors
	for len(queue) > 0 {
		d := queue[0]
		queue = queue[1:]
		ast.Inspect(d, func(n ast.Node) bool {
			if se, ok := n.(*ast.SelectorExpr); ok {
				if x, ok := se.X.(*ast.Ident); ok && x.Obj == nil {
					used[x.Name] = true // a package, not a local of that name
				}
			}
			if id, ok := n.(*ast.Ident); ok {
				if dd, ok := byName[id.Name]; ok && !kept[dd] {
					if fd, ok := dd.(*ast.FuncDecl); ok && strings.HasPrefix(fd.Name.Name, "Test") {
						return true
					}
					kept[dd] = true
					queue = append(queue, dd)
				}
			}
			return true
		})
	}
	var decls []ast.Decl
	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			var specs []ast.Spec
			for _, s := range gd.Specs {
				is := s.(*ast.ImportSpec)
				p, _ := strconv.Unquote(is.Path.Value)
				name := p[strings.LastIndex(p, "/")+1:]
				if is.Name != nil {
					name = is.Name.Name
				}
				if used[name] {
					specs = append(specs, is)
				}
			}
			if len(specs) > 0 {
				gd.Specs = specs
				decls = append(decls, gd)
			}
			continue
		}
		if kept[d] {
			decls = append(decls, d)
		}
	}
	f.Decls = decls
	f.Comments = nil
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code cut from upstream %s by the gjson benchmark's extract. DO NOT EDIT.\n\n", label)
	if err := format.Node(&buf, fs, f); err != nil {
		panic(err)
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		panic(err)
	}
}
