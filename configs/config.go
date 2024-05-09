package configs

type Config struct {
	Default struct {
		Bucket     string `yaml:"bucket"`
		UseProfile struct {
			Enabled bool   `yaml:"enabled"`
			Profile string `yaml:"profile"`
		} `yaml:"useProfile"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
		Region    string `yaml:"region"`
		Retention struct {
			Enabled         bool   `yaml:"enabled"`
			RetentionPeriod string `yaml:"retentionPeriod"`
		} `yaml:"retention"`
	} `yaml:"default"`
	Mongo struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Auth    struct {
			Enabled  bool   `yaml:"enabled"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
	} `yaml:"mongo"`
	MySQL struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Auth    struct {
			Enabled  bool   `yaml:"enabled"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
	} `yaml:"mysql"`
	PostgreSQL struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Auth    struct {
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
