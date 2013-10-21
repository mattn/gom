package main

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"sort"
	"strings"
)

const travis_yml = ".travis.yml"

func genTravisYml() error {
	_, err := os.Stat(travis_yml)
	if err == nil {
		return errors.New(".travis.yml already exists")
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

// http://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=go1.1.2#96
func isStandardImport(path string) bool {
	return !strings.Contains(path, ".")
}

func appendPkg(pkgs []string, pkg string) []string {
	for _, ele := range pkgs {
		if ele == pkg {
			return pkgs
		}
	}
	return append(pkgs, pkg)
}

func appendPkgs(pkgs, more []string) []string {
	for _, pkg := range more {
		pkgs = appendPkg(pkgs, pkg)
	}
	return pkgs
}

func scanDirectory(path, srcDir string) (ret []string, err error) {
	pkg, err := build.Import(path, srcDir, build.AllowBinary)
	if err != nil {
		return ret, err
	}

	for _, imp := range pkg.Imports {
		switch {
		case isStandardImport(imp):
			// Ignore standard packages
		case !build.IsLocalImport(imp):
			// Add the external package
			ret = appendPkg(ret, imp)
			fallthrough
		default:
			// Does the recursive walk
			pkgs, err := scanDirectory(imp, pkg.Dir)
			if err != nil {
				return ret, err
			}
			ret = appendPkgs(ret, pkgs)
		}
	}

	return ret, err
}

func genGomfile() error {
	_, err := os.Stat("Gomfile")
	if err == nil {
		return errors.New("Gomfile already exists")
	}
	f, err := os.Create("Gomfile")
	if err != nil {
		return err
	}
	defer f.Close()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	all, err := scanDirectory(".", dir)
	if err != nil {
		return err
	}
	sort.Strings(all)
	for _, pkg := range all {
		fmt.Fprintf(f, "gom '%s'\n", pkg)
	}
	return nil
}
