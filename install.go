package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
		err = os.MkdirAll(vendor, 755)
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
		if gom.tag != "" {
			cmdArgs = append(cmdArgs, "-tags", gom.tag)
			fmt.Printf("installing %s(tag=%s)\n", gom.name, gom.tag)
		} else {
			fmt.Printf("installing %s\n", gom.name)
		}
		cmdArgs = append(cmdArgs, args...)
		cmdArgs = append(cmdArgs, gom.name)
		err = gom_exec(cmdArgs, Blue)
		if err != nil {
			return err
		}
	}
	return nil
}
