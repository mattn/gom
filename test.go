package main

import (
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"path/filepath"
)

func test() error {
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	ct.ChangeColor(ct.Cyan, true, ct.None, false)
	args := []string{"go", "test", "-x"}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ct.ChangeColor(ct.Blue, true, ct.None, false)
	err = cmd.Run()
	ct.ResetColor()
	if err != nil {
		return err
	}
	return nil
}
