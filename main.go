package main

import (
	"flag"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/signal"
	"syscall"
)

func usage() {
	flag.Usage()
	fmt.Println(" Tasks:")
	fmt.Println("   gom build [options]")
	fmt.Println("   gom install [options]")
	fmt.Println("   gom test [options]")
	fmt.Println("   gom run [options]")
	fmt.Println("   gom doc [options]")
	fmt.Println("   gom gen travis-yml")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	sc := make(chan os.Signal, 10)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sc
		ct.ResetColor()
		os.Exit(0)
	}()

	var err error
	switch flag.Arg(0) {
	case "install", "i":
		err = install(flag.Args()[1:])
	case "build", "b":
		err = build(flag.Args()[1:])
	case "test", "t":
		err = test(flag.Args()[1:])
	case "run", "r":
		err = run(flag.Args()[1:])
	case "doc", "d":
		err = doc(flag.Args()[1:])
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
