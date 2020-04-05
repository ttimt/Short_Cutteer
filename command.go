package main

import (
	"sort"
)

var (
	userCommands     = make(map[string]*Command)
	maxCommandLen    int
	uniqueCommandLen = make(map[int]struct{})
	sliceCommandLen  []int
)

type Command struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Output      string `json:"output"`
}

func updateUserCommand(c Command) {
	userCommands[c.Command] = &c

	if len(c.Command) > maxCommandLen {
		maxCommandLen = len(c.Command)
	}
}

func updateAllCommandLength() {
	commands := getAllCommands()

	for _, v := range commands {
		updateCommandLength(v.Command)
	}

	sort.Ints(sliceCommandLen)
}

func updateCommandLength(command string) {
	commandLength := len(command)

	if _, ok := uniqueCommandLen[commandLength]; !ok {
		sliceCommandLen = append(sliceCommandLen, commandLength)
		uniqueCommandLen[commandLength] = struct{}{}
	}
}

func getAllCommands() []Command {
	commands := make([]Command, len(userCommands))

	i := 0
	for k := range userCommands {
		commands[i] = *userCommands[k]
		i++
	}

	return commands
}
