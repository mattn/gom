package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func has(kv map[string]string, key string) bool {
	_, ok := kv[key]
	return ok
}

func checkout(repo string, commit_or_branch_or_tag string) error {
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	p := filepath.Join(vendor, "src")
	for _, elem := range strings.Split(repo, "/") {
		p = filepath.Join(p, elem)
		if isDir(filepath.Join(p, ".git")) {
			err = vcsExec(p, "git", "checkout", "-q", commit_or_branch_or_tag)
			if err != nil {
				return err
			}
			return vcsExec(p, "go", "install")
		} else if isDir(filepath.Join(p, ".hg")) {
			err = vcsExec(p, "hg", "update", commit_or_branch_or_tag)
			if err != nil {
				return err
			}
			return vcsExec(p, "go", "install")
		}
	}
	return errors.New("gom currently support git/hg for specifying tag/branch/commit")
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
	goms, err := parseGomfile("Gomfile")
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
	for _, gom := range goms {
		cmdArgs := []string{"go", "get"}
		fmt.Printf("installing %s\n", gom.name)
		cmdArgs = append(cmdArgs, args...)
		cmdArgs = append(cmdArgs, gom.name)
		err = run(cmdArgs, Blue)
		if err != nil {
			return err
		}
	}
	for _, gom := range goms {
		commit_or_branch_or_tag := ""
		if has(gom.options, "branch") {
			commit_or_branch_or_tag = gom.options["branch"]
		}
		if has(gom.options, "tag") {
			commit_or_branch_or_tag = gom.options["tag"]
		}
		if has(gom.options, "commit") {
			commit_or_branch_or_tag = gom.options["commit"]
		}
		if commit_or_branch_or_tag != "" {
			err = checkout(gom.name, commit_or_branch_or_tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
