package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	goms, err := parseGomfile("Gomfile")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gom: ", err)
		os.Exit(1)
	}

	if flag.Arg(0) == "install" {
		err = install(goms)
		if err != nil {
			fmt.Fprintln(os.Stderr, "gom: ", err)
			os.Exit(1)
		}
	}
}
