package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func doc(args []string) error {
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	cmdArgs := []string{"godoc", "-goroot", vendor}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
