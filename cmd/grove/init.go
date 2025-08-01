package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

var (
	ErrFolderAlreadyExists = errors.New("the project folder already exists")
	ErrGitNotInstalled     = errors.New("git is not installed, please install git to use this command")
	ErrInvalidTemplate     = errors.New("invalid project template specified")
)

//go:embed templates/main.go.tmpl
var mainTemplate string

func initializeFolder(projectName string) error {
	stat, err := os.Stat(projectName)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(projectName, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if stat != nil && !stat.IsDir() {
		return ErrFolderAlreadyExists
	}

	if err := os.Mkdir(projectName+"/cmd", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(projectName+"/internal", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(projectName+"/internal/controllers", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(projectName+"/internal/models", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(projectName+"/internal/services", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(projectName+"/internal/repositories", 0755); err != nil {
		return err
	}

	if err := os.Chdir(projectName); err != nil {
		return err
	}
	return nil
}

func initializeProject(projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse main template: %w", err)
	}

	mainFile, err := os.Create("./cmd/main.go")
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	defer mainFile.Close()

	if err := tmpl.Execute(
		mainFile,
		map[string]string{
			"ProjectName": projectName,
		},
	); err != nil {
		return fmt.Errorf("failed to execute main template: %w", err)
	}

	if err := exec.Command("go", "get", "github.com/StevenAlexanderJohnson/grove@v0").Run(); err != nil {
		return fmt.Errorf("failed to tidy go module: %w", err)
	}
	return nil
}

func handleInitCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: init <project-name>")
		return
	}

	projectName := args[0]

	if err := initializeFolder(projectName); err != nil {
		fmt.Println("Error initializing folder:", err.Error())
		return
	}

	if err := initializeProject(projectName); err != nil {
		fmt.Println("Error initializing project:", err.Error())
		return
	}

	fmt.Println("Grove project initialized successfully:", projectName)
}
