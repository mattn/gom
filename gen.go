package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

const travis_yml = ".travis.yml"

func genTravisYml() error {
	_, err := os.Stat(travis_yml)
	if err == nil {
		return errors.New(".travis.yml is already exists")
	}
	f, err := os.Create(travis_yml)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(`language: go
go:
  - tip
before_install:
  - go get github.com/mattn/gom
script:
  - $HOME/gopath/bin/gom install
  - $HOME/gopath/bin/gom test
`)
	return nil
}

func scanPackages(filename string) (ret []string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return
	}
	ast.SortImports(fset, f)
	goroot := os.Getenv("GOROOT")
	for _, imp := range f.Imports {
		pkg := unquote(imp.Path.Value)
		p := filepath.Join(goroot, "src", "pkg", pkg)
		if _, err = os.Stat(p); err != nil {
			ret = appendPkg(ret, imp.Path.Value)
		}
	}
	return ret
}

func appendPkg(pkgs []string, pkg string) []string {
    for _, ele := range pkgs {
        if ele == pkg {
            return pkgs
        }
    }
    return append(pkgs, pkg)
}

func genGomfile() error {
	_, err := os.Stat("Gomfile")
	if err == nil {
		return errors.New("Gomfile is already exists")
	}
	f, err := os.Create("Gomfile")
	if err != nil {
		return err
	}
	defer f.Close()

	all := []string{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		_, last := filepath.Split(path)
		if last == "vendor" && info.IsDir() {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".go" {
			for _, pkg := range scanPackages(path) {
				all = appendPkg(all, unquote(pkg))
			}
		}
		return nil
	})
	sort.Strings(all)
	for _, pkg := range all {
		fmt.Fprintf(f, "gom '%s'\n", pkg)
	}
	return nil
}
