package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	App      *AppConfig      `mapstructure:"app"`
	Database *DatabaseConfig `mapstructure:"database"`
	Log      *LogConfig      `mapstructure:"log"`
	Minio    *MinioConfig    `mapstructure:"minio"`
	Platform *PlatformConfig `mapstructure:"platform"`
	Auth     *AuthConfig     `mapstructure:"auth"`
}

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Default  string          `mapstructure:"default"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
	Sqlite   *SqliteConfig   `mapstructure:"sqlite"`
}

type PostgresConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type SqliteConfig struct {
	Path        string `mapstructure:"path"`
	BusyTimeout int    `mapstructure:"busy_timeout"`
	ForeignKeys bool   `mapstructure:"foreign_keys"`
	JournalMode string `mapstructure:"journal_mode"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type MinioConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Endpoint   string `mapstructure:"endpoint"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	UseSSL     bool   `mapstructure:"use_ssl"`
	BucketName string `mapstructure:"bucket_name"`
	Region     string `mapstructure:"region"`
}

type PlatformConfig struct {
	Type string `mapstructure:"type"`
}

type AuthConfig struct {
	DESKey        string `mapstructure:"des_key"`
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpireHours int   `mapstructure:"jwt_expire_hours"`
}

var cfg *Config

func Init(configPath string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(filepath.Base(configPath))
	viper.AddConfigPath(filepath.Dir(configPath))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func Get() *Config {
	return cfg
}
