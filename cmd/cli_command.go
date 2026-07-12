package cmd

import "fmt"

type CliCommand struct {
	Name    string
	Args    []string
	Handler func(state *state, cmd *CliCommand) error
}

type Commands struct {
	RegisteredCommands map[string]func(*state, CliCommand) error
}

func (c *Commands) Register(commandName string, handler func(state *state, cmd CliCommand) error) {
	if c.RegisteredCommands == nil {
		c.RegisteredCommands = make(map[string]func(*state, CliCommand) error)
	}

	c.RegisteredCommands[commandName] = handler
}

func (c *Commands) RunCommand(state *state, cmd CliCommand) error {

	handler, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("Command '%s' not found", cmd.Name)
	}
	return handler(state, cmd)
}
