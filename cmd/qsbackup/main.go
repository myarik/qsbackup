package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/myarik/qsbackup/app"
	"github.com/myarik/qsbackup/app/engine"
	"github.com/myarik/qsbackup/app/store"
	"github.com/coreos/bbolt"
	"time"
	"github.com/myarik/qsbackup/pkg/logger"
	"os/signal"
	"syscall"
)

const defaultVersion = "0.0.3"

// Opts with command line flags and env
// nolint:maligned
type Opts struct {
	ConfigFile string `short:"c" long:"config" default:"/usr/local/etc/qsbackup.conf" description:"config file"`
	Version    func() `short:"v" long:"version" description:"version"`
	Show       bool   `short:"s" long:"show" description:"all backups"`
	Last       bool   `short:"l" long:"last" description:"last backups"`
	Debug      bool   `long:"debug" description:"debugger on"`
}

func main() {
	// Parse the CommandLine
	var opts Opts

	opts.Version = func() {
		fmt.Println(defaultVersion)
		os.Exit(0)
	}

	p := flags.NewParser(&opts, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
	}

	// Read and validate a config file
	source, err := ioutil.ReadFile(opts.ConfigFile)
	if err != nil {
		fmt.Printf("Can't open configuration file: %s\n", opts.ConfigFile)
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
			AWSS3Bucket:  conf.Storage.AwsBucket,
		}
	default:
		fmt.Println("storage type does not support")
		os.Exit(1)
	}

	// Setup logger
	appLogger, err := logger.Init(conf.Logfile, opts.Debug)
	defer appLogger.Close()
	if err != nil {
		fmt.Printf("can't setup a logger: %s\n", err)
		os.Exit(1)
	}
	// setup boltDb
	os.MkdirAll(conf.Home, 0755)
	db, err := store.NewBoltDB(
		conf.GetDatabasePath(),
		bolt.Options{Timeout: 30 * time.Second},
		appLogger,
	)
	if err != nil {
		fmt.Printf("can't create a database: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// init a backup object
	backup := &app.Backup{
		Logger:     appLogger,
		BackupDirs: backupDirs,
		Storage:    storage,
		DB:         db,
	}

	if opts.Last {
		backup.ShowLastBackups()
		os.Exit(0)
	} else if opts.Show {
		backup.AllBackups()
		os.Exit(0)
	} else {
		cancelChan := make(chan struct{}, 1)
		printChan := make(chan struct{}, 1)
		signalChan := make(chan os.Signal, 1)
		go func() {
			<-signalChan
			close(cancelChan)
			fmt.Print("Canceling .")
			for {
				select {
				case <-printChan:
					return
				default:
					time.Sleep(time.Second * 1)
					fmt.Print(".")
				}
			}
		}()
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		stopTask, err := backup.Run(conf.Jobs, cancelChan)
		if err != nil {
			os.Exit(1)
		}
		<-stopTask
		close(printChan)
	}
}
