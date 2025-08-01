package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: grove <command> [args]")
		return
	}

	noRepo := flag.Bool("no-repo", false, "Create resource without a repository")
	noService := flag.Bool("no-service", false, "Create resource without a service")

	flag.Parse()

	command := os.Args[1]
	switch command {
	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Usage: grove create <resource_name> <field1:type1> [<field2:type2> ...]")
			return
		}

		// Check that the resource name is capitalized to be marked as public
		// This is required for the resource to be accessible due to Go's visibility rules
		if len(os.Args[3]) == 0 || os.Args[3][0] < 'A' || os.Args[3][0] > 'Z' {
			fmt.Println("Error: Resource name must start with a capital letter in order to be public.", os.Args[3])
			return
		}

		switch os.Args[2] {
		case "resource":
			if err := handleCreateResourceCommand(os.Args[3:], *noRepo, *noService); err != nil {
				fmt.Println("Error creating resource:", err.Error())
				return
			}
		case "controller":
			if err := handleCreateControllerCommand(os.Args[3:]); err != nil {
				fmt.Println("Error creating controller:", err.Error())
				return
			}
		default:
			fmt.Println("Unknown create command:", os.Args[2])
			fmt.Println("Available commands: resource, controller")
			return
		}
	case "init":
		handleInitCommand(os.Args[2:])
	case "help":
		writeHelpMenu()
	default:
		fmt.Println("Unknown command:", command)
		writeHelpMenu()
	}
}
