package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"unsafe"

	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

var hhook HHOOK
var currentKeyStrokeSignal = make(chan rune)
var userCommands = make(map[string]string)
var maxBufferLen int

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
		fmt.Printf("Current key: %d 0x%x %c\n", currentKeyStroke, currentKeyStroke, currentKeyStroke)

		shiftKeyState, _ := GetKeyState(VK_SHIFT)
		capsLockState, _ := GetKeyState(VK_CAPITAL)

		_, char, _ := findAllKeyCode(uint16(currentKeyStroke), 0, getKeyStateBool(shiftKeyState), getKeyStateBool(capsLockState, true))
		if char != -1 {
			fmt.Printf("Real string is: %c\n", char)
		}

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
	userCommands["/akey"] = "VALUE( object.Key() )"
	userCommands["/adef"] = "VALUE( object.DefinitionName() )"

	maxBufferLen = 5
}

func main() {
	log.Println("Start")

	// Load all required DLLs
	_ = LoadDLLs()

	// Define commands
	defineCommands()

	// Setup process interrupt signal
	processInterruptSignal := make(chan os.Signal)
	signal.Notify(processInterruptSignal, os.Interrupt)

	// Setup hook annd receive message
	go receiveHook()

	// Wait for process to be interrupted
	<-processInterruptSignal

	// Unhook Windows keyboard
	fmt.Println("Removing Windows hook ......")
	_, _ = UnhookWindowsHookEx(hhook)
}
