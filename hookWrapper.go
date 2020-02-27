package main

import (
	"fmt"
	"strings"

	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

const (
	windowsNewLine = "\r\n"
)

// Create []TagInputs that can be used in SendInput() function
func createTagInputs(strToSend string, isShiftEnabled, isCapsEnabled bool) (tagInputs []TagINPUT) {

	// Store if character in iteration is SHIFT
	var isShiftNeeded bool

	for _, c := range strToSend {

		// Store the current tag input
		currentStrTag := tagInputKeyboard()
		currentStrTagUp := tagInputKeyboard()

		// Get current tag
		currentStrTag.Ki.WVk, _, isShiftNeeded = findAllKeyCode(0, c)

		// Up key
		currentStrTagUp.Ki.WVk, _, isShiftNeeded = findAllKeyCode(0, c)
		currentStrTagUp.Ki.DwFlags = KEYEVENTF_KEYUP

		// Temporary remove caps lock state if caps lock is on
		if isCapsEnabled {
			tagInputs = append(tagInputs, tagInputCapsDown())
		}

		if isShiftNeeded && !isShiftEnabled {
			tagInputs = append(tagInputs, tagInputShiftDown(), currentStrTag, tagInputShiftUp())
		} else if currentStrTag.Ki.WVk != 0 {
			tagInputs = append(tagInputs, currentStrTag, currentStrTagUp)
		}

		// Restore caps lock state
		if isCapsEnabled {
			tagInputs = append(tagInputs, tagInputCapsDown())
		}
	}

	return
}

// Find listed key code from the given character.
//
// If sending text to keyboard,
// fill the 2nd parameter, and use only 1st and 3rd return values.
//
// If receiving text from keyboard, fill 1st, 3rd and 4th paramters
// (3rd param: GetKeyState SHIFT_KEY, 4th param: GetKeyState CAPS_LOCK),
// and use only 2nd return value.
//
// Default values for input and return values:
// keyCode: 0
// char: -1
//
// Paramters: Key code, character, is shift enabled, is caps lock enabled.
//
// Return values: Key code, character, is shift needed
func findAllKeyCode(k uint16, c rune, isCapitalEnabled ...bool) (keyCode WORD, char rune, isShiftNeeded bool) {
	paramBoolLen := len(isCapitalEnabled)

	if paramBoolLen == 2 && k == 0 || paramBoolLen == 0 && c == -1 {
		panic("Wrong parameter for function findAllKeyCode")
	}

	// Receiving text
	if paramBoolLen > 0 {
		isCapitalLetter := IsCapitalLetterEnabled(isCapitalEnabled[0], isCapitalEnabled[1])

		if isCapitalEnabled[0] {
			keyCode, char = findShiftKeyCode(k, c, isCapitalLetter)
		} else {
			keyCode, char = findNonShiftKeyCode(k, c, isCapitalLetter)
		}
	} else {
		// Sending text
		keyCode, char = findShiftKeyCode(k, c)

		if isShiftNeeded = keyCode != 0; isShiftNeeded {
			return
		}

		keyCode, char = findNonShiftKeyCode(k, c)
	}

	return
}

// SHIFT needed
func findShiftKeyCode(k uint16, c rune, isCapitalLetter ...bool) (keyCode WORD, char rune) {
	if k != 0 && len(isCapitalLetter) == 0 {
		panic("No capital letter state received when receiving keystroke")
	}

	// Use ASCII value to identify character
	switch {
	case k == VK_ONE, c == 33: // Exclamation mark !
		keyCode = VK_ONE
		char = 33

	case k == VK_OEM_7, c == 34: // Double quote
		keyCode = VK_OEM_7
		char = 34

	case k == VK_THREE, c == 35: // Hash/Sharp #
		keyCode = VK_THREE
		char = 35

	case k == VK_FOUR, c == 36: // Dollar $
		keyCode = VK_FOUR
		char = 36

	case k == VK_FIVE, c == 37: // Percent %
		keyCode = VK_FIVE
		char = 37

	case k == VK_SEVEN, c == 38: // Ampersand &
		keyCode = VK_SEVEN
		char = 38

	case k == VK_NINE, c == 40: // Left parenthesis (
		keyCode = VK_NINE
		char = 40

	case k == VK_ZERO, c == 41: // Right parenthesis )
		keyCode = VK_ZERO
		char = 41

	case k == VK_EIGHT, c == 42: // Asterisk *
		keyCode = VK_EIGHT
		char = 42

	case k == VK_OEM_PLUS, c == 43: // Plus +
		keyCode = VK_OEM_PLUS
		char = 43

	case k == VK_OEM_1, c == 58: // Colon :
		keyCode = VK_OEM_1
		char = 58

	case k == VK_OEM_COMMA, c == 60: // Left angled bracket <
		keyCode = VK_OEM_COMMA
		char = 60

	case k == VK_OEM_PERIOD, c == 62: // Right angled bracket <
		keyCode = VK_OEM_PERIOD
		char = 62

	case k == VK_OEM_2, c == 63: // Question mark ?
		keyCode = VK_OEM_2
		char = 63

	case k == VK_TWO, c == 64: // At @
		keyCode = VK_TWO
		char = 64

	case VK_A <= k && k <= VK_Z, 65 <= c && c <= 90: // Capital letter: A-Z
		keyCode = WORD(c)
		char = rune(k)

		// If shift pressed and caps lock toggled
		if len(isCapitalLetter) > 0 && !isCapitalLetter[0] {
			char = rune(strings.ToLower(string(k))[0])
		}

	case k == VK_SIX, c == 94: // Caret ^
		keyCode = VK_SIX
		char = 94

	case k == VK_OEM_MINUS, c == 95: // Underscore _
		keyCode = VK_OEM_MINUS
		char = 95

	case k == VK_OEM_4, c == 123: // Left brace {
		keyCode = VK_OEM_4
		char = 123

	case k == VK_OEM_5, c == 124: // Vertical bar/Pipe |
		keyCode = VK_OEM_5
		char = 124

	case k == VK_OEM_6, c == 125: // Right brace }
		keyCode = VK_OEM_6
		char = 125

	case k == VK_OEM_3, c == 126: // Tilde ~
		keyCode = VK_OEM_3
		char = 126

	default:
		// Don't process key if not specified above
		// Or keys like backspace, delete, and weird symbols will be added to the buffer
		return 0, -1
	}

	return
}

// SHIFT not needed
func findNonShiftKeyCode(k uint16, c rune, isCapitalLetter ...bool) (keyCode WORD, char rune) {
	if k != 0 && len(isCapitalLetter) == 0 {
		panic("No capital letter state received when receiving keystroke")
	}

	// Use ASCII value to identify character
	switch {
	case k == VK_BACK, c == 8, // Backspace '\b'
		k == VK_TAB, c == 9, // horizontal tab '\t'
		k == VK_SPACE, c == 32, // spacebar
		k == VK_RETURN, c == 13, // CRLF - only carriage return for ENTER key '\r'
		VK_ZERO <= k && k <= VK_NINE, 48 <= c && c <= 57: // 0-9
		keyCode = WORD(c)
		char = rune(k)

	case k == VK_OEM_7, c == 39: // Single quote
		keyCode = VK_OEM_7
		char = 39

	case k == VK_OEM_COMMA, c == 44: // Comma
		keyCode = VK_OEM_COMMA
		char = 44

	case k == VK_OEM_MINUS, c == 45: // Hypen or minus
		keyCode = VK_OEM_MINUS
		char = 45

	case k == VK_OEM_PERIOD, c == 46: // Period
		keyCode = VK_OEM_PERIOD
		char = 46

	case k == VK_OEM_2, c == 47: // Slash or divide '/'
		keyCode = VK_OEM_2
		char = 47

	case k == VK_OEM_1, c == 59: // Semicolon
		keyCode = VK_OEM_1
		char = 59

	case k == VK_OEM_PLUS, c == 61: // Equal
		keyCode = VK_OEM_PLUS
		char = 61

	case k == VK_OEM_4, c == 91: // Opening square bracket
		keyCode = VK_OEM_4
		char = 91

	case k == VK_OEM_5, c == 92: // Backslash '\'
		keyCode = VK_OEM_5
		char = 92

	case k == VK_OEM_6, c == 93: // Closing square bracket
		keyCode = VK_OEM_6
		char = 93

	case k == VK_OEM_3, c == 96: // Grave accent '`'
		keyCode = VK_OEM_3
		char = 96

	case VK_A <= k && k <= VK_Z, 97 <= c && c <= 122: // Small letters: a-z
		keyCode = WORD(strings.ToUpper(string(c))[0]) // Keyboard code equals capital letter value
		char = rune(strings.ToLower(string(k))[0])

		// If shift no pressed but caps lock toggled
		if len(isCapitalLetter) > 0 && isCapitalLetter[0] {
			char = rune(k)
		}

	case k == VK_DELETE, c == 127: // Grave accent '`'
		keyCode = VK_DELETE
		char = 127

	default:
		// Don't process key if not specified above
		// Or keys like backspace, delete, and weird symbols will be added to the buffer
		return 0, -1
	}

	return
}

// Ex: To insert left arrow 5 times, do:
// multiplyTagInputKey(tagInputLeftArrowDown(), 5)
//
func multiplyTagInputKey(tagInputs []TagINPUT, multiplier int) []TagINPUT {
	tagInputsLength := len(tagInputs)

	newTagInputs := make([]TagINPUT, multiplier*tagInputsLength)

	if tagInputsLength == 0 {
		return newTagInputs
	}

	i := 0

	for k := range newTagInputs {
		fmt.Println(i, tagInputsLength, multiplier, len(newTagInputs))
		fmt.Println(tagInputs[i])
		newTagInputs[k] = tagInputs[i]

		if i+1 < tagInputsLength {
			i++
		} else {
			i = 0
		}
	}

	return newTagInputs
}

// Process the return value from GetKeyState.
//
// If the key state needed is not key down/up but key toggled like caps lock,
// then pass in true in the second bool parameter
func getKeyStateBool(state SHORT, checkToggle ...bool) bool {
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
