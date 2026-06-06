package unlock

import (
	"os"
	"strings"
)

var unlocked bool

func GetSecretCode() string {
	return os.Getenv("SECRET_CODE")
}

func Check(code string) bool {
	secret := GetSecretCode()
	if secret == "" {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(code), secret)
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
