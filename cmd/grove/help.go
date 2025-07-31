package main

import "strings"

func writeHelpMenu() {
	var builder strings.Builder
	builder.WriteString("Grove CLI Help Menu:\n")
	builder.WriteString("Available commands:\n")
	builder.WriteString("  create <controller|resource> <name> [<field_name:go_type> ...] - Generates a new controller or resource with the fields specified.\n")
	builder.WriteString("  init <project-name> - Initialize a new Grove project\n")
	builder.WriteString("    - <project-name> is the name of the project and will be used as the go mod name if go mod doesn't already exists.\n")
	builder.WriteString("  help - Show this help menu\n")
	println(builder.String())
}
