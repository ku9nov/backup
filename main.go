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
	logrus.Infof("Retention enabled: %t, retention period: %s", configs.Default.Retention.Enabled, configs.Default.Retention.RetentionPeriod)
	awsCfg := utils.AWSAuth(*configs.Config)
	utils.ConnectToS3(awsCfg, *configs.Config)
	// MongoDB configuration
	if configs.Mongo.Enabled {
		backup.CreateMongoBackup(configs.Mongo.Host, configs.Mongo.Port, configs.Mongo.Auth.Enabled, configs.Mongo.Auth.Username, configs.Mongo.Auth.Password)
	}
	// // MySQL configuration
	if configs.MySQL.Enabled {
		backup.CreateMySQLBackup(configs.MySQL.Host, configs.MySQL.Port, configs.MySQL.Auth.Enabled, configs.MySQL.Auth.Username, configs.MySQL.Auth.Password)
	}
	// // PostgreSQL configuration
	if configs.PostgreSQL.Enabled {
		backup.CreatePostgreSQLBackup(configs.PostgreSQL.Host, configs.PostgreSQL.Port, configs.PostgreSQL.Auth.Enabled, configs.PostgreSQL.Auth.Username, configs.PostgreSQL.Auth.Password)
	}
	// // Additional configurations
	if configs.Additional.Enabled {
		backup.CreateAdditionalFilesBackup(configs.Additional.Files, currentDate)
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
