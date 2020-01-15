package main

import (
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

type (
	SHORT     int16
	DWORD     uint32
	ULONG_PTR uint32
	LRESULT   int64
	LPARAM    int64
	HHOOK     uintptr
	HINSTANCE uintptr
	HWND      uintptr
	LPMSG     uintptr
	WPARAM    uintptr
	LPINPUT   tagINPUT

	// HOOKPROC Callback function after SendMessage function is called (Keyboard input received)
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-hookproc
	//
	// LPARAM is a pointer to a KBDLLHOOKSTRUCT struct :
	// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)
	HOOKPROC func(int, WPARAM, LPARAM) LRESULT
)

// Low-level keyboard input event
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
type tagKBDLLHOOKSTRUCT struct {
	vkCode      DWORD
	scanCode    DWORD
	flags       DWORD
	time        DWORD
	dwExtraInfo ULONG_PTR
}

// Input events
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-input
type tagINPUT struct {
	inputType uint32
	ki        KEYBDINPUT
	padding   uint64
}

// KEYBDINPUT Simulated keyboard event
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-keybdinput
type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

const (
	// Types of hook procedure:
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
	WH_KEYBOARD_LL = 13

	// Virtual key codes:
	// https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	VK_BACK       = 0x08
	VK_TAB        = 0x09
	VK_RETURN     = 0x0D // Enter key
	VK_SHIFT      = 0x10
	VK_CONTROL    = 0x11
	VK_MENU       = 0x12
	VK_CAPITAL    = 0x14
	VK_SPACE      = 0x20
	VK_END        = 0x23
	VK_HOME       = 0x24
	VK_LEFT       = 0x25
	VK_UP         = 0x26
	VK_RIGHT      = 0x27
	VK_DOWN       = 0x28
	VK_ZERO       = 0x30
	VK_ONE        = 0x31
	VK_TWO        = 0x32
	VK_THREE      = 0x33
	VK_FOUR       = 0x34
	VK_FIVE       = 0x35
	VK_SIX        = 0x36
	VK_SEVEN      = 0x37
	VK_EIGHT      = 0x38
	VK_NINE       = 0x39
	VK_LCONTROL   = 0xA2
	VK_RCONTROL   = 0xA3
	VK_OEM_1      = 0xBA // ';:' key
	VK_OEM_PERIOD = 0xBE // '.' key
	VK_OEM_2      = 0xBF // '/?' key
	VK_OEM_4      = 0xDB // '[{' key
	VK_OEM_6      = 0xDD // ']}' key
	VK_OEM_7      = 0xDE // 'single-quote/double-quote' key

	// INPUT_MOUSE Types of input event:
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-input#members
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2

	// KEYEVENTF_EXTENDEDKEY Keystroke for dwFlags in KEYBDINPUT
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-keybdinput#members
	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_SCANCODE    = 0x0008
	KEYEVENTF_UNICODE     = 0x0004
)

var (
	// Microsoft Windows DLLs
	winDLLUser32 = windows.NewLazyDLL("user32.dll")

	// User32.dll procedures
	winDLLUser32_ProcCallNextHookEx      = winDLLUser32.NewProc("CallNextHookEx")
	winDLLUser32_ProcSetWindowsHookExW   = winDLLUser32.NewProc("SetWindowsHookExW")
	winDLLUser32_ProcUnhookWindowsHookEx = winDLLUser32.NewProc("UnhookWindowsHookEx")
	winDLLUser32_GetMessageW             = winDLLUser32.NewProc("GetMessageW")
	winDLLUser32_SendInput               = winDLLUser32.NewProc("SendInput")
	winDLLUser32_GetKeyState             = winDLLUser32.NewProc("GetKeyState")
	winDLLUser32_GetForegroundWindow     = winDLLUser32.NewProc("GetForegroundWindow")
	winDLLUser32_SendMessage             = winDLLUser32.NewProc("SendMessage")
)

// LoadDLLs loads all required DLLs and panic if error(s) occurred
func LoadDLLs() {
	// Load user32.dll
	err := winDLLUser32.Load()
	if err != nil {
		panic("LoadDLL error" + err.Error())
	}
}

// CallNextHookEx Pass the hook information to the next hook procedure
// A hook procedure can call this function either before or after processing the hook information
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callnexthookex
func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	result, _, _ := winDLLUser32_ProcCallNextHookEx.Call(uintptr(hhk), uintptr(nCode), uintptr(wParam), uintptr(lParam))

	return LRESULT(result)
}

// SetWindowsHookExW Install hook procedure into a hhook chain
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
func SetWindowsHookExW(idHook int, lpfn HOOKPROC, hmod HINSTANCE, dwThreadID DWORD) HHOOK {
	result, _, _ := winDLLUser32_ProcSetWindowsHookExW.Call(uintptr(idHook), windows.NewCallback(lpfn), uintptr(hmod), uintptr(dwThreadID))

	return HHOOK(result)
}

// UnhookWindowsHookEx Remove a hook procedure installed in a hook chain by the SetWindowsHookEx function
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
func UnhookWindowsHookEx(hhk HHOOK) bool {
	result, _, err := winDLLUser32_ProcUnhookWindowsHookEx.Call(uintptr(hhk))

	if result == 0 {
		log.Println("UnhookWindowsHookEx error:", err)
		return false
	}

	return true
}

// GetMessageW Retrieves a message
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessagew
func GetMessageW(lpMsg LPMSG, hWnd HWND, wMsgFilterMin uint, wMsgFilterMax uint) bool {
	res, _, err := winDLLUser32_GetMessageW.Call(uintptr(lpMsg), uintptr(hWnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax))

	if res == 0 {
		log.Println("GetMessageW error:", err)
		return false
	}

	return true
}

// SendInput Simulate keyboard inputs to the operating system
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput
func SendInput(cInputs uint, pInputs *LPINPUT, cbSize int) uint {

	result, _, err := winDLLUser32_SendInput.Call(uintptr(cInputs), uintptr(unsafe.Pointer(pInputs)), uintptr(cbSize))

	if result == 0 {
		log.Println("SendInput error:", err)
		return 0
	}

	return uint(result)
}

// GetKeyState Retrieves the status of the specified virtual key
// The status specifies whether the key is up, down or toggled (on, off - alternating each time the key is pressed)
//
// Returned bits = 16 bits
// If high-order bit is 1, the key is down, otherwise it is up
// If low-order bit is 1, the key is toggled on, otherwise the key is off
//
// Since SHORT is int16, a negative value will indicates high-order bit is 1
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeystate
func GetKeyState(nVirtKey int) SHORT {
	result, _, _ := winDLLUser32_GetKeyState.Call(uintptr(nVirtKey))

	return SHORT(result)
}

// Retrieve a handle to the user active foreground window
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
func GetForegroundWindow() HWND {
	result, _, err := winDLLUser32_GetForegroundWindow.Call()

	if result == 0 {
		log.Println("GetForegroundWindow error:", err)
	}

	return HWND(result)
}

// Send the specified message to a window.
// The method does not return until the window procedure processed the message
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage
func SendMessage(hWnd HWND, Msg uint, wParam WPARAM, lParam LPARAM) LRESULT {
	result, _, _ := winDLLUser32_SendMessage.Call(uintptr(hWnd), uintptr(Msg), uintptr(wParam), uintptr(lParam))

	return LRESULT(result)
}
