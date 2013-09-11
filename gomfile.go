package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var qx = `'[^']*'|"[^"]*"`
var re1 = regexp.MustCompile(`^\s*gom\s+(` + qx + `)\s*$`)
var re2 = regexp.MustCompile(`^\s*gom\s+(` + qx + `)\s*((?:,\s*:[a-zA-Z][a-z0-9_]*\s=>\s*` + qx + `)+)$`)
var reOptions = regexp.MustCompile(`(,\s*:[a-zA-Z][a-z0-9_]*\s=>\s*` + qx + `)`)

func unquote(name string) string {
	name = strings.TrimSpace(name)
	unquoted, err := strconv.Unquote(name)
	if err != nil {
		return name[1:len(name)-1]
	}
	return unquoted
}

func parseOptions(line string, options map[string]string) {
	ss := reOptions.FindAllStringSubmatch(line, -1)
	for _, s := range ss {
		kvs := strings.Split(strings.TrimSpace(s[0])[1:], "=>")
		options[strings.TrimSpace(kvs[0])[1:]] = unquote(kvs[1])
	}
}

type Gom struct {
	name    string
	options map[string]string
}

func parseGomfile(filename string) ([]Gom, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(f)

	goms := make([]Gom, 0)

	n := 0
	for {
		n++
		lb, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				return goms, nil
			}
			return nil, err
		}
		line := strings.TrimSpace(string(lb))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name := ""
		options := make(map[string]string)
		var items []string
		if re1.MatchString(line) {
			items = re1.FindStringSubmatch(line)[1:]
			name = unquote(items[0])
		} else if re2.MatchString(line) {
			items = re2.FindStringSubmatch(line)[1:]
			name = unquote(items[0])
			parseOptions(items[1], options)
		} else {
			return nil, fmt.Errorf("Failed to parse Gomfile at line %d", n)
		}
		goms = append(goms, Gom{name, options})
	}
	return goms, nil
}
