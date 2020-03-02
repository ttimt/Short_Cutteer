package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	httpPort          = 8081
	htmlFilePath      = "html/"
	jqueryFilePath    = "node_modules/jquery/dist/jquery.min.js"
	jqueryUIFilePath  = "node_modules/jquery-ui-dist/jquery-ui.min.js"
	semanticFilePath  = "node_modules/fomantic-ui/dist/"
	mainHtmlFile      = "index.html"
	templateFilesPath = "html/template/"
)

var (
	httpPortStr  = ":" + strconv.Itoa(httpPort)
	httpURL      = "http://localhost" + httpPortStr
	wsUpgrader   = websocket.Upgrader{}
	wsConnection webSocketConnection
	tplt         *template.Template
)

type webSocketConnection struct {
	mux    sync.Mutex
	client *websocket.Conn
}

type httpFileSystem struct {
	fileSystem http.FileSystem
}

// Initialize the templates that will be used
func initializeTemplates() {
	// Get all template files info
	templateFiles, err := ioutil.ReadDir(templateFilesPath)
	if err != nil {
		panic(err)
	}

	// Get the name of template files
	templateFilesName := make([]string, len(templateFiles)+1)

	templateFilesName[0] = htmlFilePath + mainHtmlFile
	for k := range templateFiles {
		templateFilesName[k+1] = templateFilesPath + templateFiles[k].Name()
	}

	// Parse template files
	tplt, err = template.ParseFiles(templateFilesName...)
	if err != nil {
		panic(err)
	}
}

// Setup and run HTTP server
func setupHTTPServer() {
	// Setup mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Server template file
		err := tplt.ExecuteTemplate(w, mainHtmlFile, nil)
		if err != nil {
			panic(err)
		}
	})

	mux.HandleFunc("/dist/jquery.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, jqueryFilePath)
	})

	mux.HandleFunc("/dist/jquery-ui.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, jqueryUIFilePath)
	})

	mux.Handle("/html/", http.StripPrefix("/html", http.FileServer(httpFileSystem{http.Dir(htmlFilePath)})))
	mux.Handle("/dist/", http.StripPrefix("/dist", http.FileServer(httpFileSystem{http.Dir(semanticFilePath)})))

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "icons/icon.ico")
	})

	mux.HandleFunc("/ws", handleWebSocket)

	// Concurrently run web server
	go func() {
		log.Println("Started listening on", httpPort)
		log.Fatal(http.ListenAndServe(httpPortStr, mux))
	}()
}

// Handle web socket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade GET request to a web socket
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Assign the web socket client
	wsConnection.client = wsConn

	// Send collectionCommands to UI
	webSocketWriteMessage(createMessage(messageKindCommand, messageOperationWrite, getAllCommands()))

	// Start reading message
	webSocketReadMessage()
}

// Read message from web socket
func webSocketReadMessage() {
	for {
		// Read message
		var m Message
		err := wsConnection.client.ReadJSON(&m)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Println("Web socket connection closed:", err)
			} else {
				log.Println("Web socket error:", err)
			}

			break
		}

		// Read success
		processIncomingMessage(m)
	}

	// Close the connection at the end if read fails
	_ = wsConnection.client.Close()
}

// Write message to web socket
func webSocketWriteMessage(msg Message) {
	wsConnection.mux.Lock()

	// Check any client exist
	if wsConnection.client == nil {
		wsConnection.mux.Unlock()
		log.Println("Unable to write message: no client exist!")
		return
	}

	err := wsConnection.client.WriteJSON(msg)
	if err != nil {
		log.Println("Write connection error:", err)
	} else {
		log.Println("Write message succeeded:", time.Now().Format(time.Kitchen))
	}

	wsConnection.mux.Unlock()
}

// Override Open method to hide file directory and serve index.html if it exist
func (fs httpFileSystem) Open(name string) (http.File, error) {
	file, err := fs.fileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		index := strings.TrimSuffix(name, "/") + "/index.html"

		if _, err := fs.fileSystem.Open(index); err != nil {
			return nil, err
		}
	}

	return file, nil
}
