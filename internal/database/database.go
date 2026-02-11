package database

import (
	"fmt"
	"os"
	"path/filepath"

	"robot_scheduler/internal/config"
	"robot_scheduler/internal/logger"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	var dialector gorm.Dialector
	var err error

	switch cfg.Database.Default {
	case "postgres":
		dialector, err = initPostgres(cfg.Database.Postgres)
	case "sqlite":
		dialector, err = initSqlite(cfg.Database.Sqlite)
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Default)
	}

	if err != nil {
		return err
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	}

	// 连接数据库
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取通用数据库对象
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池
	if cfg.Database.Default == "postgres" && cfg.Database.Postgres != nil {
		sqlDB.SetMaxIdleConns(cfg.Database.Postgres.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.Database.Postgres.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.Postgres.ConnMaxLifetime) * time.Second)
	}

	logger.Info(fmt.Sprintf("connected to %s database", cfg.Database.Default))
	return nil
}

func initPostgres(cfg *config.PostgresConfig) (gorm.Dialector, error) {
	if cfg == nil {
		return nil, fmt.Errorf("postgres config is nil")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	return postgres.Open(dsn), nil
}

func initSqlite(cfg *config.SqliteConfig) (gorm.Dialector, error) {
	if cfg == nil {
		return nil, fmt.Errorf("sqlite config is nil")
	}

	// 创建数据目录
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	config := &sqlite.Config{
		DriverName: "sqlite3",
		DSN:        cfg.Path + "?" + buildSqliteParams(cfg),
	}

	return sqlite.Open(config.DSN), nil
}

func buildSqliteParams(cfg *config.SqliteConfig) string {
	params := ""
	if cfg.BusyTimeout > 0 {
		params += fmt.Sprintf("_busy_timeout=%d&", cfg.BusyTimeout)
	}
	if !cfg.ForeignKeys {
		params += "_foreign_keys=off&"
	}
	if cfg.JournalMode != "" {
		params += fmt.Sprintf("_journal_mode=%s&", cfg.JournalMode)
	}
	return params
}
