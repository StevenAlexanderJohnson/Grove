package main

import (
	"flag"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: grove <command> [args]")
		return
	}

	noRepo := flag.Bool("no-repo", false, "Create resource without a repository")

	flag.Parse()

	command := os.Args[1]
	switch command {
	case "create-resource":
		handleCreateResourceCommand(os.Args[2:], *noRepo)
	case "create-controller":
		handleCreateControllerCommand(os.Args[2:])
	case "init":
		handleInitCommand(os.Args[2:])
	case "help":
		writeHelpMenu()
	default:
		println("Unknown command:", command)
		writeHelpMenu()
	}
}
