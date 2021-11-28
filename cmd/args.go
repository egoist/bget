package main

import (
	"fmt"
	"os"
	"strings"
)

type AppArgs struct {
	Repo    string
	BinDir  string
	BinName string
}

func ParseArgs(args []string) (*AppArgs, error) {
	appArgs := AppArgs{}

	for index := 0; index < len(args); index++ {
		v := args[index]
		if v == "-v" || v == "--version" {
			printVersion()
			os.Exit(0)
		} else if v == "-h" || v == "--help" {
			printHelp()
			os.Exit(0)
		} else if v == "-d" || v == "--dir" {
			if index+1 >= len(args) {
				return nil, fmt.Errorf("--dir requires an argument")
			}
			appArgs.BinDir = args[index+1]
			index++
		} else if v == "-b" || v == "--bin" {
			if index+1 >= len(args) {
				return nil, fmt.Errorf("--bin requires an argument")
			}
			appArgs.BinName = args[index+1]
			index++
		} else if strings.HasPrefix(v, "-") {
			return nil, fmt.Errorf("unknown option: %s", v)
		} else {
			if appArgs.Repo == "" {
				appArgs.Repo = v
			} else {
				return nil, fmt.Errorf("too many arguments")
			}
		}
	}

	if appArgs.Repo == "" {
		printHelp()
		os.Exit(0)
	}

	return &appArgs, nil
}

const VERSION string = "1.0.4"

func printVersion() {
	println(VERSION)
}

func printHelp() {
	fmt.Printf(`
bget v%s	

Flags:
	-b, --bin <name>	 The name of the binary file to output (default: repo name)
	-d, --dir <dir>	 	 The directory to install the binary to (default: /usr/local/bin)
	-v, --version		 Print version
	-h, --help		     Print help

`, VERSION)
}
