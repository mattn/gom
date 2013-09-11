package main

import (
	"github.com/daviddengcn/go-colortext"
	"os"
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
	ct.ChangeColor(ct.Blue, true, ct.None, false)
	err = exec(cmdArgs)
	ct.ResetColor()
	return err
}
