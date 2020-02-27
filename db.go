package main

import (
	"encoding/json"
	"fmt"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/HouzuoGuo/tiedot/dberr"
)

const (
	dbCollectionCommand             = "Commands"
	dbCollectionCommandFieldTitle   = "title"
	dbCollectionCommandFieldCommand = "command"
)

var (
	myDB               *db.DB
	collectionCommands *db.Col
)

// Initialize collections and indexes
func setupDB() {
	// Track if command collection exists
	existCommandCollection := myDB.ColExists(dbCollectionCommand)

	// Create collection 'Commands' if does not exists
	if !existCommandCollection {
		// Create collection "Commands"
		if err := myDB.Create(dbCollectionCommand); err != nil {
			panic(err)
		}
	}

	// Select collection 'Commands' for usage
	collectionCommands = myDB.Use(dbCollectionCommand)

	// Create index "title" for querying in 'Commands'
	if !existCommandCollection {
		if err := collectionCommands.Index([]string{dbCollectionCommandFieldTitle}); err != nil {
			panic(err)
		}
	}

	// Read from DB and import to in memory struct
	readAndImportFromDB()
}

// Reset collection by the given name in DB
func resetDBCollection(nameCollection string) {
	// Drop the collection
	_ = myDB.Drop(nameCollection)
}

// Read all Commands from DB and import to in memory struct
func readAndImportFromDB() {
	// Temporary store unmarshalled Command
	var c Command

	// Read all documents
	collectionCommands.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		// Convert Command from stored JSON format
		_ = json.Unmarshal(doc, &c)

		// Add newly read Command to in memory userCommands struct
		updateUserCommand(c)

		// Continue to the next item in the collection
		return true
	})
}

// Write data to the specified collection
func writeToDB(collection *db.Col, data map[string]interface{}) {
	if _, err := collection.Insert(data); err != nil {
		panic(err)
	}
}

// Write a command to the Command collection
func writeCommandToDB(title, description, command, output string) {
	writeToDB(collectionCommands, map[string]interface{}{
		dbCollectionCommandFieldTitle: title,
		"description":                 description,
		"command":                     command,
		"output":                      output,
	})
}

// Delete a command from Command collection
func deleteCommandFromDB(title string) {
	// Get the Command with the title
	queryResult := retrieveCommandsFromTitle(title)

	for id := range queryResult {
		// Get the commandStr(string) of the Command to use for deleting the in memory commandStr in userCommands struct
		commandToBeDeleted, _ := collectionCommands.Read(id)
		commandStr := commandToBeDeleted[dbCollectionCommandFieldCommand].(string)

		// Delete command from collection and check for error
		if err := collectionCommands.Delete(id); dberr.Type(err) == dberr.ErrorNoDoc {
			fmt.Println("The document was already deleted")
		} else if err != nil {
			panic(err)
		} else {
			// Delete from in memory struct userCommands
			delete(userCommands, commandStr)
		}
	}
}

// Return Commands by querying with the specified title
func retrieveCommandsFromTitle(title string) map[int]struct{} {
	// Store the query result
	queryResult := make(map[int]struct{})

	// Retrieve the query string
	query := query(title, dbCollectionCommandFieldTitle)

	// Execute the query
	if err := db.EvalQuery(query, collectionCommands, &queryResult); err != nil {
		panic(err)
	}

	return queryResult
}

// Create a query to find a specific value from the specified index
// Index must first be set in setupDB()
func query(value, index string) (query interface{}) {
	// Create the query with the index
	_ = json.Unmarshal([]byte(`[{"eq": "`+value+`", "in": ["`+index+`"]}]`), &query)

	return
}
