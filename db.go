package main

import (
	"encoding/json"
	"fmt"

	"github.com/HouzuoGuo/tiedot/dberr"
)

func setupDB() {
	// Reset DB
	// myDB.Drop("Commands")

	// Check if collection "Commands" exists
	if !myDB.ColExists(dbCollectionCommand) {
		// Create collection "Commands"
		if err := myDB.Create(dbCollectionCommand); err != nil {
			panic(err)
		}
	}

	// Use collection: Commands
	commands = myDB.Use(dbCollectionCommand)

	// Create indexing "title" for querying
	if !myDB.ColExists(dbCollectionCommand) {
		if err := commands.Index([]string{dbCollectionCommandFieldTitle}); err != nil {
			panic(err)
		}
	}
}

func readDB() {
	var cs []Command
	var c Command

	// Read documents
	commands.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		_ = json.Unmarshal(doc, &c)
		cs = append(cs, c)

		return true
	})

	webSocketWriteMessage(cs)
}

func writeDB(data map[string]interface{}) {
	if _, err := commands.Insert(data); err != nil {
		panic(err)
	}
}

func writeCommandToDB(title, description, command, output string) {
	writeDB(map[string]interface{}{
		dbCollectionCommandFieldTitle: title,
		"description":                 description,
		"command":                     command,
		"output":                      output,
	})
}

func testDB() {
	// Create two collections: Feeds and Votes
	if err := myDB.Create("Feeds"); err != nil {
		panic(err)
	}

	// if err := myDB.Create("Votes"); err != nil {
	// 	panic(err)
	// }

	// What collections do I now have?
	// for _, name := range myDB.AllCols() {
	// 	fmt.Printf("I have a collection called %s\n", name)
	// }

	// Rename collection "Votes" to "Points"
	// if err := myDB.Rename("Votes", "Points"); err != nil {
	// 	panic(err)
	// }

	// Drop (delete) collection "Points"
	// if err := myDB.Drop("Points"); err != nil {
	// 	panic(err)
	// }

	// Start using a collection (the reference is valid until DB schema changes or Scrub is carried out)
	feeds := myDB.Use("Feeds")

	// Insert document (afterwards the docID uniquely identifies the document and will never change)
	docID, err := feeds.Insert(map[string]interface{}{
		"name": "Go 1.2 is released",
		"url":  "golang.org"})
	if err != nil {
		panic(err)
	}

	// Read document
	readBack, err := feeds.Read(docID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document", docID, "is", readBack)

	// Update document
	err = feeds.Update(docID, map[string]interface{}{
		"name": "Go is very popular",
		"url":  "google.com"})
	if err != nil {
		panic(err)
	}

	// Process all documents (note that document order is undetermined)
	feeds.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		fmt.Println("For each document", id, "is", string(docContent))
		return true // move on to the next document OR
		// return false // do not move on to the next document
	})

	// More complicated error handing - identify the error Type.
	// In this example, the error code tells that the document no longer exists.
	if err := feeds.Delete(docID); dberr.Type(err) == dberr.ErrorNoDoc {
		fmt.Println("The document was already deleted")
	}

	// Drop (delete) collection "Feeds"
	if err := myDB.Drop("Feeds"); err != nil {
		panic(err)
	}
}
