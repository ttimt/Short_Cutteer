package main

var (
	userCommands = make(map[string]*Command)
	maxBufferLen int
	bufferStr    string
)

type Command struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Output      string `json:"output"`
}

func updateUserCommand(c Command) {
	userCommands[c.Command] = &c

	if len(c.Command) > maxBufferLen {
		maxBufferLen = len(c.Command)
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
