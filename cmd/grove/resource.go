package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/model.go.tmpl
var modelTemplate string

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func isPrivate(name string) bool {
	if name == "" {
		return false
	}

	return name[0] >= 'a' && name[0] <= 'z'
}

func resourceHelp() {
	println("Resource command help:")
	println("Usage: create-resource <resource-name> [<resource-field-name>:<go type> ...]")
	println("The CodeGen will use the field names as provided, so lower case names will not be exported.")
	println("This command creates a new resource for Grove project management.")
}

func createRepo(resourceName string, moduleName string) error {
	return fmt.Errorf("not implemented")
}

func parseResourceFields(args []string) map[string]string {
	fields := make(map[string]string)
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			println("Invalid field format:", arg)
			continue
		}
		fields[parts[0]] = parts[1]
	}
	return fields
}

// Gets the module name from go.mod
// Returns an error if go.mod does not exist or is not formatted correctly.
func getModuleName() (string, error) {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		println("Error: This command must be run in a Grove project directory with a go.mod file.")
		return "", err
	}

	modFile, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer modFile.Close()

	scanner := bufio.NewScanner(modFile)
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	line := scanner.Text()
	if !strings.HasPrefix(line, "module ") {
		return "", nil
	}

	moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
	if moduleName == "" {
		return "", nil
	}

	return moduleName, nil
}

// Creates the model file for the resource.
func createModel(resourceName string, fields map[string]string) error {
	resourcePath := "internal/models/" + resourceName + ".go"
	templateData := map[string]interface{}{
		"ResourceName":   resourceName,
		"ResourceFields": fields,
	}

	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse model template: %w", err)
	}
	modelFile, err := os.Create(resourcePath)
	if err != nil {
		return fmt.Errorf("failed to create model file: %w", err)
	}
	defer modelFile.Close()

	if err := tmpl.Execute(modelFile, templateData); err != nil {
		return fmt.Errorf("failed to execute model template: %w", err)
	}
	println("Model created successfully at", resourcePath)

	return nil
}

func handleCreateResourceCommand(args []string, noRepo bool) {
	if len(args) < 2 {
		resourceHelp()
		return
	}
	if (args[0] == "help" || args[0] == "--help") && len(args) == 1 {
		resourceHelp()
		return
	}

	moduleName, err := getModuleName()
	if err != nil {
		println("Error reading go.mod:", err.Error())
		return
	}

	resourceName := args[0]

	if err := createModel(resourceName, parseResourceFields(args[1:])); err != nil {
		println("Error creating resource:", err.Error())
		return
	}

	if !noRepo {
		if err := createRepo(resourceName, moduleName); err != nil {

		}
	}
}
