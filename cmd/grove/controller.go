package main

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/controller.go.tmpl
var controllerTemplate string

func createController(controllerPath, resourceModelPath, resourceServicePath string) error {
	// Initialize the controller for Grove project management
	println("Controller initialized for Grove project management.")
	templateData := map[string]interface{}{
		"Package":             "controllers",
		"ResourceModelPath":   resourceModelPath,
		"ResourceServicePath": resourceServicePath,
		"ResourceName":        "TestResource",
	}

	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(controllerPath)
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

func handleCreateControllerCommand(args []string) {
	if len(args) < 3 {
		controllerHelp()
		return
	}

	controllerPath := args[0]
	resourceModelPath := args[1]
	resourceServicePath := args[2]

	if err := createController(controllerPath, resourceModelPath, resourceServicePath); err != nil {
		println("Error creating controller:", err.Error())
		return
	}

	println("Controller created successfully at", controllerPath)
}
