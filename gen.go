package main

import (
	"errors"
	"os"
)

const travis_yml = ".travis.yml"

func gen_travis_yml() error {
	_, err := os.Stat(travis_yml)
	if err == nil {
		return errors.New(".travis.yml is already exists")
	}
	f, err := os.Create(travis_yml)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(`language: go
go:
  - tip
before_install:
  - go get github.com/mattn/gom
script:
  - $HOME/gopath/bin/gom install
  - $HOME/gopath/bin/gom test
`)
	return nil
}
