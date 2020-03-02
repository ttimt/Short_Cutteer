package main

import (
	"log"

	. "github.com/ttimt/Short_Cutteer/hook"
	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

// Create []TagInputs that can be used in SendInput() function to send keys that
// can be represented in string or escape sequences
func createTagInputs(strToSend string, isShiftEnabled, isCapsEnabled bool) (tagInputs []TagINPUT) {
	// Shift Key
	shiftKey := getKeyByKeyCode(VK_SHIFT)

	for _, c := range strToSend {
		key := getKeyByChar(c)

		if key == nil {
			continue
		}

		if key.IsShiftNeeded && !isShiftEnabled || key.IsCapitalLetter && !IsCapitalLetterEnabled(isShiftEnabled, isCapsEnabled) {
			tagInputs = append(tagInputs, shiftKey.KeyHold())
		}

		tagInputs = append(tagInputs, key.KeyPress()...)

		// Release all keys
		if key.IsShiftNeeded && !isShiftEnabled || key.IsCapitalLetter && !IsCapitalLetterEnabled(isShiftEnabled, isCapsEnabled) {
			tagInputs = append(tagInputs, shiftKey.KeyRelease())
		}
	}

	return
}

// Process the return value from GetKeyState.
//
// If the key state needed is not key down/up but key toggled like caps lock,
// then pass in true in the second bool parameter
func getKeyState(state SHORT, checkToggle ...bool) bool {
	if len(checkToggle) > 0 {
		return state&1 == 1
	}

	return state < 0
}

// Check whether key is a capital letter:
// First state is SHIFT state
// Second state is CAPS state
func IsCapitalLetterEnabled(shiftKeyState, capsLockKeyState bool) bool {
	if !capsLockKeyState && !shiftKeyState || capsLockKeyState && shiftKeyState {
		return false
	}

	return true
}

// Search through existing hook keys by its key code.
// Input first parameter as true if shift key enabled.
// Input second parameter as true if capital key enabled
func getKeyByKeyCode(keyCode uint16, modifiers ...bool) *Key {
	if len(modifiers) > 2 {
		panic("Error input parameter to getKeyByKeyCode")
	}

	var key *Key

	if len(modifiers) > 0 && (modifiers[0] || IsCapitalLetterEnabled(modifiers[0], modifiers[1])) {
		key, _ = keysByKeyCodeWithShiftOrCapital[keyCode]
	}

	if key == nil || (!IsCapitalLetterEnabled(modifiers[0], modifiers[1]) && 65 <= key.Char && key.Char <= 90) {
		key, _ = keysByKeyCodeWithoutShift[keyCode]
	}

	if key == nil {
		log.Printf("Key code does not exist: 0x%0x", keyCode)
	}

	return key
}

// Search through existing hook keys by its char
func getKeyByChar(char rune) *Key {
	key, _ := keysByChar[char]

	if key == nil {
		log.Printf("Key code does not exist: %s", string(char))
	}

	return key
}

// Construct all necessary keys
func createAllHookKeys() {
	// Initialize memory
	keys = make([]Key, nrOfKeys)

	// Counter
	i := 0

	// Function to increment counter
	incrementCounter := func(key Key) {
		keys[i] = key

		// Update maps
		if key.IsShiftNeeded || key.IsCapitalLetter {
			keysByKeyCodeWithShiftOrCapital[key.KeyCode] = &key
		} else {
			keysByKeyCodeWithoutShift[key.KeyCode] = &key
		}

		if key.Char != invalidRune {
			keysByChar[key.Char] = &key
		}

		i++
	}

	incrementCounter(CreateHookKey(VK_BACK, '\b'))
	incrementCounter(CreateHookKey(VK_TAB, '\t'))
	incrementCounter(CreateHookKey(VK_RETURN, '\n'))
	incrementCounter(CreateHookKey(VK_SHIFT, invalidRune))
	incrementCounter(CreateHookKey(VK_CONTROL, invalidRune))
	incrementCounter(CreateHookKey(VK_MENU, invalidRune))
	incrementCounter(CreateHookKey(VK_CAPITAL, invalidRune))
	incrementCounter(CreateHookKey(VK_ESCAPE, invalidRune))
	incrementCounter(CreateHookKey(VK_SPACE, ' '))
	incrementCounter(CreateHookKey(VK_END, invalidRune))
	incrementCounter(CreateHookKey(VK_HOME, invalidRune))
	incrementCounter(CreateHookKey(VK_LEFT, invalidRune))
	incrementCounter(CreateHookKey(VK_UP, invalidRune))
	incrementCounter(CreateHookKey(VK_RIGHT, invalidRune))
	incrementCounter(CreateHookKey(VK_DOWN, invalidRune))
	incrementCounter(CreateHookKey(VK_DELETE, invalidRune))
	incrementCounter(CreateHookKey(VK_ZERO, '0'))
	incrementCounter(CreateHookKey(VK_ZERO, ')', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_ONE, '1'))
	incrementCounter(CreateHookKey(VK_ONE, '!', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_TWO, '2'))
	incrementCounter(CreateHookKey(VK_TWO, '@', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_THREE, '3'))
	incrementCounter(CreateHookKey(VK_THREE, '#', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_FOUR, '4'))
	incrementCounter(CreateHookKey(VK_FOUR, '$', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_FIVE, '5'))
	incrementCounter(CreateHookKey(VK_FIVE, '%', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_SIX, '6'))
	incrementCounter(CreateHookKey(VK_SIX, '^', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_SEVEN, '7'))
	incrementCounter(CreateHookKey(VK_SEVEN, '&', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_EIGHT, '8'))
	incrementCounter(CreateHookKey(VK_EIGHT, '*', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_NINE, '9'))
	incrementCounter(CreateHookKey(VK_NINE, '(', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_A, 'A', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_A, 'a'))
	incrementCounter(CreateHookKey(VK_B, 'B', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_B, 'b'))
	incrementCounter(CreateHookKey(VK_C, 'C', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_C, 'c'))
	incrementCounter(CreateHookKey(VK_D, 'D', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_D, 'd'))
	incrementCounter(CreateHookKey(VK_E, 'E', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_E, 'e'))
	incrementCounter(CreateHookKey(VK_F, 'F', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_F, 'f'))
	incrementCounter(CreateHookKey(VK_G, 'G', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_G, 'g'))
	incrementCounter(CreateHookKey(VK_H, 'H', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_H, 'h'))
	incrementCounter(CreateHookKey(VK_I, 'I', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_I, 'i'))
	incrementCounter(CreateHookKey(VK_J, 'J', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_J, 'j'))
	incrementCounter(CreateHookKey(VK_K, 'K', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_K, 'k'))
	incrementCounter(CreateHookKey(VK_L, 'L', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_L, 'l'))
	incrementCounter(CreateHookKey(VK_M, 'M', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_M, 'm'))
	incrementCounter(CreateHookKey(VK_N, 'N', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_N, 'n'))
	incrementCounter(CreateHookKey(VK_O, 'O', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_O, 'o'))
	incrementCounter(CreateHookKey(VK_P, 'P', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_P, 'p'))
	incrementCounter(CreateHookKey(VK_Q, 'Q', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_Q, 'q'))
	incrementCounter(CreateHookKey(VK_R, 'R', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_R, 'r'))
	incrementCounter(CreateHookKey(VK_S, 'S', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_S, 's'))
	incrementCounter(CreateHookKey(VK_T, 'T', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_T, 't'))
	incrementCounter(CreateHookKey(VK_U, 'U', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_U, 'u'))
	incrementCounter(CreateHookKey(VK_V, 'V', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_V, 'v'))
	incrementCounter(CreateHookKey(VK_W, 'W', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_W, 'w'))
	incrementCounter(CreateHookKey(VK_X, 'X', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_X, 'x'))
	incrementCounter(CreateHookKey(VK_Y, 'Y', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_Y, 'y'))
	incrementCounter(CreateHookKey(VK_Z, 'Z', IsCapitalLetter()))
	incrementCounter(CreateHookKey(VK_Z, 'z'))
	incrementCounter(CreateHookKey(VK_NUMPAD0, '0'))
	incrementCounter(CreateHookKey(VK_NUMPAD1, '1'))
	incrementCounter(CreateHookKey(VK_NUMPAD2, '2'))
	incrementCounter(CreateHookKey(VK_NUMPAD3, '3'))
	incrementCounter(CreateHookKey(VK_NUMPAD4, '4'))
	incrementCounter(CreateHookKey(VK_NUMPAD5, '5'))
	incrementCounter(CreateHookKey(VK_NUMPAD6, '6'))
	incrementCounter(CreateHookKey(VK_NUMPAD7, '7'))
	incrementCounter(CreateHookKey(VK_NUMPAD8, '8'))
	incrementCounter(CreateHookKey(VK_NUMPAD9, '9'))
	incrementCounter(CreateHookKey(VK_MULTIPLY, '*'))
	incrementCounter(CreateHookKey(VK_ADD, '+'))
	incrementCounter(CreateHookKey(VK_SUBTRACT, '-'))
	incrementCounter(CreateHookKey(VK_DECIMAL, '.'))
	incrementCounter(CreateHookKey(VK_DIVIDE, '/'))
	incrementCounter(CreateHookKey(VK_LSHIFT, invalidRune))
	incrementCounter(CreateHookKey(VK_RSHIFT, invalidRune))
	incrementCounter(CreateHookKey(VK_LCONTROL, invalidRune))
	incrementCounter(CreateHookKey(VK_RCONTROL, invalidRune))
	incrementCounter(CreateHookKey(VK_OEM_1, ';'))
	incrementCounter(CreateHookKey(VK_OEM_1, ':', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_PLUS, '='))
	incrementCounter(CreateHookKey(VK_OEM_PLUS, '+', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_COMMA, ','))
	incrementCounter(CreateHookKey(VK_OEM_COMMA, '<', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_MINUS, '-'))
	incrementCounter(CreateHookKey(VK_OEM_MINUS, '_', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_PERIOD, '.'))
	incrementCounter(CreateHookKey(VK_OEM_PERIOD, '>', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_2, '/'))
	incrementCounter(CreateHookKey(VK_OEM_2, '?', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_3, '`'))
	incrementCounter(CreateHookKey(VK_OEM_3, '~', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_4, '['))
	incrementCounter(CreateHookKey(VK_OEM_4, '{', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_5, '\\'))
	incrementCounter(CreateHookKey(VK_OEM_5, '|', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_6, ']'))
	incrementCounter(CreateHookKey(VK_OEM_6, '}', IsShiftNeeded()))
	incrementCounter(CreateHookKey(VK_OEM_7, '\''))
	incrementCounter(CreateHookKey(VK_OEM_7, '"', IsShiftNeeded()))
}
