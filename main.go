package main

import (
	"fmt"
	"os"

	"github.com/thalesraymond/gator/cmd"
)

func main() {
	state, err := cmd.NewState()

	if err != nil {
		fmt.Println(err)
		return
	}

	commands := cmd.Commands{}
	commands.Register("login", cmd.HandleLogin)

	cliCmd := cmd.CliCommand{
		Name: "login",
		Args: os.Args[1:],
	}

	error := commands.RunCommand(state, cliCmd)

	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
}
