package qsbackup

import (
	"fmt"
	"github.com/go-yaml/yaml"
)

// BackupConfig represents types capable of read a config file
type BackupConfig struct {
	Name        string
	Description string `yaml:"description,omitempty"`
	Logfile     string `yaml:"logfile,omitempty"`
	Storage     BackupStorage
	Dirs        []Dir  `yaml:"dirs,omitempty"`
}

// Dir represents types capable of read a backup directory path
type Dir struct {
	Name        string
	Description string `yaml:"description,omitempty"`
	Path        string
}

// BackupStorage represents types capable of represent a backup storage
type BackupStorage struct {
	Type      string
	DestPath  string `yaml:"dest_path,omitempty"`
	AwsRegion string `yaml:"aws_region,omitempty"`
	AwsBucket string `yaml:"aws_bucket,omitempty"`
	AwsKey    string `yaml:"aws_key,omitempty"`
	AwsSecret string `yaml:"aws_secret,omitempty"`
}

func validatedConfig(c *BackupConfig) (*BackupConfig, error) {
	for _, backupDir := range c.Dirs {
		if exist, _ := IsExists(backupDir.Path); !exist {
			return nil, fmt.Errorf("the directory %s does not exist", backupDir.Path)
		}
	}
	switch c.Storage.Type {
	case "local":
		if c.Storage.DestPath == "" {
			return nil, fmt.Errorf("dest_path does not set for the local storage")
		}
	case "aws":
		if c.Storage.AwsRegion == "" {
			return nil, fmt.Errorf("aws_region does not set for the aws storage")
		}
		if c.Storage.AwsBucket == "" {
			return nil, fmt.Errorf("aws_bucket does not set for the aws storage")
		}
		if c.Storage.AwsKey == "" {
			return nil, fmt.Errorf("aws_key does not set for the aws storage")
		}
		if c.Storage.AwsSecret == "" {
			return nil, fmt.Errorf("aws_secret does not set for the aws storage")
		}

	default:
		return nil, fmt.Errorf("storage type does not support")
	}
	return c, nil
}

// ConfigLoad will unmarshal and validate a config data
func ConfigLoad(input []byte) (*BackupConfig, error) {
	var config BackupConfig

	if err := yaml.Unmarshal(input, &config); err != nil {
		return nil, fmt.Errorf("can't parse the config file: %v", err)
	}
	return validatedConfig(&config)
}
