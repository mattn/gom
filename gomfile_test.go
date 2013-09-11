package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func tempGomfile(content string) (string, error) {
	f, err := ioutil.TempFile("", "gom")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}
	name := f.Name()
	return name, nil
}

func TestGomfile1(t *testing.T) {
	filename, err := tempGomfile(`
gom 'github.com/mattn/go-sqlite3', '>3.33'
`)
	if err != nil {
		t.Fatal(err)
	}
	goms, err := parseGomfile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/go-sqlite3", tag: ">3.33", options: make(map[string]string)},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile2(t *testing.T) {
	filename, err := tempGomfile(`
gom 'github.com/mattn/go-sqlite3', '>3.33'
gom 'github.com/mattn/go-gtk'
`)
	if err != nil {
		t.Fatal(err)
	}
	goms, err := parseGomfile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/go-sqlite3", tag: ">3.33", options: make(map[string]string)},
		{name: "github.com/mattn/go-gtk", tag: "", options: make(map[string]string)},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile3(t *testing.T) {
	filename, err := tempGomfile(`
gom 'github.com/mattn/go-sqlite3', '3.14', :commit => 'asdfasdf'
gom 'github.com/mattn/go-gtk', :foobar => 'barbaz'
`)
	if err != nil {
		t.Fatal(err)
	}
	goms, err := parseGomfile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/go-sqlite3", tag: "3.14", options: map[string]string{"commit": "asdfasdf"}},
		{name: "github.com/mattn/go-gtk", tag: "", options: map[string]string{"foobar": "barbaz"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}
