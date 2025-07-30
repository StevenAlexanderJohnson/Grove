package main

import (
	"errors"
	"os"
	"os/exec"
)

var (
	ErrFolderAlreadyExists = errors.New("the project folder already exists")
	ErrGitNotInstalled     = errors.New("git is not installed, please install git to use this command")
	ErrInvalidTemplate     = errors.New("invalid project template specified")
)

func initializeFolder(projectName string) error {
	stat, err := os.Stat(projectName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(projectName, 0755)
		} else {
			return err
		}
	}
	if stat != nil && !stat.IsDir() {
		return ErrFolderAlreadyExists
	}
	return err
}

func initializeProject(projectTemplate string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return ErrGitNotInstalled
	}

	// Clone into current directory
	cmd := exec.Command("git", "clone", "--filter=blob:none", "--no-checkout", "https://github.com/StevenAlexanderJohnson/grove.git", ".")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Enable sparse-checkout
	cmd = exec.Command("git", "sparse-checkout", "init", "--cone")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set sparse-checkout to the template folder
	templatePath := "v1/cmd/template/skeleton"
	switch projectTemplate {
	case "skeleton":
		// already set
	default:
		return ErrInvalidTemplate
	}

	cmd = exec.Command("git", "sparse-checkout", "set", templatePath)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "checkout")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Move files from templatePath to current directory
	cmd = exec.Command("bash", "-c", "shopt -s dotglob && mv "+templatePath+"/* . && rmdir "+templatePath)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Remove .git directory to clean up
	cmd = exec.Command("rm", "-rf", ".git")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
