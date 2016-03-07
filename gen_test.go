package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGomfile(t *testing.T) {
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
	err = os.MkdirAll(filepath.Join(vendor, "src"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	oldstdout := stdout
	defer func() {
		stdout = oldstdout
	}()
	stdout = f
	err = run([]string{"go", "env", "GOPATH"}, None)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	stdout = oldstdout
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	gopath := string(b)
	found := false
	for _, s := range strings.Split(gopath, string(filepath.ListSeparator)) {
		if filepath.Clean(s) == filepath.Clean(vendor) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Expected %v, but %v:", vendor, found)
	}
}
