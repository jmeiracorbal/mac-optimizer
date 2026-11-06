package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DryRun  bool
	Quiet   bool
	Force   bool
	LogFile string
}

func Default() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return &Config{
		LogFile: filepath.Join(home, ".mac-optimizer.log"),
	}
}
