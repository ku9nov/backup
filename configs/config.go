package configs

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDeviceIDFile = "device_id"
	defaultStateDirName = ".backup"
)

type Config struct {
	Default struct {
		Host            string `yaml:"host"`
		StorageProvider string `yaml:"storageProvider"`
		Bucket          string `yaml:"bucket"`
		UseProfile      struct {
			Enabled bool   `yaml:"enabled"`
			Profile string `yaml:"profile"`
		} `yaml:"useProfile"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
		Region    string `yaml:"region"`
		BackupDir string `yaml:"backupDir"`
		Retention struct {
			Enabled                bool `yaml:"enabled"`
			DryRun                 bool `yaml:"dryRun"`
			RetentionPeriodDaily   int  `yaml:"retentionPeriodDaily"`
			RetentionPeriodWeekly  int  `yaml:"retentionPeriodWeekly"`
			RetentionPeriodMonthly int  `yaml:"retentionPeriodMonthly"`
		} `yaml:"retention"`
	} `yaml:"default"`
	Minio struct {
		Secure     bool   `yaml:"secure"`
		S3Endpoint string `yaml:"s3Endpoint"`
	} `yaml:"minio"`
	Azure struct {
		StorageAccountName string `yaml:"storageAccountName"`
		StorageAccountKey  string `yaml:"storageAccountKey"`
	} `yaml:"azure"`
	Mongo struct {
		Enabled   bool     `yaml:"enabled"`
		Host      string   `yaml:"host"`
		Port      string   `yaml:"port"`
		DumpTool  string   `yaml:"dumpTool"`
		Databases []string `yaml:"databases"`
		Auth      struct {
			Enabled      bool   `yaml:"enabled"`
			Username     string `yaml:"username"`
			Password     string `yaml:"password"`
			AuthDatabase string `yaml:"authDatabase"`
		} `yaml:"auth"`
	} `yaml:"mongo"`
	MySQL struct {
		Enabled   bool     `yaml:"enabled"`
		Host      string   `yaml:"host"`
		Port      string   `yaml:"port"`
		DumpTool  string   `yaml:"dumpTool"`
		Databases []string `yaml:"databases"`
		Auth      struct {
			Enabled  bool   `yaml:"enabled"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
		DumpFlags string `yaml:"dumpFlags"`
	} `yaml:"mysql"`
	PostgreSQL struct {
		Enabled   bool     `yaml:"enabled"`
		Host      string   `yaml:"host"`
		Port      string   `yaml:"port"`
		DumpTool  string   `yaml:"dumpTool"`
		Databases []string `yaml:"databases"`
		Auth      struct {
			Enabled  bool   `yaml:"enabled"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
		DumpFlags string `yaml:"dumpFlags"`
	} `yaml:"postgresql"`
	Redis struct {
		Enabled      bool   `yaml:"enabled"`
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		RedisCliTool string `yaml:"redisCliTool"`
		Auth         struct {
			Enabled  bool   `yaml:"enabled"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
	} `yaml:"redis"`
	Additional struct {
		Enabled bool     `yaml:"enabled"`
		Files   []string `yaml:"files"`
	} `yaml:"additional"`
	ExtraBackups struct {
		Enabled         bool   `yaml:"enabled"`
		StorageProvider string `yaml:"storageProvider"`
		Bucket          string `yaml:"bucket"`
		UseProfile      struct {
			Enabled bool   `yaml:"enabled"`
			Profile string `yaml:"profile"`
		} `yaml:"useProfile"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
		Region    string `yaml:"region"`
		Retention struct {
			Enabled              bool `yaml:"enabled"`
			DryRun               bool `yaml:"dryRun"`
			RetentionPeriodDaily int  `yaml:"retentionPeriodDaily"`
		} `yaml:"retention"`
	} `yaml:"extraBackups"`
	Slack struct {
		Enabled        bool   `yaml:"enabled"`
		SlackToken     string `yaml:"slackToken"`
		SlackChannelID string `yaml:"slackChannelID"`
	} `yaml:"slack"`
	Zabbix struct {
		Enabled     bool   `yaml:"enabled"`
		ZabbixUrl   string `yaml:"zabbixUrl"`
		ZabbixPort  int    `yaml:"zabbixPort"`
		ZabbixKey   string `yaml:"zabbixKey"`
		ZabbixValue string `yaml:"zabbixValue"`
	} `yaml:"zabbix"`
	Upgrade struct {
		Server string `yaml:"server"`
		Owner  string `yaml:"owner"`
		TUF    bool   `yaml:"tuf"`
		App    string `yaml:"app"`
	} `yaml:"upgrade"`
}

func StateDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	homeDir = strings.TrimSpace(homeDir)
	if homeDir == "" {
		return "", errors.New("home directory is empty")
	}
	return filepath.Join(homeDir, defaultStateDirName), nil
}

func EnsureDeviceID() (string, error) {
	stateDir, err := StateDir()
	if err != nil {
		return "", err
	}

	devicePath := filepath.Join(stateDir, defaultDeviceIDFile)
	raw, err := os.ReadFile(devicePath)
	if err == nil {
		deviceID := strings.TrimSpace(string(raw))
		if deviceID != "" {
			return deviceID, nil
		}
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return "", err
	}

	deviceID, err := generateDeviceID()
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(devicePath, []byte(deviceID+"\n"), 0o600); err != nil {
		return "", err
	}

	return deviceID, nil
}

func generateDeviceID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
