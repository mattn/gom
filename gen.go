package main

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
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

func vcsScan(p, target string) (*vcsCmd, string, string) {
	name := ""
	for _, elem := range strings.Split(target, "/") {
		var vcs *vcsCmd
		p = filepath.Join(p, elem)
		if name == "" {
			name = elem
		} else {
			name += `/` + elem
		}
		if isDir(filepath.Join(p, ".git")) {
			vcs = git
		} else if isDir(filepath.Join(p, ".hg")) {
			vcs = hg
		} else if isDir(filepath.Join(p, ".bzr")) {
			vcs = bzr
		}
		if vcs != nil {
			return vcs, name, p
		}
	}
	return nil, "", ""
}

func genGomfile() error {
	_, err := os.Stat("Gomfile")
	if err == nil {
		return errors.New("Gomfile already exists")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	all, err := scanDirectory(".", dir)
	if err != nil {
		return err
	}
	sort.Strings(all)
	goms := make([]Gom, 0)
	for _, pkg := range all {
		for _, p := range filepath.SplitList(os.Getenv("GOPATH")) {
			var vcs *vcsCmd
			var n string
			vcs, n, p = vcsScan(filepath.Join(p, "src"), pkg)
			if vcs != nil {
				found := false
				for _, gom := range goms {
					if gom.name == n {
						found = true
						break
					}
				}
				if !found {
					gom := Gom{name: n, options: make(map[string]interface{})}
					rev, err := vcs.Revision(p)
					if err == nil && rev != "" {
						gom.options["commit"] = rev
					}
					goms = append(goms, gom)
				}
			}
		}
	}

	return writeGomfile("Gomfile", goms)
}

func genGomfileLock() error {
	allGoms, err := parseGomfile("Gomfile")
	if err != nil {
		return err
	}
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	goms := make([]Gom, 0)
	for _, gom := range allGoms {
		if group, ok := gom.options["group"]; ok {
			if !matchEnv(group) {
				continue
			}
		}
		if goos, ok := gom.options["goos"]; ok {
			if !matchOS(goos) {
				continue
			}
		}
		goms = append(goms, gom)
	}

	for _, gom := range goms {
		var vcs *vcsCmd
		var p string
		vcs, _, p = vcsScan(vendorSrc(vendor), gom.name)
		if vcs != nil {
			rev, err := vcs.Revision(p)
			if err == nil && rev != "" {
				gom.options["commit"] = rev
			}
		}
	}
	err = writeGomfile("Gomfile.lock", goms)
	if err == nil {
		fmt.Println("Gomfile.lock is generated")
	}
	return err
}
