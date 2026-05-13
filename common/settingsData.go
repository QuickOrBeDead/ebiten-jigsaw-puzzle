package common

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SettingsData struct {
	SFXVolume   float64 `json:"sfx_volume"`
	MusicVolume float64 `json:"music_volume"`
}

var (
	settings     *SettingsData
	settingsPath string
)

func init() {
	settings = &SettingsData{
		SFXVolume:   0.7,
		MusicVolume: 0.5,
	}
	dir, err := os.UserConfigDir()
	if err == nil {
		settingsPath = filepath.Join(dir, "ebiten-jigsaw-puzzle", "settings.json")
	}
	loadSettings()
}

func loadSettings() {
	if settingsPath == "" {
		return
	}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return
	}
	var loaded SettingsData
	if err := json.Unmarshal(data, &loaded); err != nil {
		return
	}
	settings = &loaded
}

func SaveSettings() {
	if settingsPath == "" {
		return
	}
	dir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(settingsPath, data, 0644)
}

func GetSettings() *SettingsData {
	return settings
}
