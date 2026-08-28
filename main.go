package main

import (
	"fmt"
	"os"
)

const defaultBranch = "main"

func CreateDirs(paths []string) error {
	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateEmptyFiles(paths []string) error {
	for _, path := range paths {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func InitMainDir() error {
	// Initalizing default dirs
	paths := []string{
		".bit/objects/commits",
		".bit/objects/trees",
		".bit/objects/blobs",
		".bit/branches",
	}
	// Initalizing default files
	err := CreateDirs(paths)
	if err != nil {
		return err
	}
	files := []string{
		".bit/branches/" + defaultBranch,
		".bit/head",
		".bit/stage",
	}
	err = CreateEmptyFiles(files)
	if err != nil {
		return err
	}
	// Head file should contain the main branch at the beggining
	err = os.WriteFile(".bit/head", []byte(defaultBranch), 0644)
	err = os.WriteFile(".bit/stage", []byte("{}"), 0644)
	return err
}

func RunCommand(command string, args []string) {
	switch command {
	case "init":
		err := InitMainDir()
		if err != nil {
			fmt.Println("Error while tring to init project: " + err.Error())
		} else {
			fmt.Println("Project initalized successfully")
		}
	case "stage":
		err := StageFiles(args)
		if err != nil {
			fmt.Println("Error while tring to stage file/s: " + err.Error())
		}
	default:
		fmt.Println("Unknown command: " + command)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bit <command> <args>")
		return
	}
	command := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	RunCommand(command, args)
}
