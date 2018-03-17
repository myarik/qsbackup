package config

// Config contains configuration information to do a backup
type BackupConfig struct {
	Name string
	Description string `yaml:"description,omitempty"`
	Logfile string `yaml:"logfile,omitempty"`
	Dirs []Dir `yaml:"dirs,omitempty"`
}

type Dir struct {
	Name string
	Description string `yaml:"description,omitempty"`
	Path string
}
