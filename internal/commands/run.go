package commands

import (
	"fmt"

	"github.com/noamkalmar/bit/internal/files"
)

func RunCommand(command string, args []string) {
	switch command {
	case "init":
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
