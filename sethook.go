package main

import (
	"fmt"
	"log"
	"unsafe"

	. "github.com/ttimt/Short_Cutteer/hook"
	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

const (
	nrOfKeys    = 129 // Number of possible keys
	invalidRune = -1
	invalidWord = 0
)

var (
	currentKeyStrokeSignal = make(chan rune)
	hhook                  HHOOK
	autoCompleteJustDone   bool

	// Store hook keys
	keys                            []Key
	keysByKeyCodeWithShiftOrCapital = make(map[uint16]*Key)
	keysByKeyCodeWithoutShift       = make(map[uint16]*Key)
	keysByChar                      = make(map[rune]*Key)
)

// Initialize DLLs and start listening to keyboard hook
func setupWindowsHook() {
	log.Println("Keyboard listener started ......")

	// Load all required DLLs
	_ = LoadDLLs()

	// Setup hook annd receive message
	go receiveHook()
}

// Detect hook
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
		isShiftEnabled := getKeyState(shiftKeyState)
		isCapsEnabled := getKeyState(capsLockState, true)

		key := getKeyByKeyCode(uint16(currentKeyStroke), isShiftEnabled, isCapsEnabled)
		var char rune
		if key == nil {
			char = invalidRune
		} else {
			char = key.Char
		}

		if isAutoComplete(char) {
			// Process auto complete
			tagInputs := processAutoComplete(char, isShiftEnabled, isCapsEnabled)

			// Send input
			_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
			continue
		}

		// Reset if character is not a letter/symbol/number
		if char == -1 {
			bufferStr = ""

			if currentKeyStroke != VK_SHIFT && currentKeyStroke != VK_LSHIFT && currentKeyStroke != VK_RSHIFT {
				autoCompleteJustDone = false
			}

			continue
		}

		// User pre-collectionCommands can be CTRL, ALT, SHIFT and 1 letter afterwards
		// User can create shortcut MODIFIER + a key or text collectionCommands + tab or space (remove collectionCommands)
		switch char {
		case '\b':
			if len(bufferStr) > 0 {
				bufferStr = bufferStr[:len(bufferStr)-1]
			}

			if autoCompleteJustDone {
				// Send input
				tagInputs := getKeyByKeyCode(VK_DELETE).KeyPress()
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
			}
		case '\r':
			bufferStr += windowsNewLine
			bufferStr = ""
		case '\t':
			if str, ok := userCommands[bufferStr]; ok {
				// Send input
				tagInputs := createTagInputs(str.Output, isShiftEnabled, isCapsEnabled)
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
			}

			bufferStr = ""
		case ' ':
			if str, ok := userCommands[bufferStr]; ok {
				// Send input
				tagInputsBackspace := getKeyByKeyCode(VK_BACK).KeyPress(len(bufferStr) + 1)
				_, _ = SendInput(uint(len(tagInputsBackspace)), (*LPINPUT)(&tagInputsBackspace[0]), int(unsafe.Sizeof(tagInputsBackspace[0])))

				tagInputs := createTagInputs(str.Output, isShiftEnabled, isCapsEnabled)
				_, _ = SendInput(uint(len(tagInputs)), (*LPINPUT)(&tagInputs[0]), int(unsafe.Sizeof(tagInputs[0])))
			}

			bufferStr = ""
		default:
			// If buffer full, trim
			if len(bufferStr) >= maxBufferLen && len(bufferStr) > 0 {
				bufferStr = bufferStr[1:]
			}

			bufferStr += string(char)
		}

		autoCompleteJustDone = false
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
