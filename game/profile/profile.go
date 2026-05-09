package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	AppDirName         = "snake_golang"
	ConfigFileName     = "config.json"
	MaxPlayerNameRunes = 20
)

type Config struct {
	PlayerName string `json:"player_name"`
}

var (
	ErrEmptyName   = errors.New("empty name")
	ErrNameTooLong = errors.New("name too long")
	ErrBadName     = errors.New("name contains unsupported characters")
)

func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(dir, AppDirName, ConfigFileName), nil
}

func NormalizePlayerName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrEmptyName
	}
	runes := []rune(name)
	if len(runes) > MaxPlayerNameRunes {
		return "", ErrNameTooLong
	}
	for _, r := range runes {
		if unicode.IsControl(r) {
			return "", ErrBadName
		}
	}
	return name, nil
}

func LoadConfig() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, fmt.Errorf("get config path: %w", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("read config file: %w", err)
	}
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}
	config.PlayerName, err = NormalizePlayerName(config.PlayerName)
	if err != nil {
		return Config{}, fmt.Errorf("normalize player name: %w", err)
	}
	return config, nil
}

func SaveConfig(config Config) error {
	name, err := NormalizePlayerName(config.PlayerName)
	if err != nil {
		return fmt.Errorf("normalize player name: %w", err)
	}
	config.PlayerName = name

	path, err := ConfigPath()
	if err != nil {
		return fmt.Errorf("get config path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0600)
}
