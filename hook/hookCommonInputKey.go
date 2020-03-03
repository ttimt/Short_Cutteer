package hook

import (
	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

// Key stores all the possible keys
type Key struct {
	KeyCode         uint16
	Char            rune
	IsShiftNeeded   bool
	IsCapitalLetter bool
}

// Option type for constructing a key
type Option func(key *Key)

// Create base tag input keyboard template
func tagInputKeyboard() TagINPUT {
	return TagINPUT{
		InputType: INPUT_KEYBOARD,
	}
}

// Simulate key pressed down
func (k *Key) keyDown() TagINPUT {
	tagInput := tagInputKeyboard()
	tagInput.Ki.WVk = WORD(k.KeyCode)

	return tagInput
}

// Simulate key press released
func (k *Key) keyUp() TagINPUT {
	tagInput := k.keyDown()
	tagInput.Ki.DwFlags = KEYEVENTF_KEYUP

	return tagInput
}

// Simulate a key press: Key down and up.
// With a specified number of times (multiplier[0])
func (k *Key) KeyPress(multiplier ...int) []TagINPUT {
	realMultiplier := 0

	// Check input
	if len(multiplier) == 0 {
		realMultiplier = 1
	} else if len(multiplier) > 1 {
		panic("Only 1 number for multiplier allowed")
	} else {
		realMultiplier = multiplier[0]
	}

	// Create tag inputs based on number of times to simulate
	tagInputs := make([]TagINPUT, 2*realMultiplier)

	for i := 0; i < realMultiplier; i++ {
		tagInputs[i*2] = k.keyDown()
		tagInputs[i*2+1] = k.keyUp()
	}

	return tagInputs
}

// Simulate a key hold: Key down
// Ex: CTRL key hold down
func (k *Key) KeyHold() TagINPUT {
	return k.keyDown()
}

// Simulate a key hold release: Key up
// Ex: CTRL key hold release
func (k *Key) KeyRelease() TagINPUT {
	return k.keyUp()
}

// Construct a hook key.
// First and second parameter are mandatory.
// Remaining parameters are optional, to set isShiftNeeded or isCapitalLetter
//
// Ex: CreateHookKey(0x123, 'a', hook.IsShiftNeeded())
func CreateHookKey(keyCode uint16, char rune, options ...Option) Key {
	key := Key{
		KeyCode: keyCode,
		Char:    char,
	}

	// Set the options
	for k := range options {
		options[k](&key)
	}

	return key
}

// Set is shift needed. Set to false for capital letters
func IsShiftNeeded() Option {
	return func(k *Key) {
		k.IsShiftNeeded = true
	}
}

// Set is capital letter
func IsCapitalLetter() Option {
	return func(k *Key) {
		k.IsCapitalLetter = true
	}
}
