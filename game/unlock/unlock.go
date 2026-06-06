package unlock

import "strings"

const secretCode = "qwerty"

var unlocked bool

func Check(code string) bool {
	return strings.EqualFold(strings.TrimSpace(code), secretCode)
}

func IsUnlocked() bool {
	return unlocked
}

func SetUnlocked(value bool) {
	unlocked = value
}

func TryUnlock(code string) bool {
	if Check(code) {
		SetUnlocked(true)
		return true
	}
	return false
}
