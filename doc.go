package main

import (
	"os"
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
	return exec(cmdArgs)
}
