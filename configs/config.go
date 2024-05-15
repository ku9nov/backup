package configs

type Config struct {
	Default struct {
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
			RetentionPeriodDaily   int  `yaml:"retentionPeriodDaily"`
			RetentionPeriodWeekly  int  `yaml:"retentionPeriodWeekly"`
			RetentionPeriodMonthly int  `yaml:"retentionPeriodMonthly"`
		} `yaml:"retention"`
	} `yaml:"default"`
	Minio struct {
		Secure     bool   `yaml:"secure"`
		S3Endpoint string `yaml:"s3Endpoint"`
	} `yaml:"minio"`
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
	} `yaml:"postgresql"`
	Additional struct {
		Enabled bool     `yaml:"enabled"`
		Files   []string `yaml:"files"`
	} `yaml:"additional"`
}
