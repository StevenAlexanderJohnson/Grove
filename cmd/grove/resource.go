package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var (
	ErrResourceInvalidArguments      = fmt.Errorf("invalid arguments for create-resource command")
	ErrResourceNoFieldsProvided      = fmt.Errorf("no fields provided for resource")
	ErrResourceNoNameProvided        = fmt.Errorf("no resource name provided")
	ErrResourceModuleNotFound        = fmt.Errorf("go.mod file not found or not formatted correctly")
	ErrResourceCreationFailed        = fmt.Errorf("failed to create resource")
	ErrResourceRepoCreationFailed    = fmt.Errorf("failed to create repository for resource")
	ErrResourceServiceCreationFailed = fmt.Errorf("failed to create service for resource")
	ErrResourceIsPrivate             = fmt.Errorf("resource name must be public (start with an uppercase letter)")
)

//go:embed templates/model.go.tmpl
var modelTemplate string

//go:embed templates/service.go.tmpl
var serviceTemplate string

//go:embed templates/repo.go.tmpl
var repoTemplate string

func isPrivate(name string) bool {
	if name == "" {
		return false
	}

	return name[0] >= 'a' && name[0] <= 'z'
}

func capitalize(name string) string {
	if len(name) == 0 {
		return name
	}
	if name[0] >= 'a' && name[0] <= 'z' {
		return string(name[0]-32) + name[1:]
	}
	return name
}

func resourceHelp() {
	println("Resource command help:")
	println("Usage: create-resource <resource-name> [<resource-field-name>:<go type> ...]")
	println("The CodeGen will use the field names as provided, so lower case names will not be exported.")
	println("This command creates a new resource for Grove project management.")
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

	tmpl, err := template.New("model").Funcs(template.FuncMap{
		"isPrivate":  isPrivate,
		"capitalize": capitalize,
	}).Parse(modelTemplate)
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

func createRepo(resourceName string, moduleName string) error {
	if resourceName == "" {
		return ErrResourceNoNameProvided
	}
	if moduleName == "" {
		return ErrResourceModuleNotFound
	}
	println("Creating repository for resource:", resourceName)

	templateData := map[string]interface{}{
		"ModuleName":   moduleName,
		"ResourceName": resourceName,
	}

	tmpl, err := template.New("repo").Parse(repoTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse repo template: %w", err)
	}
	repoPath := "internal/repositories/" + resourceName + ".go"
	repoFile, err := os.Create(repoPath)
	if err != nil {
		return fmt.Errorf("failed to create repo file: %w", err)
	}
	defer repoFile.Close()

	if err := tmpl.Execute(repoFile, templateData); err != nil {
		return fmt.Errorf("failed to execute repo template: %w", err)
	}
	println("Repository created successfully at", repoPath)
	return nil
}

func createService(resourceName string, moduleName string) error {
	if resourceName == "" {
		return ErrResourceNoNameProvided
	}
	if moduleName == "" {
		return ErrResourceModuleNotFound
	}
	println("Creating service for resource:", resourceName)

	templateData := map[string]interface{}{
		"ModuleName":   moduleName,
		"ResourceName": resourceName,
	}

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}
	servicePath := "internal/services/" + resourceName + ".go"
	serviceFile, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer serviceFile.Close()

	if err := tmpl.Execute(serviceFile, templateData); err != nil {
		return fmt.Errorf("failed to execute service template: %w", err)
	}
	println("Service created successfully at", servicePath)
	return nil
}

func handleCreateResourceCommand(args []string, noRepo bool, noService bool) error {
	if len(args) < 2 {
		resourceHelp()
		return ErrResourceInvalidArguments
	}
	if (args[0] == "help" || args[0] == "--help") && len(args) == 1 {
		resourceHelp()
		return nil
	}
	// Check if the first argument is a valid resource name
	if args[0] == "" {
		println("Error: Resource name must be provided.")
		return ErrResourceNoNameProvided
	} else if isPrivate(args[0]) {
		println("Error: Resource name must be public (start with an uppercase letter).")
		return fmt.Errorf("%w: %v", ErrResourceIsPrivate, args[0])
	}

	moduleName, err := getModuleName()
	if err != nil {
		println("Error reading go.mod:", err.Error())
		return fmt.Errorf("%w: %v", ErrResourceModuleNotFound, err)
	}

	resourceName := args[0]

	if err := createModel(resourceName, parseResourceFields(args[1:])); err != nil {
		println("Error creating resource:", err.Error())
		return fmt.Errorf("%w: %v", ErrResourceCreationFailed, err)
	}

	fmt.Println(noRepo, noService)
	if !noRepo {
		if err := createRepo(resourceName, moduleName); err != nil {
			return fmt.Errorf("%w: %v", ErrResourceRepoCreationFailed, err)
		}
	}
	if !noService {
		if err := createService(resourceName, moduleName); err != nil {
			return fmt.Errorf("%w: %v", ErrResourceServiceCreationFailed, err)
		}
	}

	println("Resource created successfully:", resourceName)
	return nil
}
