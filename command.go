package main

var (
	userCommands  = make(map[string]*Command)
	maxCommandLen int
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

func getAllCommands() *[]Command {
	commands := make([]Command, len(userCommands))

	i := 0
	for k := range userCommands {
		commands[i] = *userCommands[k]
		i++
	}

	return &commands
}
