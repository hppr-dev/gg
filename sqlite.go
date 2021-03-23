package gg

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

// SQLiteConfig represents an SQLite database configuration
type SQLiteConfig struct {
  File string
}

func (cfg SQLiteConfig) Configure() gorm.Dialector {
  return sqlite.Open(cfg.GetDSN())
}

func (cfg SQLiteConfig) GetDSN() string {
  return cfg.File
}
