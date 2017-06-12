package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kshvmdn/fsql"
)

var options struct {
	version     bool
	interactive bool
}

func readInput() string {
	if len(flag.Args()) > 1 {
		return strings.Join(flag.Args(), " ")
	}

	return flag.Args()[0]
}

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [options] query\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&options.interactive, "interactive", false,
		"run in interactive mode (Ctrl+D to exit)")
	flag.BoolVar(&options.version, "version", false, "print version and exit")
	flag.Parse()

	if options.version {
		fmt.Printf("fsql v%v\n", fsql.Version)
		os.Exit(0)
	}

	if options.interactive {
		if err := fsql.RunInteractive(); err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := fsql.Run(readInput()); err != nil {
		log.Fatal(err.Error())
	}
}
