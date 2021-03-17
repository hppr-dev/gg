package gg

import (
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

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

func (cfg PostgresConfig) GetDSN() string {
  sslMode := "disable"
  if cfg.SSL {
    sslMode = "enable"
  }
  return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
    cfg.Host,
    cfg.User,
    cfg.Password,
    cfg.DatabaseName,
    cfg.Port,
    sslMode,
    cfg.TimeZone,
  )
}
