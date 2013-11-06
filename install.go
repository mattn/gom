package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type vcsCmd struct {
	checkout []string
	update   []string
}

var (
	hg = &vcsCmd{
		[]string{"hg", "update"},
		[]string{"hg", "pull"},
	}
	git = &vcsCmd{
		[]string{"git", "checkout", "-q"},
		[]string{"git", "fetch"},
	}
	bzr = &vcsCmd{
		[]string{"bzr", "revert", "-r"},
		[]string{"bzr", "pull"},
	}
)

func (vcs *vcsCmd) Checkout(p, destination string) error {
	args := append(vcs.checkout, destination)
	return vcsExec(p, args...)
}

func (vcs *vcsCmd) Update(p string) error {
	return vcsExec(p, vcs.update...)
}

func (vcs *vcsCmd) Sync(p, destination string) error {
	err := vcs.Checkout(p, destination)
	if err != nil {
		err = vcs.Update(p)
		if err != nil {
			return err
		}
		err = vcs.Checkout(p, destination)
	}
	return err
}

func vcsExec(dir string, args ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func has(c interface{}, key string) bool {
	if m, ok := c.(map[string]interface{}); ok {
		_, ok := m[key]
		return ok
	} else if a, ok := c.([]string); ok {
		for _, s := range a {
			if ok && s == key {
				return true
			}
		}
	}
	return false
}

func (gom *Gom) Clone(args []string) error {
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	if command, ok := gom.options["command"].(string); ok {
		target, ok := gom.options["target"].(string)
		if !ok {
			target = gom.name
		}

		srcdir := filepath.Join(vendor, "src", target)
		customCmd := strings.Split(command, " ")
		customCmd = append(customCmd, srcdir)

		fmt.Printf("fetching %s (%v)\n", gom.name, customCmd)
		err = run(customCmd, Blue)
		if err != nil {
			return err
		}
	} else if private, ok := gom.options["private"].(string); ok {
		if private == "true" {
			target, ok := gom.options["target"].(string)
			if !ok {
				target = gom.name
			}
			srcdir := filepath.Join(vendor, "src", target)
			if _, err := os.Stat(srcdir); err != nil {
				if os.IsExist(err) {
					if err := gom.pullPrivate(srcdir); err != nil {
						return err
					}
				} else {
					if err := gom.clonePrivate(srcdir); err != nil {
						return err
					}
				}
			}
		}
	}

	cmdArgs := []string{"go", "get", "-d"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, gom.name)

	fmt.Printf("downloading %s\n", gom.name)
	return run(cmdArgs, Blue)
}

func (gom *Gom) pullPrivate(srcdir string) (err error) {
	fmt.Printf("fetching private repo %s\n", gom.name)
	pullCmd := fmt.Sprintf("git --work-tree=%s, --git-dir=%s/.git pull origin",
		srcdir, srcdir)
	pullArgs := strings.Split(pullCmd, " ")
	err = run(pullArgs, Blue)
	if err != nil {
		return
	}

	return
}

func (gom *Gom) clonePrivate(srcdir string) (err error) {
	name := strings.Split(gom.name, "/")
	privateUrl := fmt.Sprintf("git@%s:%s/%s", name[0], name[1], name[2])

	fmt.Printf("fetching private repo %s\n", gom.name)
	cloneCmd := []string{"git", "clone", privateUrl, srcdir}
	err = run(cloneCmd, Blue)
	if err != nil {
		return
	}

	return
}

func (gom *Gom) Checkout() error {
	commit_or_branch_or_tag := ""
	if has(gom.options, "branch") {
		commit_or_branch_or_tag, _ = gom.options["branch"].(string)
	}
	if has(gom.options, "tag") {
		commit_or_branch_or_tag, _ = gom.options["tag"].(string)
	}
	if has(gom.options, "commit") {
		commit_or_branch_or_tag, _ = gom.options["commit"].(string)
	}
	if commit_or_branch_or_tag == "" {
		return nil
	}
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	p := filepath.Join(vendor, "src")
	for _, elem := range strings.Split(gom.name, "/") {
		var vcs *vcsCmd
		p = filepath.Join(p, elem)
		if isDir(filepath.Join(p, ".git")) {
			vcs = git
		} else if isDir(filepath.Join(p, ".hg")) {
			vcs = hg
		} else if isDir(filepath.Join(p, ".bzr")) {
			vcs = bzr
		}
		if vcs != nil {
			p = filepath.Join(vendor, "src", gom.name)
			return vcs.Sync(p, commit_or_branch_or_tag)
		}
	}
	fmt.Printf("Warning: don't know how to checkout for %v\n", gom.name)
	return errors.New("gom currently support git/hg/bzr for specifying tag/branch/commit")
}

func (gom *Gom) Build(args []string) error {
	installCmd := append([]string{"go", "install"}, args...)
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	p := filepath.Join(vendor, "src", gom.name)
	return vcsExec(p, installCmd...)
}

func isFile(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && !fi.IsDir() {
		return true
	}
	return false
}

func isDir(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func install(args []string) error {
	allGoms, err := parseGomfile("Gomfile")
	if err != nil {
		return err
	}
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	_, err = os.Stat(vendor)
	if err != nil {
		err = os.MkdirAll(vendor, 0755)
		if err != nil {
			return err
		}
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}

	// 1. Filter goms to install
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

	// 2. Clone the repositories
	for _, gom := range goms {
		err = gom.Clone(args)
		if err != nil {
			return err
		}
	}

	// 3. Checkout the commit/branch/tag if needed
	for _, gom := range goms {
		err = gom.Checkout()
		if err != nil {
			return err
		}
	}

	// 4. Build and install
	for _, gom := range goms {
		err = gom.Build(args)
		if err != nil {
			return err
		}
	}

	return nil
}
