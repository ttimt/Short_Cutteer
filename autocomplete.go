package main

import (
	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

// Return if the character is eligible for an auto complete action
func isAutoComplete(char rune) bool {
	return char ==
		'(' || char == // Disable this if using parameter completion
		'[' || char ==
		'{' || char ==
		'\'' || char ==
		'"'
}

// Process auto complete
func processAutoComplete(char rune, isShiftEnabled, isCapsEnabled bool) (tagInputs []TagINPUT) {
	// Process auto complete
	switch char {
	case '(':
		// Parenthesis
		tagInputs = createTagInputs(" )", isShiftEnabled, isCapsEnabled)
		tagInputs = append(tagInputs, getKeyByKeyCode(VK_LEFT).KeyPress(2)...)
	case '{':
		// Scope body
		tagInputs = createTagInputs(windowsNewLine, isShiftEnabled, isCapsEnabled)
		tagInputs = append(tagInputs, createTagInputs("}", isShiftEnabled, isCapsEnabled)...)
		tagInputs = append(tagInputs, getKeyByKeyCode(VK_LEFT).KeyPress()...)
		tagInputs = append(tagInputs, createTagInputs(windowsNewLine, isShiftEnabled, isCapsEnabled)...)
		tagInputs = append(tagInputs, getKeyByKeyCode(VK_LEFT).KeyPress()...)
		tagInputs = append(tagInputs, createTagInputs("  ", isShiftEnabled, isCapsEnabled)...)
	case '[':
		tagInputs = createTagInputs("]", isShiftEnabled, isCapsEnabled)
		tagInputs = append(tagInputs, getKeyByKeyCode(VK_LEFT).KeyPress()...)
		autoCompleteJustDone = true
	case '\'':
		tagInputs = createTagInputs("'", isShiftEnabled, isCapsEnabled)
		tagInputs = append(tagInputs, getKeyByKeyCode(VK_LEFT).KeyPress()...)
		autoCompleteJustDone = true
	case '"':
		tagInputs = createTagInputs("\"", isShiftEnabled, isCapsEnabled)
		tagInputs = append(tagInputs, getKeyByKeyCode(VK_LEFT).KeyPress()...)
		autoCompleteJustDone = true
	default:
		panic("Auto complete not match with isAutoComplete(char) function!")
	}

	return
}
