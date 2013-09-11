package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	flag.Usage()
	fmt.Println(" Tasks:")
	fmt.Println("   gom build   [options]")
	fmt.Println("   gom install [options]")
	fmt.Println("   gom test    [options]")
	fmt.Println("   gom run     [options]")
	fmt.Println("   gom doc     [options]")
	fmt.Println("   gom exec    [arguments]")
	fmt.Println("   gom gen travis-yml")
	fmt.Println("   gom gen gomfile")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

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
