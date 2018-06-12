package main

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"os"
	"github.com/myarik/qsbackup/app"
	"io/ioutil"
	"fmt"
	"flag"
	"github.com/myarik/qsbackup/app/engine"
	"github.com/coreos/bbolt"
	"time"
	"github.com/myarik/qsbackup/app/store"
)

type Options struct {
	ConfigFile     string
	Debug, Version bool
}

var (
	cmdOptions Options
)

const defaultVersion = "0.0.1"

func init() {
	flag.StringVar(&cmdOptions.ConfigFile, "config", "", "Path to the config file")
	flag.StringVar(&cmdOptions.ConfigFile, "c", "", "Path to the config file")

	flag.BoolVar(&cmdOptions.Version, "v", false, "Show the program version")
	flag.BoolVar(&cmdOptions.Debug, "debug", false, "Debug mode")
}

func main() {
	// Parse the CommandLine
	flag.Parse()
	if cmdOptions.Version == true {
		fmt.Printf("Version %s\n", defaultVersion)
		os.Exit(0)
	}
	if cmdOptions.ConfigFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read and validate a config file
	source, err := ioutil.ReadFile(cmdOptions.ConfigFile)
	if err != nil {
		fmt.Printf("Can't open configuration file: %s\n", cmdOptions.ConfigFile)
		os.Exit(1)
	}
	conf, err := app.ConfigLoad(source)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	// backup dirs
	var backupDirs []string
	for _, backupDir := range conf.Dirs {
		backupDirs = append(backupDirs, backupDir.Path)
	}

	// setup storage
	var storage engine.Storage
	switch conf.Storage.Type {
	case "local":
		storage = &engine.LocalStorage{
			Archiver: engine.ZIP,
			DestPath: conf.Storage.DestPath,
		}
	case "aws":
		storage = &engine.AwsStorage{
			Archiver:     engine.ZIP,
			Region:       conf.Storage.AwsRegion,
			AccessKeyID:  conf.Storage.AwsKey,
			AccessSecret: conf.Storage.AwsSecret,
			Bucket:       conf.Storage.AwsBucket,
		}
	default:
		fmt.Println("storage type does not support")
		os.Exit(1)
	}

	// Setup logger
	logger, err := log.Init(conf.Logfile, cmdOptions.Debug)
	defer logger.Close()
	if err != nil {
		fmt.Printf("can't setup a logger: %s\n", err)
		os.Exit(1)
	}
	// setup boltDb
	os.MkdirAll(conf.Home, 0755)
	db, err := store.NewBoltDB(
		conf.GetDatabasePath(),
		bolt.Options{Timeout: 30 * time.Second},
		logger,
	)
	if err != nil {
		fmt.Printf("can't create a database: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	backup := &app.Backup{
		Logger:     logger,
		BackupDirs: backupDirs,
		Storage:    storage,
		DB:         db,
	}
	backup.Run(conf.Jobs)
}
