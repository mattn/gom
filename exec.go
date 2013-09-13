package main

import (
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

type Color int

const (
	None Color = Color(ct.None)
	Red  Color = Color(ct.Red)
	Blue Color = Color(ct.Blue)
)

func handleSignal() {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sc
		ct.ResetColor()
		os.Exit(0)
	}()
}

func ready() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	vendor, err := filepath.Abs("vendor")
	if err != nil {
		return err
	}
	for {
		file := filepath.Join(dir, "Gomfile")
		if isFile(file) {
			vendor = filepath.Join(dir, "vendor")
			break
		}
		next := filepath.Clean(filepath.Join(dir, ".."))
		if next == dir {
			break
		}
		dir = next
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	return nil
}

var stdout = os.Stdout
var stderr = os.Stderr

func run(args []string, c Color) error {
	if err := ready(); err != nil {
		return err
	}
	if len(args) == 0 {
		usage()
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	err := cmd.Run()
	ct.ResetColor()
	return err
}
