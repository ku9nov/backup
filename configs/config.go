package configs

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
}
