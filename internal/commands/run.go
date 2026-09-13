package commands

import (
	"fmt"

	"github.com/noamkalmar/bit/internal/files"
)

func RunCommand(command string, args []string) {
	if command != "init" && !files.IsProjectInitialized() {
		fmt.Println("Error: please run 'bit init' to initalize the project before running other commands")
		return
	}
	switch command {
	case "init":
		if files.IsProjectInitialized() {
			fmt.Println("Error: project is already initialized")
			return
		}
		err := files.InitMainDir()
		if err != nil {
			fmt.Println("Error while tring to init project: " + err.Error())
		} else {
			fmt.Println("Project initialized successfully")
		}
	case "stage":
		err := StageFiles(args)
		if err != nil {
			fmt.Println("Error while tring to stage file/s: " + err.Error())
		}
	case "unstage":
		err := UnstageFiles(args)
		if err != nil {
			fmt.Println("Error while tring to unstage file/s: " + err.Error())
		}
	case "commit":
		err := CreateCommit(args)
		if err != nil {
			fmt.Println("Error while tring to create a new commit: " + err.Error())
		}
	default:
		fmt.Println("Unknown command: " + command)
	}
}
