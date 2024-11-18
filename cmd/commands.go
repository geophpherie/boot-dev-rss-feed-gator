package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*State, command) error
}

func (c *commands) register(name string, f func(*State, command) error) {
	/* Registers a new handler function for a given command name*/
	c.commands[name] = f
}

func (c *commands) run(s *State, cmd command) error {
	/* Runs a given command provided it exists*/
	f, ok := c.commands[cmd.name]
	if !ok {
		return errors.New("Function not registered")
	}
	return f(s, cmd)
}
