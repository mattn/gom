package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var qx = `'[^']*'|"[^"]*"`
var re_group = regexp.MustCompile(`\s*group\s+((?::[a-z][a-z0-9_]*\s*|,\s*:[a-z][a-z0-9]*\s*)*)\s*do\s*$`)
var re_end = regexp.MustCompile(`\s*end\s*$`)
var re_gom1 = regexp.MustCompile(`^\s*gom\s+(` + qx + `)\s*$`)
var re_gom2 = regexp.MustCompile(`^\s*gom\s+(` + qx + `)\s*((?:,\s*:[a-z][a-z0-9_]*\s=>\s*` + qx + `)+)$`)
var re_options = regexp.MustCompile(`(,\s*:[a-z][a-z0-9_]*\s=>\s*` + qx + `)`)

func unquote(name string) string {
	name = strings.TrimSpace(name)
	if len(name) > 2 {
		if (name[0] == '\'' && name[len(name)-1] == '\'') || (name[0] == '"' && name[len(name)-1] == '"') {
			return name[1 : len(name)-1]
		}
	}
	return name
}

func isEnv(envs []string, env string) bool {
	for _, e := range envs {
		e = strings.TrimSpace(e)
		if e[1:] == env {
			return true
		}
	}
	return false
}

func parseOptions(line string, options map[string]string) {
	ss := re_options.FindAllStringSubmatch(line, -1)
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
	skip := 0
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
		if re_group.MatchString(line) {
			envs := strings.Split(re_group.FindStringSubmatch(line)[1], ",")
			switch {
			case isEnv(envs, "production") && *productionEnv:
				continue
			case isEnv(envs, "development") && *developmentEnv:
				continue
			case isEnv(envs, "test") && *testEnv:
				continue
			}
			skip++
			continue
		} else if re_end.MatchString(line) {
			if skip > 0 {
				skip--
			}
			continue
		} else if skip > 0 {
			continue
		} else if re_gom1.MatchString(line) && skip == 0 {
			items = re_gom1.FindStringSubmatch(line)[1:]
			name = unquote(items[0])
		} else if re_gom2.MatchString(line) && skip == 0 {
			items = re_gom2.FindStringSubmatch(line)[1:]
			name = unquote(items[0])
			parseOptions(items[1], options)
		} else {
			return nil, fmt.Errorf("Failed to parse Gomfile at line %d", n)
		}
		goms = append(goms, Gom{name, options})
	}
	return goms, nil
}
