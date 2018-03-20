package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/myarik/qsbackup/pkg/util/file"
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

func validatedConfig(c *BackupConfig) (*BackupConfig, error) {
	for _, backupDir := range c.Dirs {
		if exist, _ := file.IsExists(backupDir.Path); !exist {
			return nil, fmt.Errorf("the directory %s does not exist\n", backupDir.Path)
		}
	}
	return c, nil
}

func Load(input []byte) (*BackupConfig, error) {
	var config BackupConfig

	if err := yaml.Unmarshal(input, &config); err != nil {
		return nil, fmt.Errorf("can't parse the config file\n")
	}
	return validatedConfig(&config)
}
