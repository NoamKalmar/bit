package commands

import (
	"fmt"

	"github.com/noamkalmar/bit/internal/files"
)

func RunCommand(command string, args []string) {
	if command != "init" && !files.IsProjectInitalized() {
		fmt.Println("Error: please run 'bit init' to initalize the project before running other commands")
		return
	}
	switch command {
	case "init":
		if files.IsProjectInitalized() {
			fmt.Println("Error: project is already initalized")
			return
		}
		err := files.InitMainDir()
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
