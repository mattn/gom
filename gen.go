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

type importPackages []importPackage
type importPackage struct {
	path       string
	isTestFile bool
}

func (slice importPackages) Len() int {
	return len(slice)
}

func (slice importPackages) Less(i, j int) bool {
	return slice[i].path < slice[j].path
}

func (slice importPackages) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

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

func appendPkg(pkgs []importPackage, pkg string) ([]importPackage, bool) {
	for _, ele := range pkgs {
		if ele.path == pkg {
			return pkgs, false
		}
	}
	return append(pkgs, importPackage{path: pkg}), true
}

func appendPkgs(pkgs, more []importPackage) []importPackage {
	for _, pkg := range more {
		pkgs, _ = appendPkg(pkgs, pkg.path)
	}
	return pkgs
}

func scanDirectory(path, srcDir string) (ret []importPackage, err error) {
	pkg, err := build.Import(path, srcDir, build.AllowBinary)
	if err != nil {
		return ret, err
	}
	for _, imp := range pkg.Imports {
		switch {
		case pkg.Goroot:
			// Ignore standard packages
		case !build.IsLocalImport(imp):
			// Add the external package
			ret, _ = appendPkg(ret, imp)
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
	retTests := []importPackage{}
	isAdd := false
	for _, imp := range pkg.TestImports {
		switch {
		case pkg.Goroot:
			// Ignore standard packages
			break
		case !build.IsLocalImport(imp):
			// Add the external package
			retTests, isAdd = appendPkg(retTests, imp)
			if isAdd {
				retTests[len(retTests)-1].isTestFile = true
			}
		}
	}
	ret = append(ret, retTests...)
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
	sort.Sort(importPackages(all))
	goms := make([]Gom, 0)
	for _, pkg := range all {
		for _, p := range filepath.SplitList(os.Getenv("GOPATH")) {
			var vcs *vcsCmd
			var n string
			vcs, n, p = vcsScan(filepath.Join(p, "src"), pkg.path)
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
					if pkg.isTestFile {
						gom.options["group"] = "test"
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
