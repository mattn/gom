package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func checkout(repo string, commit_or_branch_or_tag string, args []string) error {
	installCmd := append([]string{"go", "install"}, args...)
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	p := filepath.Join(vendor, "src")
	for _, elem := range strings.Split(repo, "/") {
		p = filepath.Join(p, elem)
		if isDir(filepath.Join(p, ".git")) {
			p = filepath.Join(vendor, "src", repo)
			err = vcsExec(p, "git", "checkout", "-q", commit_or_branch_or_tag)
			if err != nil {
				return err
			}
			return vcsExec(p, installCmd...)
		} else if isDir(filepath.Join(p, ".hg")) {
			p = filepath.Join(vendor, "src", repo)
			err = vcsExec(p, "hg", "update", commit_or_branch_or_tag)
			if err != nil {
				return err
			}
			return vcsExec(p, installCmd...)
		} else if isDir(filepath.Join(p, ".bzr")) {
			p = filepath.Join(vendor, "src", repo)
			err = vcsExec(p, "bzr", "revert", "-r", commit_or_branch_or_tag)
			if err != nil {
				return err
			}
			return vcsExec(p, installCmd...)
		}
	}
	fmt.Printf("Warning: don't know how to checkout for %v", repo)
	return errors.New("gom currently support git/hg/bzr for specifying tag/branch/commit")
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
	err = cmd.Run()
	if cmd.Process == nil {
		return err
	}
	return nil
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
		}

		cmdArgs := []string{"go", "get"}
		cmdArgs = append(cmdArgs, args...)
		cmdArgs = append(cmdArgs, gom.name)

		fmt.Printf("installing %s\n", gom.name)
		err = run(cmdArgs, Blue)
		if err != nil {
			return err
		}
	}
	for _, gom := range goms {
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
		if commit_or_branch_or_tag != "" {
			err = checkout(gom.name, commit_or_branch_or_tag, args)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
