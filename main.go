package main

import (
	"fmt"
	"os"

	"github.com/thalesraymond/gator/cmd"

	_ "github.com/lib/pq"
)

func main() {
	state, err := cmd.NewState()

	if err != nil {
		fmt.Println(err)
		return
	}

	commands := cmd.Commands{}
	commands.Register("login", cmd.HandleLogin)
	commands.Register("register", cmd.HandleRegister)
	commands.Register("reset", cmd.HandleReset)
	commands.Register("users", cmd.HandleUsers)
	commands.Register("agg", cmd.HandleAggregator)
	commands.Register("addfeed", cmd.HandleAddFeed)
	commands.Register("feeds", cmd.HandlerListFeeds)
	commands.Register("follow", cmd.HandleFollow)
	commands.Register("following", cmd.HandleFollowing)
	

	if len(os.Args) < 2 {
		return
	}

	cliCmd := cmd.CliCommand{
		Name: os.Args[1],
		Args: os.Args[1:],
	}

	err = commands.RunCommand(state, cliCmd)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
