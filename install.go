package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"path/filepath"
)

func install(goms []Gom) error {
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
		args := []string{"go", "get", "-x"}
		if gom.tag != "" {
			args = append(args, "-tags", gom.tag)
		}
		args = append(args, gom.name)
		cmd := exec.Command(args[0], args[1:]...)
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
