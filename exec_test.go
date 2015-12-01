package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	dir, err := ioutil.TempDir("", "gom")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	f, err := ioutil.TempFile(dir, "gom")
	if err != nil {
		t.Fatal(err)
	}
	vendor := filepath.Join(dir, vendorFolder)
	err = os.MkdirAll(vendorSrc(vendor), 0755)
	if err != nil {
		t.Fatal(err)
	}
	oldstdout := stdout
	defer func() {
		stdout = oldstdout
	}()
	stdout = f
	err = run([]string{"go", "env"}, None)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	stdout = oldstdout
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	gopath := ""
	for _, line := range strings.Split(string(b), "\n") {
		if runtime.GOOS == "windows" {
			item := strings.SplitN(line, " ", 2)
			if len(item) < 2 {
				continue
			}
			if strings.HasPrefix(item[1], "GOPATH=") {
				gopath = item[1][7:]
			}
		} else if strings.HasPrefix(line, "GOPATH=") {
			gopath, _ = strconv.Unquote(line[7:])
		}
	}
	found := false
	vendorInfo, _ := os.Stat(vendor)
	for _, s := range strings.Split(gopath, string(filepath.ListSeparator)) {
		currentInfo, _ := os.Stat(s)
		if os.SameFile(vendorInfo, currentInfo) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Expected %v, but %v:", vendor, gopath)
	}
}
