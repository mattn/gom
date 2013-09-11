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
	fmt.Println("   gom build")
	fmt.Println("   gom install")
	fmt.Println("   gom test")
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
	case "install":
		err = install()
	case "build":
		err = build()
	case "test":
		err = test()
	case "gen":
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
