package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
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

func GetHash(content []byte, hasher hash.Hash) string {
	hasher.Write(content)
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func StageFile(path string, hasher hash.Hash) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	hash := GetHash(content, hasher)
	file, err := os.Create(".bit/objects/blobs/" + hash)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(string(content))
	return nil
}

func StageFiles(paths []string) error {
	hasher := sha1.New()
	for _, path := range paths {
		err := StageFile(path, hasher)
		if err != nil {
			return err
		}
	}
	return nil
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
	command := os.Args[1]
	args := os.Args[2:]
	RunCommand(command, args)
}
