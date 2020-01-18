package main

import (
	"strings"

	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

// Create base tag input keyboard template
func tagInputKeyboard() TagINPUT {
	return TagINPUT{
		InputType: INPUT_KEYBOARD,
	}
}

// Create tag input for SHIFT key down
func tagInputShiftDown() TagINPUT {
	tagInput := tagInputKeyboard()
	tagInput.Ki.WVk = VK_SHIFT

	return tagInput
}

// Create tag input for SHIFT key up by adding key event flag with SHIFT key down
func tagInputShiftUp() TagINPUT {
	tagInput := tagInputShiftDown()
	tagInput.Ki.DwFlags = KEYEVENTF_KEYUP

	return tagInput
}

// Create []TagInputs that can be used in SendInput() function
func createTagInputs(strToSend string) (tagInputs []TagINPUT) {

	// Store if character in iteration is SHIFT
	var isShiftNeeded bool

	for _, c := range strToSend {

		// Store the current tag input
		currentStrTag := tagInputKeyboard()

		// Get current tag
		currentStrTag.Ki.WVk, _, isShiftNeeded = findAllKeyCode(0, c)

		if isShiftNeeded {
			tagInputs = append(tagInputs, tagInputShiftDown(), currentStrTag, tagInputShiftUp())
		} else if currentStrTag.Ki.WVk != 0 {
			tagInputs = append(tagInputs, currentStrTag)
		}
	}

	return
}

// Find listed key code from the given character
// If sending text to keyboard, fill the 2nd parameter, and use only 1st and 3rd return values
// If receiving text from keyboard, fill 1st and 3rd paramters (3rd paramter: GetKeyState: SHIFT_KEY XOR CAPS_LOCK), and use only 2nd return value
// Not used parameters put 0 and _
//
// If empty, return value key code is 0 and char is -1
//
// Paramters: Key code, character, is caps lock enabled
func findAllKeyCode(k uint16, c rune, isCapsLockEnabled ...bool) (keyCode uint16, char rune, isShiftNeeded bool) {
	if len(isCapsLockEnabled) > 0 && !isCapsLockEnabled[0] {
		keyCode, char = findNonShiftKeyCode(k, c)
	} else {
		keyCode, char = findShiftKeyCode(k, c)

		if isShiftNeeded = keyCode != 0; isShiftNeeded || len(isCapsLockEnabled) > 0 {
			return
		}

		keyCode, char = findNonShiftKeyCode(k, c)
	}

	return
}

// SHIFT needed
func findShiftKeyCode(k uint16, c rune) (keyCode uint16, char rune) {
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

	case k == VK_FOUR, c == 36:
		keyCode = VK_FOUR
		char = 37

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
		keyCode = uint16(c)
		char = rune(k)

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
func findNonShiftKeyCode(k uint16, c rune) (keyCode uint16, char rune) {
	// Use ASCII value to identify character
	switch {
	case k == VK_BACK, c == 8, // Backspace
		k == VK_TAB, c == 9, // horizontal tab
		k == VK_SPACE, c == 32, // spacebar
		VK_ZERO <= k && k <= VK_NINE, 48 <= c && c <= 57: // 0-9
		keyCode = uint16(c)
		char = rune(k)

	case k == VK_RETURN, c == 10: // Line feed '\n'
		keyCode = VK_RETURN
		char = 10

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
		keyCode = uint16(strings.ToUpper(string(c))[0]) // Keyboard code equals capital letter value
		char = rune(strings.ToLower(string(k))[0])

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
