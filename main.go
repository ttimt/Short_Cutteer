// +build windows

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"unsafe"

	"github.com/ttimt/systray"

	_ "github.com/lxn/walk"

	. "github.com/ttimt/Short_Cutteer/hook/windows"
	icon "github.com/ttimt/Short_Cutteer/icons"
)

const (
	htmlFilePath = "html/index.html"
	httpPort     = 8080
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
)

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

	maxBufferLen = 5
}

func init() {
	// Setup process interrupt signal
	signal.Notify(processInterruptSignal, os.Interrupt)
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
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Handle web socket
	_, _ = w.Write([]byte("Hello world!!"))
}

func setupHTTPServer() {
	// Setup mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, htmlFilePath)
	})

	mux.HandleFunc("/ws", handleWebSocket)

	// Concurrently run web server
	go func() {
		log.Println("Started listening on", httpPort)
		log.Fatal(http.ListenAndServe(httpPortStr, mux))
	}()
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

	// Exit
	os.Exit(1)
}
