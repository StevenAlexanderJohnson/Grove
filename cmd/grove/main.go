package main

import (
	"flag"
	"os"
)

var allowedTemplates map[string]bool = map[string]bool{
	"skeleton": true,
}

func main() {
	projectName := flag.String("project", "", "Name of the Grove project")
	projectTemplate := flag.String("template", "", "Template to use for the Grove project")
	flag.Parse()

	if !allowedTemplates[*projectTemplate] {
		println("Error: Invalid template specified. Allowed templates are:")
		for template := range allowedTemplates {
			println("- " + template)
		}
		return
	}

	if *projectName == "" {
		println("Error: Project name must be specified.")
		return
	}

	println("Creating Grove project:", *projectName, "with template:", *projectTemplate)
	if err := initializeFolder(*projectName); err != nil {
		println("Error initializing folder:", err.Error())
		return
	}
	err := os.Chdir(*projectName)
	if err != nil {
		println("Error changing directory:", err.Error())
		return
	}

	if err := initializeProject(*projectTemplate); err != nil {
		println("Error initializing project:", err.Error())
		return
	}
}
