package main

import (
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"path/filepath"
)

func build(args []string) error {
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	ct.ChangeColor(ct.Cyan, true, ct.None, false)
	cmdArgs := []string{"go", "build"}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
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
