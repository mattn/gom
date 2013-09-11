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
		for _, gom := range goms {
			fmt.Printf("installing %s(tag: %s, options: %s)\n",
				gom.name,
				gom.tag,
				gom.options)
		}
	}
}
