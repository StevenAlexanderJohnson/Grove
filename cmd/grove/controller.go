package main

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/controller.go.tmpl
var controllerTemplate string

func createController(resourceName string) error {
	// Initialize the controller for Grove project management
	println("Controller initialized for Grove project management.")

	packageName, err := getModuleName()

	if err != nil {
		return err
	}

	templateData := map[string]interface{}{
		"ProjectName":  packageName,
		"ResourceName": resourceName,
	}

	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create("internal/controllers/" + resourceName + "Controller.go")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, templateData); err != nil {
		return err
	}

	return nil
}

func controllerHelp() {
	println("Controller command help:")
	println("Usage: create-controller <controller-path> <resource-model-path> <resource-service-path>")
	println("This command creates a new controller for Grove project management.")
}

func handleCreateControllerCommand(args []string) error {
	if len(args) < 3 || args[0] == "help" || args[0] == "--help" {
		controllerHelp()
		return nil
	}

	if err := handleCreateResourceCommand(args, false, false); err != nil {
		println("Error creating resource:", err.Error())
		return err
	}

	resourceName := args[0]

	if err := createController(resourceName); err != nil {
		println("Error creating controller:", err.Error())
		return err
	}

	println("Controller created successfully for", resourceName)
	return nil
}
