// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/gorilla/websocket"

	"github.com/ttimt/systray"

	_ "github.com/HouzuoGuo/tiedot/db"
	_ "github.com/lxn/walk"

	. "github.com/ttimt/Short_Cutteer/hook/windows"
	icon "github.com/ttimt/Short_Cutteer/icons"
)

const (
	htmlFilePath     = "html/"
	jqueryFilePath   = "node_modules/jquery/dist/jquery.min.js"
	jqueryUIFilePath = "node_modules/jquery-ui-dist/jquery-ui.min.js"
	semanticFilePath = "node_modules/fomantic-ui/dist/"

	mainHtmlFile      = "index.html"
	templateFilesPath = "html/template/"

	dbPath = "db"

	httpPort = 8080

	dbCollectionCommand           = "Commands"
	dbCollectionCommandFieldTitle = "title"

	messageKindCommand     = "command"
	messageOperationRead   = "read"
	messageOperationWrite  = "write"
	messageOperationDelete = "delete"
)

var (
	hhook                  HHOOK
	currentKeyStrokeSignal = make(chan rune)
	userCommands           = make(map[string]*Command)
	maxBufferLen           int
	bufferStr              string

	httpPortStr = ":" + strconv.Itoa(httpPort)
	httpURL     = "http://localhost" + httpPortStr

	processInterruptSignal = make(chan os.Signal)

	wsUpgrader   = websocket.Upgrader{}
	wsConnection webSocketConnection

	myDB               *db.DB
	commandsCollection *db.Col

	t *template.Template
)

type webSocketConnection struct {
	mux    sync.Mutex
	client *websocket.Conn
}

type httpFileSystem struct {
	fileSystem http.FileSystem
}

type Command struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Output      string `json:"output"`
}

type Message struct {
	Kind      string      `json:"kind"`
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

func receiveHook() {
	// Declare a keyboard hook callback function (type HOOKPROC)
	hookCallback := func(code int, wParam WPARAM, lParam LPARAM) LRESULT {
		// If keystroke is pressed down
		if wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN {
			// Retrieve the keyboard hook struct
			keyboardHookData := (*TagKBDLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))

			// Retrieve current keystroke from keyboard hook struct's vkCode
			currentKeystroke := rune((*keyboardHookData).VkCode)

			// Send the keystroke to be processed
			select {
			case currentKeyStrokeSignal <- currentKeystroke:

			default:
				// Skip if the channel currentKeyStrokeSignal is busy,
				// as it means the current keystroke is sent by the processHook while processing,
				// we want to ignore keystrokes that we sent ourself
			}
		}

		// Return CallNextHookEx result to allow keystroke to be displayed on user screen
		return CallNextHookEx(0, code, wParam, lParam)
	}

	// Install a Windows hook that listen to keyboard input
	hhook, _ = SetWindowsHookExW(WH_KEYBOARD_LL, hookCallback, 0, 0)
	if hhook == 0 {
		panic("Failed to set Windows hook")
	}

	// Run hook processing goroutine
	go processHook()

	// Start retrieving message from the hook
	if b, _ := GetMessageW(0, 0, 0, 0); !b {
		panic("Failed to get message")
	}
}

// Process your received keystroke here
func processHook() {
	for {
		// Receive keystroke as rune from channel
		currentKeyStroke := <-currentKeyStrokeSignal

		// Process keystroke
		fmt.Printf("Current key: %d 0x0%X %c\n", currentKeyStroke, currentKeyStroke, currentKeyStroke)

		shiftKeyState, _ := GetKeyState(VK_SHIFT)
		capsLockState, _ := GetKeyState(VK_CAPITAL)

		_, char, _ := findAllKeyCode(uint16(currentKeyStroke), 0, getKeyStateBool(shiftKeyState), getKeyStateBool(capsLockState, true))

		// Reset if character is not a letter/symbol/number
		if char == -1 {
			bufferStr = ""

			continue
		}

		// User pre-commandsCollection can be CTRL, ALT, SHIFT and 1 letter afterwards
		// User can create shortcut MODIFIER + a key or text commandsCollection + tab or space (remove commandsCollection)
		switch char {
		case '\b':
			if len(bufferStr) > 0 {
				bufferStr = bufferStr[:len(bufferStr)-1]
			}
		case '\r':
			bufferStr += windowsNewLine
			bufferStr = ""
		case '\t':
			if str, ok := userCommands[bufferStr]; ok {
				// Send input
				tagInputs := createTagInputs(str.Output)
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
				bufferStr = ""
			}
		case ' ':
			if str, ok := userCommands[bufferStr]; ok {
				// Send input
				tagInputsBackspace := multiplyTagInputKey(tagInputBackspaceDown(), len(bufferStr)+1)
				_, _ = SendInput(uint(len(tagInputsBackspace)), (*LPINPUT)(&tagInputsBackspace[0]), int(unsafe.Sizeof(tagInputsBackspace[0])))

				tagInputs := createTagInputs(str.Output)
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))

				bufferStr = ""
			}
		default:
			// If buffer full, trim
			if len(bufferStr) >= maxBufferLen {
				bufferStr = bufferStr[1:]
			}

			bufferStr += string(char)

			//  TODO CHECK SHORTCUT KET EXIST EX: CTRL + ALT + F
		}

		fmt.Println("Buffer string:", bufferStr)

		// If left bracket, left bracket, space x2, right bracket, left arrow x2
		// If first double/single quotes, right quote, left arrow
		// If after ( or [ (bracket already complete) and enter, enter and left arrow x2? depends on default enter behavior
		// if after {  } and enter, delete, backspace, left arrow, enter, right arrow, enter x2, left arrow
		// if command xxx and tab, enter yyy
		// if command xxx and space, backspace xxx len and enter yyy
		// If brackets or quotes, just can copy text (mouse + keyboard to detect there is text selected),
		//      copy, left bracket/quote, space, paste, space, right bracket/quote
		// If 'shortcut key', copy and format and paste
	}
}

