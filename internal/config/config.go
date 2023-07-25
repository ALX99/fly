package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Config represents all of the configuration options
type Config struct {
	Settings Settings `yaml:"settings"`
}

type Settings struct {
	ScrollPadding   int      `yaml:"scrollPadding"`
	ShowHiddenFiles bool     `yaml:"showHiddenFiles"`
	Keybinds        Keybinds `yaml:"keybinds"`
}
type Keybinds struct {
	NavUp    string `yaml:"up"`
	NavDown  string `yaml:"down"`
	NavLeft  string `yaml:"left"`
	NavRight string `yaml:"right"`
	Delete   string `yaml:"delete"`
}

// GetConfig reads, pareses and returns the configuration
func GetConfig() (Config, error) {
	cfgLoc := getCfgFileLoc()

	f, err := os.ReadFile(cfgLoc)
	if err != nil {
		if os.IsNotExist(err) {
			// Safe default if no config found
			return Config{
				Settings: Settings{
					ScrollPadding: 2,
					Keybinds: Keybinds{
						NavUp:    "k",
						NavDown:  "j",
						NavLeft:  "k",
						NavRight: "l",
						Delete:   "d",
					},
				},
			}, nil
		}
		return Config{}, err
	}

	cfg := Config{}
	if err = yaml.Unmarshal(f, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// getCfgFileLoc returns the configuration file location
func getCfgFileLoc() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir != "" {
		return path.Join(path.Clean(dir), "fly") + "/config.yaml"
	}

	return "~/.config/fly/config.yaml"
}
