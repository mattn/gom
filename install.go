package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
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
		ct.ChangeColor(ct.Cyan, true, ct.None, false)
		if gom.tag != "" {
			cmdArgs = append(cmdArgs, "-tags", gom.tag)
			fmt.Printf("installing %s(tag=%s)\n", gom.name, gom.tag)
		} else {
			fmt.Printf("installing %s\n", gom.name)
		}
		ct.ResetColor()
		cmdArgs = append(cmdArgs, args...)
		cmdArgs = append(cmdArgs, gom.name)
		ct.ChangeColor(ct.Blue, true, ct.None, false)
		err = exec(cmdArgs)
		ct.ResetColor()
		if err != nil {
			return err
		}
	}
	return nil
}
