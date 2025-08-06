package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ChunkyConfig struct {
	Preamble  Section `toml:"preamble"`
	Postamble Section `toml:"postamble"`
}

type Section struct {
	Text string `toml:"text"`
}

// LoadConfig loads .chunkyconfig from the given repo path if it exists
func LoadConfig(repoPath string) (*ChunkyConfig, error) {
	cfgPath := filepath.Join(repoPath, ".chunkyconfig")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, nil // no config, use defaults
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var cfg ChunkyConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
