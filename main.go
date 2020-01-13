package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"unsafe"
)

var wg sync.WaitGroup
var hhook HHOOK
var keyDown bool

func receiveHook(ctx context.Context, ch chan *tagKBDLLHOOKSTRUCT) {
	var fn HOOKPROC

	fn = func(code int, wParam WPARAM, lParam LPARAM) LRESULT {
		if wParam == 256 {
			// Retrieve the KBDLLHOOKSTRUCT
			char := (*tagKBDLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))

			// Process your received key here
			curChar := byte((*char).vkCode)
			fmt.Println("Current character:", curChar)

			if curChar == '9' && !keyDown {
				shiftKeyState := GetKeyState(VK_SHIFT) >> 15
				fmt.Println("shift state", shiftKeyState)
				if shiftKeyState == 1 {
					keyDown = true
					var input tagINPUT
					input.inputType = 1
					input.ki.WVk = '0'

					var input2 tagINPUT
					input2.inputType = 1
					input2.ki.WVk = '9'

					var input5 tagINPUT
					input5.inputType = 1
					input5.ki.WVk = 0x25

					allInput := make([]tagINPUT, 0)

					allInput = append(allInput, input2, input, input5)

					SendInput(uint(len(allInput)), LPINPUT(allInput[0]), int(unsafe.Sizeof(allInput[0])))

					keyDown = false

					return -1
				}
			}

			ch <- char
		}

		return CallNextHookEx(0, code, wParam, lParam)
	}

	go func() {
		hhook = SetWindowsHookExW(WH_KEYBOARD_LL, fn, 0, 0)
		if hhook == 0 {
			panic("Failed to set windows hook")
		}

		GetMessageW(0, 0, 0, 0)
	}()

	<-ctx.Done()
}

func main() {
	fmt.Println("Start")

	// Load user32.dll
	LoadDLLs()

	var isInterrupted bool
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	ch := make(chan *tagKBDLLHOOKSTRUCT, 1)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		receiveHook(ctx, ch)
		wg.Done()
	}()

	for {
		if isInterrupted {
			cancel()
			break
		}

		select {
		case <-signalChan:
			isInterrupted = true
		case c := <-ch:
			fmt.Printf("Char: %q\n", byte((*c).vkCode))
		}
	}

	wg.Wait()
	UnhookWindowsHookEx(hhook)
}
