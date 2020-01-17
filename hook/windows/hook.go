package windows

import (
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
	LPINPUT   TagINPUT

	// HOOKPROC Callback function after SendMessageW function is called (Keyboard input received)
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-hookproc
	//
	// LPARAM is a pointer to a KBDLLHOOKSTRUCT struct :
	// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)
	HOOKPROC func(int, WPARAM, LPARAM) LRESULT
)

// Low-level keyboard input event
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
type TagKBDLLHOOKSTRUCT struct {
	vkCode      DWORD
	scanCode    DWORD
	flags       DWORD
	time        DWORD
	dwExtraInfo ULONG_PTR
}

// Input events
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-input
type TagINPUT struct {
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
	VK_BACK       = 0x08 // Key: BACKSPACE
	VK_TAB        = 0x09 // Key: TAB
	VK_RETURN     = 0x0D // Key: Enter
	VK_SHIFT      = 0x10 // Key: SHIFT
	VK_CONTROL    = 0x11 // Key: CTRL
	VK_MENU       = 0x12 // Key: ALT
	VK_CAPITAL    = 0x14 // Key: CAPS LOCK
	VK_ESCAPE     = 0x1B // Key: ESC
	VK_SPACE      = 0x20 // Key: SPACEBAR
	VK_END        = 0x23 // Key: END
	VK_HOME       = 0x24 // Key: HOME
	VK_LEFT       = 0x25 // Key: LEFT ARROW
	VK_UP         = 0x26 // Key: UP ARROW
	VK_RIGHT      = 0x27 // Key: RIGHT ARROW
	VK_DOWN       = 0x28 // Key: DOWN ARROW
	VK_DELETE     = 0x2E // Key: DEL
	VK_ZERO       = 0x30 // Key: ZERO
	VK_ONE        = 0x31 // Key: ONE
	VK_TWO        = 0x32 // Key: TWO
	VK_THREE      = 0x33 // Key: THREE
	VK_FOUR       = 0x34 // Key: FOUR
	VK_FIVE       = 0x35 // Key: FIVE
	VK_SIX        = 0x36 // Key: SIX
	VK_SEVEN      = 0x37 // Key: SEVEN
	VK_EIGHT      = 0x38 // Key: EIGHT
	VK_NINE       = 0x39 // Key: NINE
	VK_A          = 0x41 // Key: A
	VK_B          = 0x42 // Key: B
	VK_C          = 0x43 // Key: C
	VK_D          = 0x44 // Key: D
	VK_E          = 0x45 // Key: E
	VK_F          = 0x46 // Key: F
	VK_G          = 0x47 // Key: G
	VK_H          = 0x48 // Key: H
	VK_I          = 0x49 // Key: I
	VK_J          = 0x4A // Key: J
	VK_K          = 0x4B // Key: K
	VK_L          = 0x4C // Key: L
	VK_M          = 0x4D // Key: M
	VK_N          = 0x4E // Key: N
	VK_O          = 0x4F // Key: O
	VK_P          = 0x50 // Key: P
	VK_Q          = 0x51 // Key: Q
	VK_R          = 0x52 // Key: R
	VK_S          = 0x53 // Key: S
	VK_T          = 0x54 // Key: T
	VK_U          = 0x55 // Key: U
	VK_V          = 0x56 // Key: V
	VK_W          = 0x57 // Key: W
	VK_X          = 0x58 // Key: X
	VK_Y          = 0x59 // Key: Y
	VK_Z          = 0x5A // Key: Z
	VK_NUMPAD0    = 0x60 // Key: Numeric keypad 0
	VK_NUMPAD1    = 0x61 // Key: Numeric keypad 1
	VK_NUMPAD2    = 0x62 // Key: Numeric keypad 2
	VK_NUMPAD3    = 0x63 // Key: Numeric keypad 3
	VK_NUMPAD4    = 0x64 // Key: Numeric keypad 4
	VK_NUMPAD5    = 0x65 // Key: Numeric keypad 5
	VK_NUMPAD6    = 0x66 // Key: Numeric keypad 6
	VK_NUMPAD7    = 0x67 // Key: Numeric keypad 7
	VK_NUMPAD8    = 0x68 // Key: Numeric keypad 8
	VK_NUMPAD9    = 0x69 // Key: Numeric keypad 9
	VK_MULTIPLY   = 0x6A // Key: Numeric keypad *
	VK_ADD        = 0x6B // Key: Numeric keypad +
	VK_SUBTRACT   = 0x6D // Key: Numeric keypad -
	VK_DECIMAL    = 0x6E // Key: Numeric keypad .
	VK_DIVIDE     = 0x6F // Key: Numeric keypad /
	VK_LSHIFT     = 0xA0 // Key: Left SHIFT
	VK_RSHIFT     = 0xA1 // Key: Right SHIFT
	VK_LCONTROL   = 0xA2 // Key: Left CONTROL
	VK_RCONTROL   = 0xA3 // Key: RIGHT CONTROL
	VK_OEM_1      = 0xBA // Key: ";:
	VK_OEM_PLUS   = 0xBB // Key: +
	VK_OEM_COMMA  = 0xBC // Key: ,
	VK_OEM_MINUS  = 0xBD // Key: -
	VK_OEM_PERIOD = 0xBE // Key: .
	VK_OEM_2      = 0xBF // Key: "/?"
	VK_OEM_3      = 0xc0 // Key: "`~"
	VK_OEM_4      = 0xDB // Key: "[{"
	VK_OEM_5      = 0xDC // Key: "\|"
	VK_OEM_6      = 0xDD // Key: "]}"
	VK_OEM_7      = 0xDE // Key: "'""

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

	// Window Messages: Code to get and send messages between applications
	// https://docs.microsoft.com/en-us/windows/win32/winmsg/window-messages
	WM_SETTEXT       = 0x000C
	WM_GETTEXT       = 0x000D
	WM_GETTEXTLENGTH = 0x000E
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
	winDLLUser32_SendMessageW            = winDLLUser32.NewProc("SendMessageW")
)

