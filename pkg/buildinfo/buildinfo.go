package buildinfo

import (
	"flag"
	"fmt"
	"os"
)

var version = flag.Bool("version", false, "show version")

// Commit must be set via -ldflags '-X'
var Commit string

// Init must be called after flag.Parse call.
func Init() {
	if *version {
		printVersion()
		os.Exit(0)
	}
}

func init() {
	oldUsage := flag.Usage
	flag.Usage = func() {
		printVersion()
		fmt.Println()
		oldUsage()
	}
}

func printVersion() {
	fmt.Fprintf(flag.CommandLine.Output(), "Commit: %s\n", Commit)
}
