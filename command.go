package main

import (
	"fmt"
)

var (
	userCommands     = make(map[string]*Command)
	maxCommandLen    int
	maxCommandLength *CommandLength
	minCommandLength *CommandLength
)

type Command struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Command       string `json:"command"`
	Output        string `json:"output"`
	commandLength *CommandLength
}

type CommandLength struct {
	length  int
	next    *CommandLength
	command *Command // Parent
}

func updateUserCommand(c Command) {
	userCommands[c.Command] = &c

	if len(c.Command) > maxCommandLen {
		maxCommandLen = len(c.Command)
	}
}

func updateCommandLength() {
	commands := getAllCommands()

	// Process
	fmt.Println(commands)
}

func getAllCommands() *[]Command {
	commands := make([]Command, len(userCommands))

	i := 0
	for k := range userCommands {
		commands[i] = *userCommands[k]
		i++
	}

	return &commands
}
