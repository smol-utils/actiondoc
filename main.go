package main

import (
	"fmt"
	"os"

	"github.com/smol-utils/actiondoc/cmd"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "generate", "gen":
		err = cmd.Generate(os.Args[2:])
	case "version":
		fmt.Println("actiondoc", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `actiondoc - generate documentation for GitHub Actions workflows

Usage:
  actiondoc <command> [arguments]

Commands:
  generate    Generate Markdown docs from workflow files
  version     Print version
  help        Show this help

Run "actiondoc generate --help" for subcommand details.
`)
}