// LoadDLLs loads all required DLLs and return if error(s) occurred
func LoadDLLs() error {
	// Load user32.dll
	err := winDLLUser32.Load()

	return err
}

// CallNextHookEx Pass the hook information to the next hook procedure
// A hook procedure can call this function either before or after processing the hook information
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callnexthookex
func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) (LRESULT, error) {
	result, _, err := winDLLUser32_ProcCallNextHookEx.Call(uintptr(hhk), uintptr(nCode), uintptr(wParam), uintptr(lParam))

	return LRESULT(result), err
}

// SetWindowsHookExW Install hook procedure into a hhook chain
// Result is null if error
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
func SetWindowsHookExW(idHook int, lpfn HOOKPROC, hmod HINSTANCE, dwThreadID DWORD) (HHOOK, error) {
	result, _, err := winDLLUser32_ProcSetWindowsHookExW.Call(uintptr(idHook), windows.NewCallback(lpfn), uintptr(hmod), uintptr(dwThreadID))

	return HHOOK(result), err
}

// UnhookWindowsHookEx Remove a hook procedure installed in a hook chain by the SetWindowsHookEx function
// Result is zero if error
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
func UnhookWindowsHookEx(hhk HHOOK) (bool, error) {
	result, _, err := winDLLUser32_ProcUnhookWindowsHookEx.Call(uintptr(hhk))

	return result != 0, err
}

// GetMessageW Retrieves a message
// Result is -1 if error
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessagew
func GetMessageW(lpMsg LPMSG, hWnd HWND, wMsgFilterMin uint, wMsgFilterMax uint) (bool, error) {
	result, _, err := winDLLUser32_GetMessageW.Call(uintptr(lpMsg), uintptr(hWnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax))

	return result != -1, err
}

// SendInput Simulate keyboard inputs to the operating system
// Result is zero if error
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput
func SendInput(cInputs uint, pInputs *LPINPUT, cbSize int) (uint, error) {

	result, _, err := winDLLUser32_SendInput.Call(uintptr(cInputs), uintptr(unsafe.Pointer(pInputs)), uintptr(cbSize))

	return uint(result), err
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
func GetKeyState(nVirtKey int) (SHORT, error) {
	result, _, err := winDLLUser32_GetKeyState.Call(uintptr(nVirtKey))

	return SHORT(result), err
}

// Retrieve a handle to the user active foreground window
// Result is null if empty handle
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
func GetForegroundWindow() (HWND, error) {
	result, _, err := winDLLUser32_GetForegroundWindow.Call()

	return HWND(result), err
}

// Send the specified message to a window.
// The method does not return until the window procedure processed the message
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage
func SendMessageW(hWnd HWND, Msg uint, wParam WPARAM, lParam LPARAM) (LRESULT, error) {
	result, _, err := winDLLUser32_SendMessageW.Call(uintptr(hWnd), uintptr(Msg), uintptr(wParam), uintptr(lParam))

	return LRESULT(result), err
}
