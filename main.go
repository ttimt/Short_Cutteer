package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"unsafe"
)

var hhook HHOOK
var currentKeyStrokeSignal = make(chan rune)
var userCommands = make(map[string]string)
var bufferCommand string
var bufferLen int
var maxBufferLen int
var commandReady bool

func receiveHook() {
	// Declare a keyboard hook callback function (type HOOKPROC)
	hookCallback := func(code int, wParam WPARAM, lParam LPARAM) LRESULT {
		// If keystroke is pressed down
		if wParam == 256 {
			// Retrieve the keyboard hook struct
			keyboardHookData := (*tagKBDLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))

			// Retrieve current keystroke from keyboard hook struct's vkCode
			currentKeystroke := rune((*keyboardHookData).vkCode)

			// Send the keystroke to be processed
			select {
			case currentKeyStrokeSignal <- currentKeystroke:

			default:
				// Skip if the channel currentKeyStrokeSignal is busy
				// as it means that the keystroke is sent from the processHook
				// We want to ignore keystrokes that we sent ourself
			}
		}

		// Return CallNextHookEx result to allow keystroke to be displayed on user screen
		return CallNextHookEx(0, code, wParam, lParam)
	}

	// Install a Windows hook that listen to keyboard input
	hhook = SetWindowsHookExW(WH_KEYBOARD_LL, hookCallback, 0, 0)
	if hhook == 0 {
		panic("Failed to set Windows hook")
	}

	// Run hook processing goroutine
	go processHook()

	// Start retrieving message from the hook
	if !GetMessageW(0, 0, 0, 0) {
		panic("Failed to get message")
	}
}

// Process your received keystroke here
func processHook() {
	for {
		// Receive keystroke from channel
		currentKeyStroke := <-currentKeyStrokeSignal

		// Process keystroke
		fmt.Printf("Current key: %d\n", currentKeyStroke)

		if commandReady && (currentKeyStroke == VK_SPACE || currentKeyStroke == VK_TAB) {
			fmt.Println("Command ready in", uint16(userCommands[bufferCommand][0]))
			switch currentKeyStroke {
			case VK_SPACE:
			case VK_TAB:
				tagInputs := createKeyboardTagInputs(userCommands[bufferCommand])
				SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
			}

			bufferLen = 0
			bufferCommand = ""
			commandReady = false
			continue
		}

		commandReady = false
		if bufferLen >= maxBufferLen {
			bufferCommand = bufferCommand[1:]
			bufferLen--
		}

		switch {
		case currentKeyStroke == VK_OEM_2:
			bufferCommand += "/"
		case 65 <= currentKeyStroke && currentKeyStroke <= 90 && GetKeyState(VK_SHIFT)>>15 == 1: // Capital letter
			bufferCommand += string(currentKeyStroke)
		case 65 <= currentKeyStroke && currentKeyStroke <= 90: // Small letters
			bufferCommand += strings.ToLower(string(currentKeyStroke))
		default:
			bufferLen--
			// bufferCommand += string(currentKeyStroke)
		}
		bufferLen++
		fmt.Println("Current buffer:", bufferCommand)

		if _, ok := userCommands[bufferCommand]; ok {
			commandReady = true
		}

		// If left bracket, left bracket, space x2, right bracket, left arrow x2
		// If first double/single quotes, right quote, left arrow
		// If after ( or [ (bracket already complete) and enter, enter and left arrow x2? depends on default enter behavior
		// if after {  } and enter, delete, backspace, left arrow, enter, right arrow, enter x2, left arrow
		// if command xxx and tab, enter yyy
		// if command xxx and space, backspace xxx len and enter yyy
		// If brackets or quotes, just can copy text (mouse + keyboard to detect there is text selected), copy, left bracket/quote, space, paste, space, right bracket/quote
		// If 'shortcut key', copy and format and paste
	}
}

func defineCommands() {
	userCommands["/akey"] = ""
	userCommands["/adef"] = ""

	maxBufferLen = 5
}

func createKeyboardTagInputs(str string) []tagINPUT {
	var tagInputs []tagINPUT

	shiftDownInput := tagINPUT{
		inputType: INPUT_KEYBOARD,
		ki: KEYBDINPUT{
			WVk: VK_SHIFT,
		},
	}

	shiftUpInput := tagINPUT{
		inputType: INPUT_KEYBOARD,
		ki: KEYBDINPUT{
			WVk:     VK_SHIFT,
			DwFlags: KEYEVENTF_KEYUP,
		},
	}

	for _, v := range str {
		switch {
		case 65 <= v && v <= 90: // Capital letter
			key := tagINPUT{
				inputType: INPUT_KEYBOARD,
				ki: KEYBDINPUT{
					WVk: uint16(v),
				},
			}

			tagInputs = append(tagInputs, shiftDownInput, key, shiftUpInput)

		case 97 <= v && v <= 122: // Small letters
			key := tagINPUT{
				inputType: INPUT_KEYBOARD,
				ki: KEYBDINPUT{
					WVk: uint16(strings.ToUpper(string(v))[0]),
				},
			}

			tagInputs = append(tagInputs, key)

		case v == 40: // Left bracket
			key := tagINPUT{
				inputType: INPUT_KEYBOARD,
				ki: KEYBDINPUT{
					WVk: VK_NINE,
				},
			}

			tagInputs = append(tagInputs, shiftDownInput, key, shiftUpInput)
		case v == 41: // Right bracket
			key := tagINPUT{
				inputType: INPUT_KEYBOARD,
				ki: KEYBDINPUT{
					WVk: VK_ZERO,
				},
			}

			tagInputs = append(tagInputs, shiftDownInput, key, shiftUpInput)

		case v == 46: // Period
			key := tagINPUT{
				inputType: INPUT_KEYBOARD,
				ki: KEYBDINPUT{
					WVk: VK_OEM_PERIOD,
				},
			}

			tagInputs = append(tagInputs, key)

		case v == 32: // Space
			key := tagINPUT{
				inputType: INPUT_KEYBOARD,
				ki: KEYBDINPUT{
					WVk: VK_SPACE,
				},
			}

			tagInputs = append(tagInputs, key)

		default:
			// Don't process key if not specified above
			// Or keys like backspace, delete, and weird symbols will be added to the buffer
		} // END switch
	}

	return tagInputs
}

func main() {
	log.Println("Start")

	// Load all required DLLs
	LoadDLLs()

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
	UnhookWindowsHookEx(hhook)
}
