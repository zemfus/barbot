package config

// TelegramConfig represents the Telegram-related configuration
type TelegramConfig struct {
	Token string `yaml:"token"`
	Debug bool   `yaml:"debug"`
}

// DatabaseConfig represents the database-related configuration
type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	UserName       string `yaml:"userName"`
	Password       string `yaml:"password"`
	DbName         string `yaml:"dbName"`
	MigrationsPath string `yaml:"migrationsPath" default:"migrations/"`
}

// AppConfig represents the application-related configuration
type AppConfig struct {
	LogLevel string `yaml:"logLevel"`
	AdminID  int64  `yaml:"adminID"`
	BarmenID int64  `yaml:"barmenID"`
	ChatID   int64  `yaml:"chatID"`
}

// Configuration is the main configuration structure that includes all other config structs
type Configuration struct {
	Telegram *TelegramConfig `yaml:"telegram"`
	Database *DatabaseConfig `yaml:"database"`
	App      *AppConfig      `yaml:"app"`
}
