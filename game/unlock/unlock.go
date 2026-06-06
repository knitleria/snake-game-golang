package unlock

import (
	"os"
	"strings"
)

var unlocked bool

// BuildSecretCode can be injected at build time via -ldflags -X.
// It takes precedence over the SECRET_CODE environment variable so that
// distributed release binaries carry the code without shipping a .env file.
var BuildSecretCode string

func GetSecretCode() string {
	if BuildSecretCode != "" {
		return BuildSecretCode
	}
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
