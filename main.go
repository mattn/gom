package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf(`Usage of %s:
 Tasks:
   gom build   [options]   : Build with vendor packages
   gom install [options]   : Install bundled packages into vendor directory
   gom test    [options]   : Run tests with bundles
   gom run     [options]   : Run go file with bundles
   gom doc     [options]   : Run godoc for bundles
   gom exec    [arguments] : Execute command with bundle environment
   gom gen travis-yml      : Generate .travis.yml which uses "gom test"
   gom gen gomfile         : Scan packages from current directory as root
                              recursively, and generate Gomfile
`, os.Args[0])
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}
	handleSignal()

	var err error
	switch flag.Arg(0) {
	case "install", "i":
		err = install(flag.Args()[1:])
	case "build", "b":
		err = gom_exec(append([]string{"go", "build"}, flag.Args()[1:]...), None)
	case "test", "t":
		err = gom_exec(append([]string{"go", "test"}, flag.Args()[1:]...), None)
	case "run", "r":
		err = gom_exec(append([]string{"go", "run"}, flag.Args()[1:]...), None)
	case "doc", "d":
		err = gom_exec(append([]string{"godoc"}, flag.Args()[1:]...), None)
	case "exec", "e":
		err = gom_exec(flag.Args()[1:], None)
	case "gen", "g":
		switch flag.Arg(1) {
		case "travis-yml":
			err = gen_travis_yml()
		case "gomfile":
			err = gen_gomfile()
		default:
			usage()
		}
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "gom: ", err)
		os.Exit(1)
	}
}
