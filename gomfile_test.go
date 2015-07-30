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
gom 'github.com/mattn/go-sqlite3', :tag => '>3.33'
`)
	if err != nil {
		t.Fatal(err)
	}
	goms, err := parseGomfile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/go-sqlite3", options: map[string]interface{}{"tag": ">3.33"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile2(t *testing.T) {
	filename, err := tempGomfile(`
gom 'github.com/mattn/go-sqlite3', :tag => '>3.33'
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
		{name: "github.com/mattn/go-sqlite3", options: map[string]interface{}{"tag": ">3.33"}},
		{name: "github.com/mattn/go-gtk", options: map[string]interface{}{}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile3(t *testing.T) {
	filename, err := tempGomfile(`
gom 'github.com/mattn/go-sqlite3', :tag => '3.14', :commit => 'asdfasdf'
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
		{name: "github.com/mattn/go-sqlite3", options: map[string]interface{}{"tag": "3.14", "commit": "asdfasdf"}},
		{name: "github.com/mattn/go-gtk", options: map[string]interface{}{"foobar": "barbaz"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile4(t *testing.T) {
	filename, err := tempGomfile(`
group :development do
	gom 'github.com/mattn/go-sqlite3', :tag => '3.14', :commit => 'asdfasdf'
end

group :test do
	gom 'github.com/mattn/go-gtk', :foobar => 'barbaz'
end
`)
	if err != nil {
		t.Fatal(err)
	}

	*developmentEnv = true
	goms, err := parseGomfile(filename)
	*developmentEnv = false

	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/go-sqlite3", options: map[string]interface{}{"tag": "3.14", "commit": "asdfasdf"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile5(t *testing.T) {
	filename, err := tempGomfile(`
group :custom_one do
	gom 'github.com/mattn/go-sqlite3', :tag => '3.14', :commit => 'asdfasdf'
end
group :custom_two do
	gom 'github.com/mattn/go-gtk', :foobar => 'barbaz'
end
`)
	if err != nil {
		t.Fatal(err)
	}

	customGroupList = []string{"custom_one", "custom_two"}
	goms, err := parseGomfile(filename)

	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/go-sqlite3", options: map[string]interface{}{"tag": "3.14", "commit": "asdfasdf"}},
		{name: "github.com/mattn/go-gtk", options: map[string]interface{}{"foobar": "barbaz"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}

	customGroupList = []string{"custom_one"}
	goms, err = parseGomfile(filename)

	if err != nil {
		t.Fatal(err)
	}
	expected = []Gom{
		{name: "github.com/mattn/go-sqlite3", options: map[string]interface{}{"tag": "3.14", "commit": "asdfasdf"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}

func TestGomfile99(t *testing.T) {
	filename, err := tempGomfile(`
gom 'github.com/mattn/gom', :fork => 'github.com/dicefm/gom'
`)
	if err != nil {
		t.Fatal(err)
	}
	goms, err := parseGomfile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Gom{
		{name: "github.com/mattn/gom", options: map[string]interface{}{"fork": "github.com/dicefm/gom"}},
	}
	if !reflect.DeepEqual(goms, expected) {
		t.Fatalf("Expected %v, but %v:", expected, goms)
	}
}
