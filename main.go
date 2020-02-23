// +build windows

package main

import (
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
	"github.com/HouzuoGuo/tiedot/dberr"
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
)

var (
	hhook                  HHOOK
	currentKeyStrokeSignal = make(chan rune)
	userCommands           = make(map[string]string)
	maxBufferLen           int
	bufferStr              string

	httpPortStr = ":" + strconv.Itoa(httpPort)
	httpURL     = "http://localhost" + httpPortStr

	processInterruptSignal = make(chan os.Signal)

	wsUpgrader   = websocket.Upgrader{}
	wsConnection webSocketConnection

	myDB *db.DB

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
	Title       string
	Description string
	Command     string
	Output      string
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

		// User pre-commands can be CTRL, ALT, SHIFT and 1 letter afterwards
		// User can create shortcut MODIFIER + a key or text commands + tab or space (remove commands)
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
				tagInputs := createTagInputs(str)
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
				bufferStr = ""
			}
		case ' ':
			if str, ok := userCommands[bufferStr]; ok {
				// Send input
				tagInputsBackspace := multiplyTagInputKey(tagInputBackspaceDown(), len(bufferStr)+1)
				_, _ = SendInput(uint(len(tagInputsBackspace)), (*LPINPUT)(&tagInputsBackspace[0]), int(unsafe.Sizeof(tagInputsBackspace[0])))
				tagInputs := createTagInputs(str)
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
				bufferStr = ""
			}
		case '`':
			var c []Command
			c = append(c, Command{
				Title:       "hey",
				Description: "description!",
				Command:     "/akey",
				Output:      "VALUEOBJECTKEY",
			})
			c = append(c, Command{
				Title:       "he nonono",
				Description: "description!!!!!",
				Command:     "/adef",
				Output:      "VALUE OBJECT DEF",
			})

			fmt.Println("Sending:", c)
			webSocketWriteMessage(c)
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

func defineCommands() {
	// To move to UI
	userCommands["akey"] = "VALUE( object.Key() )"
	userCommands["adef"] = "VALUE( object.DefinitionName() )"

	maxBufferLen = 4
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

	// Test DB
	testDB()
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

	// Start reading message
	webSocketReadMessage()
}

func webSocketReadMessage() {
	for {
		// Read message
		_, m, err := wsConnection.client.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Println("Web socket connection closed:", err)
			} else {
				log.Println("Web socket error:", err)
			}

			break
		}

		// Read success
		fmt.Println("Read message succeeded!", m, time.Now().Format(time.Kitchen))
	}

	// Close the connection at the end if read fails
	_ = wsConnection.client.Close()
}

func webSocketWriteMessage(jsonMsg interface{}) {
	wsConnection.mux.Lock()

	// Check any client exist
	if wsConnection.client == nil {
		wsConnection.mux.Unlock()
		log.Println("Unable to write message: no client exist!")
		return
	}

	err := wsConnection.client.WriteJSON(jsonMsg)
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

	// Define commands
	defineCommands()

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
