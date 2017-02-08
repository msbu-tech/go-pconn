package main

import (
	"flag"
	"fmt"

	"github.com/msbu-tech/go-pconn/cmd/version"
)

var (
	helpPtr    *bool = flag.Bool("help", false, "Help")
	versionPtr *bool = flag.Bool("version", false, "Version Info")
)

func main() {

	flag.Usage = usage
	flag.Parse()

	if *helpPtr == true {
		usage()

		return
	} else if *versionPtr == true {
		fmt.Println("go-pconn", version.Version)
		fmt.Println("Copyright (c) MSBU-Tech, 2017")

		return
	} else {
		usage()

		return
	}
}

func usage() {
	fmt.Println("go-pconn")
	fmt.Println("")
	fmt.Println("usage: go-sf [commands] [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("\t-help         # Show usage")
	fmt.Println("\t-version      # Show version")
	fmt.Println("")
}
