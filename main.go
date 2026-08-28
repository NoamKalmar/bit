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

func InitMainDir() error {
	paths := []string{
		".bit/objects/commits",
		".bit/objects/trees",
		".bit/objects/blobs",
		".bit/branches",
	}
	err := CreateDirs(paths)
	if err != nil {
		return err
	}
	headFile, err := os.Create(".bit/head")
	if err != nil {
		return err
	}
	defer headFile.Close()
	_, err = headFile.WriteString(defaultBranch)
	if err != nil {
		return err
	}
	stageFile, err := os.Create(".bit/stage")
	if err != nil {
		return err
	}
	stageFile.Close()
	return nil
}

func main() {
	command := os.Args[1]
	switch command {
	case "init":
		err := InitMainDir()
		if err != nil {
			fmt.Println("Error while tring to init project: " + err.Error())
		}
		return
	default:
		fmt.Println("Unknown command: " + command)
		return
	}
}
