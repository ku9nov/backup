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
	notify "github.com/ku9nov/backup/utils/notifications"
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
	s3Cfg, extraS3Cfg := utils.SetStorageClient(*configs.Config)
	var failedBackups []string
	// MongoDB configuration
	if configs.Mongo.Enabled {
		if !backup.CreateMongoBackup(*configs.Config, currentDate, s3Cfg, extraS3Cfg) {
			failedBackups = append(failedBackups, "MongoDB")
		}
	}
	// MySQL configuration
	if configs.MySQL.Enabled {
		if !backup.CreateMySQLBackup(*configs.Config, currentDate, s3Cfg, extraS3Cfg) {
			failedBackups = append(failedBackups, "MySQL")
		}
	}
	// PostgreSQL configuration
	if configs.PostgreSQL.Enabled {
		if !backup.CreatePostgreSQLBackup(*configs.Config, currentDate, s3Cfg, extraS3Cfg) {
			failedBackups = append(failedBackups, "PostgreSQL")
		}
	}
	// Additional configurations
	if configs.Additional.Enabled {
		if !backup.CreateAdditionalFilesBackup(*configs.Config, currentDate, s3Cfg, extraS3Cfg) {
			failedBackups = append(failedBackups, "Additional Files")
		}
	}
	if len(failedBackups) == 0 && configs.Default.Retention.Enabled {
		utils.CheckOldFilesInS3(*configs.Config, s3Cfg, false)
	}
	if len(failedBackups) == 0 && configs.ExtraBackups.Enabled && configs.ExtraBackups.Retention.Enabled {
		utils.CheckOldFilesInS3(*configs.Config, extraS3Cfg, true)
	}
	if len(failedBackups) > 0 && configs.Slack.Enabled {
		notify.SendMessageToSlack(*configs.Config, failedBackups)
	}
	if len(failedBackups) > 0 && configs.Zabbix.Enabled {
		notify.ZabbixSender(*configs.Config)
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
