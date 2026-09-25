package commands

import (
	"fmt"

	"github.com/noamkalmar/bit/internal/files"
	"github.com/noamkalmar/bit/internal/utils"
)

func RunCommand(command string, args []string) {
	if command != "init" && !files.IsProjectInitialized() {
		fmt.Println("Error: please run 'bit init' to initialize the project before running other commands")
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
		err := StageFiles(utils.SliceToPathSet(args))
		if err != nil {
			fmt.Println("Error while tring to stage file/s: " + err.Error())
		}
	case "unstage":
		err := UnstageFiles(utils.SliceToPathSet(args))
		if err != nil {
			fmt.Println("Error while tring to unstage file/s: " + err.Error())
		}
	case "commit":
		err := CreateCommit(args)
		if err != nil {
			fmt.Println("Error while tring to create a new commit: " + err.Error())
		}
	case "status":
		err := PrintStatus()
		if err != nil {
			fmt.Println("Error while trying to display status: " + err.Error())
		}
	default:
		fmt.Println("Unknown command: " + command)
	}
}
