package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	backup "github.com/ku9nov/backup/backups"
	"github.com/ku9nov/backup/configs"
	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var logLevel string

var configPath string

type mainConfig struct {
	*configs.Config
}

func init() {
	flag.StringVar(&logLevel, "loglevel", "info", "log level (debug, info, warn, error, fatal, panic)")
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	logrus.New()

	// Parse log level from the command line and set it
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Println("Invalid log level specified:", err)
		os.Exit(1)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func NewConfig(configPath string) (*mainConfig, error) {
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	cfg := &configs.Config{}

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}
	mainCfg := &mainConfig{Config: cfg}

	return mainCfg, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}
	return nil
}

func (configs mainConfig) Run() {
	currentDate := time.Now().Format("20060102")
	logrus.Infof("Default bucket: %s", configs.Default.Bucket)
	logrus.Infof("Retention enabled: %t.", configs.Default.Retention.Enabled)
	s3Cfg := utils.CreateStorageClient(*configs.Config)
	allBackupsSuccessful := true
	// MongoDB configuration
	if configs.Mongo.Enabled {
		if !backup.CreateMongoBackup(*configs.Config, currentDate, s3Cfg) {
			allBackupsSuccessful = false
		}
	}
	// MySQL configuration
	if configs.MySQL.Enabled {
		if !backup.CreateMySQLBackup(*configs.Config, currentDate, s3Cfg) {
			allBackupsSuccessful = false
		}
	}
	// PostgreSQL configuration
	if configs.PostgreSQL.Enabled {
		if !backup.CreatePostgreSQLBackup(*configs.Config, currentDate, s3Cfg) {
			allBackupsSuccessful = false
		}
	}
	// Additional configurations
	if configs.Additional.Enabled {
		if !backup.CreateAdditionalFilesBackup(*configs.Config, currentDate, s3Cfg) {
			allBackupsSuccessful = false
		}
	}
	if allBackupsSuccessful && configs.Default.Retention.Enabled {
		utils.CheckOldFilesInS3(*configs.Config, s3Cfg)
	}
}

func main() {
	logrus.Info("Starting backups...")

	if err := ValidateConfigPath(configPath); err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg.Run()
	logrus.Info("Backups finished.")
}
