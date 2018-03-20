package config

import (
	"github.com/go-yaml/yaml"
	"fmt"
)

// Config contains configuration information to do a backup
type BackupConfig struct {
	Name        string
	Description string `yaml:"description,omitempty"`
	Logfile     string `yaml:"logfile,omitempty"`
	Dirs        []Dir  `yaml:"dirs,omitempty"`
}

type Dir struct {
	Name        string
	Description string `yaml:"description,omitempty"`
	Path        string
}

func Load(input []byte) (*BackupConfig, error) {
	var config BackupConfig

	if err := yaml.Unmarshal(input, &config); err != nil {
		return nil, fmt.Errorf("Can't parse the config file")
	}
	return &config, nil
}
