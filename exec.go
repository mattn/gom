package main

import (
	"github.com/daviddengcn/go-colortext"
	"os"
	osexec "os/exec"
	"path/filepath"
)

type Color int

const (
	None Color = Color(ct.None)
	Red Color = Color(ct.Red)
	Blue Color = Color(ct.Blue)
)

func ready() error {
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	return nil
}

func gom_exec(args []string, c Color) error {
	if err := ready(); err != nil {
		return err
	}
	if len(args) == 0 {
		usage()
	}
	cmd := osexec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	err := cmd.Run()
	ct.ResetColor()
	if cmd.Process == nil {
		return err
	}
	return nil
}

func run(args []string) error {
	if err := ready(); err != nil {
		return err
	}
	cmdArgs := []string{"go", "run"}
	return gom_exec(append(cmdArgs, args...), None)
}
