package gg

import (
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

// PostgresConfig represents a Postgres Database configuration 
type PostgresConfig struct {
  Host string
  Port int
  User string
  DatabaseName string
  Password string
  SSL bool
  TimeZone string
}

func (cfg PostgresConfig) Configure() gorm.Dialector {
  return postgres.Open(cfg.GetDSN())
}

func (cfg *PostgresConfig) setDefaults() {
  if cfg.Host == "" {
    cfg.Host = "localhost"
  }
  if cfg.Port == 0 {
    cfg.Port = 5432
  }
  if cfg.TimeZone == "" {
    cfg.TimeZone = "UTC"
  }
}

func (cfg PostgresConfig) GetDSN() string {
  sslMode := "disable"
  if cfg.SSL {
    sslMode = "enable"
  }
  cfg.setDefaults()
  return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
    cfg.Host,
    cfg.User,
    cfg.Password,
    cfg.DatabaseName,
    cfg.Port,
    sslMode,
    cfg.TimeZone,
  )
}
