package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"unsafe"
)

var hhook HHOOK
var keyDown bool
var isProcessInterrupted bool
var currentKeyStrokeSignal = make(chan rune)

func receiveHook() {
	// Declare a keyboard hook callback function (type HOOKPROC)
	hookCallback := func(code int, wParam WPARAM, lParam LPARAM) LRESULT {
		if wParam == 256 {
			// Retrieve the KBDLLHOOKSTRUCT
			keyboardHookData := (*tagKBDLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))

			// Retrieve current keystroke from keyboard hook struct's vkCode
			currentKeystroke := (*keyboardHookData).vkCode

			// Process your received key here
			if currentKeystroke == VK_NINE && !keyDown {
				shiftKeyState := GetKeyState(VK_SHIFT) >> 15

				if shiftKeyState == -1 {
					keyDown = true
					var input tagINPUT
					input.inputType = INPUT_KEYBOARD
					input.ki.WVk = VK_ZERO

					var input2 tagINPUT
					input2.inputType = INPUT_KEYBOARD
					input2.ki.WVk = VK_NINE

					var input5 tagINPUT
					input5.inputType = INPUT_KEYBOARD
					input5.ki.WVk = VK_LEFT

					var input6 tagINPUT
					input6.inputType = INPUT_KEYBOARD
					input6.ki.WVk = VK_SHIFT
					input6.ki.DwFlags = KEYEVENTF_KEYUP

					allInput := make([]tagINPUT, 0)

					allInput = append(allInput, input2, input, input6, input5)

					SendInput(uint(len(allInput)), (*LPINPUT)(&allInput[0]), int(unsafe.Sizeof(allInput[0])))
					keyDown = false

					// Call CallNextHookEx to allow other applications using Windows hook to process the keystroke as well
					CallNextHookEx(0, code, wParam, lParam)

					//
					return -1
				}
			}

			// Send the keystroke to be printed
			currentKeyStrokeSignal <- rune(currentKeystroke)
		}

		// Return CallNextHookEx result to allow keystroke to be displayed on user screen
		return CallNextHookEx(0, code, wParam, lParam)
	}

	// Install a Windows hook that listen to keyboard input
	hhook = SetWindowsHookExW(WH_KEYBOARD_LL, hookCallback, 0, 0)
	if hhook == 0 {
		panic("Failed to set Windows hook")
	}

	// Start retrieving message from the hook
	if !GetMessageW(0, 0, 0, 0) {
		panic("Failed to get message")
	}
}

func main() {
	log.Println("Start")

	// Load user32.dll
	LoadDLLs()

	// Setup process interrupt signal
	processInterruptSignal := make(chan os.Signal)
	signal.Notify(processInterruptSignal, os.Interrupt)

	// Setup hook annd receive message
	go receiveHook()

	// Detect process interrupt or keystroke input
	for {
		// Unhook Windows keyboard and break the loop if process is interrupted
		if isProcessInterrupted {
			UnhookWindowsHookEx(hhook)
			break
		}

		select {
		case <-processInterruptSignal:
			isProcessInterrupted = true
		case c := <-currentKeyStrokeSignal:
			fmt.Printf("Current keystroke: %q\n", c)
		}
	}
}
