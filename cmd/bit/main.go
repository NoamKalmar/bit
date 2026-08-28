package main

import (
	"fmt"
	"os"

	"github.com/noamkalmar/bit/internal/commands"
)

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
	commands.RunCommand(command, args)
}
