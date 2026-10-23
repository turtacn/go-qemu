package vmconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFromHCI loads a VMConfig from a YAML configuration file.
func LoadFromHCI(path string) (*VMConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := &VMConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", path, err)
	}

	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config in %s: %w", path, err)
	}

	return cfg, nil
}
