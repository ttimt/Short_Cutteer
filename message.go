package main

import (
	"encoding/json"
	"sort"
)

const (
	messageKindCommand     = "command"
	messageOperationRead   = "read"
	messageOperationWrite  = "write"
	messageOperationDelete = "delete"
)

type Message struct {
	Kind      string      `json:"kind"`
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

// Process incoming message from web socket
func processIncomingMessage(m Message) {
	dataStr := m.Data.(string)

	if m.Kind == messageKindCommand {
		switch m.Operation {
		case messageOperationWrite:
			var c Command
			_ = json.Unmarshal([]byte(dataStr), &c)

			updateUserCommand(c)
			updateCommandLength(c.Command)
			sort.Ints(sliceCommandLen)
			writeCommandToDB(c.Title, c.Description, c.Command, c.Output)
		case messageOperationDelete:
			deleteCommandFromDB(dataStr)
		}
	}
}

// Create message struct to be written to web socket
func createMessage(kind string, operation string, jsonMsg interface{}) Message {
	msg := Message{
		Kind:      kind,
		Operation: operation,
		Data:      jsonMsg,
	}

	return msg
}
