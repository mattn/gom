package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
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
		ct.ChangeColor(ct.Cyan, true, ct.None, false)
		fmt.Printf("installing %s(tag=%s)\n",
			gom.name,
			gom.tag)
		ct.ResetColor()
		cmdArgs := []string{"go", "get"}
		if gom.tag != "" {
			cmdArgs = append(cmdArgs, "-tags", gom.tag)
		}
		cmdArgs = append(cmdArgs, args...)
		cmdArgs = append(cmdArgs, gom.name)
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ct.ChangeColor(ct.Blue, true, ct.None, false)
		err = cmd.Run()
		ct.ResetColor()
		if err != nil {
			return err
		}
	}
	return nil
}