func updateUserCommand(c Command) {
	userCommands[c.Command] = &c

	if len(c.Command) > maxBufferLen {
		maxBufferLen = len(c.Command)
	}
}

func init() {
	// Setup process interrupt signal
	signal.Notify(processInterruptSignal, os.Interrupt)

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
	t, err = template.ParseFiles(templateFilesName...)
	if err != nil {
		panic(err)
	}

	// Load DB
	// Create or open the database
	myDB, err = db.OpenDB(dbPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Call systray GUI
	systray.Run(onReady, nil)
}

func onReady() {
	// Run the server
	setupHTTPServer()

	// Start low level keyboard listener
	setupWindowsHook()

	// Setup system tray icon
	setupTrayIcon()

	// Setup DB
	setupDB()
}

func setupHTTPServer() {
	// Setup mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Server template file
		err := t.ExecuteTemplate(w, mainHtmlFile, nil)
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

	// Send commandsCollection to UI
	webSocketWriteMessage(messageKindCommand, messageOperationWrite, getAllCommands())

	// Start reading message
	webSocketReadMessage()
}

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

func getAllCommands() *[]Command {
	commands := make([]Command, len(userCommands))

	i := 0
	for k := range userCommands {
		commands[i] = *userCommands[k]
		i++
	}

	return &commands
}

func processIncomingMessage(m Message) {
	dataStr := m.Data.(string)

	if m.Kind == messageKindCommand {
		switch m.Operation {
		case messageOperationWrite:
			var c Command
			_ = json.Unmarshal([]byte(dataStr), &c)

			updateUserCommand(c)
			writeCommandToDB(c.Title, c.Description, c.Command, c.Output)
		case messageOperationDelete:
			deleteCommandFromDB(dataStr)
		}
	}
}

func webSocketWriteMessage(kind string, operation string, jsonMsg interface{}) {
	wsConnection.mux.Lock()

	// Check any client exist
	if wsConnection.client == nil {
		wsConnection.mux.Unlock()
		log.Println("Unable to write message: no client exist!")
		return
	}

	err := wsConnection.client.WriteJSON(Message{
		Kind:      kind,
		Operation: operation,
		Data:      jsonMsg,
	})
	if err != nil {
		log.Println("Write connection error:", err)
	}

	log.Println("Write message succeeded:", time.Now().Format(time.Kitchen))

	wsConnection.mux.Unlock()
}

func setupTrayIcon() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Short Cutteer")

	// Add default menu items in sequence
	menuLaunchUI := systray.AddMenuItem("Launch UI", "", true)
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", "", false)

	go func() {
		for {
			select {
			case <-menuLaunchUI.ClickedCh:
				x := exec.Command("explorer", httpURL).Start()
				fmt.Println(x)

			case <-menuQuit.ClickedCh:
				processInterrupted()

			case <-processInterruptSignal:
				processInterrupted()
			}
		}
	}()
}

func setupWindowsHook() {
	log.Println("Keyboard listener started ......")

	// Load all required DLLs
	_ = LoadDLLs()

	// Setup hook annd receive message
	go receiveHook()
}

func processInterrupted() {
	// Unhook Windows keyboard
	log.Println("Removing Windows hook ......")
	_, _ = UnhookWindowsHookEx(hhook)

	// Quit system tray
	log.Println("Removing sytem tray ......")
	systray.Quit()

	// Close db
	err := myDB.Close()
	if err != nil {
		panic(err)
	}

	// Exit
	os.Exit(1)
}

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
