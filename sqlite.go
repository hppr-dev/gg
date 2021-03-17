package gg

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

type SQLiteConfig struct {
  File string
}

func (cfg SQLiteConfig) Configure() gorm.Dialector {
  return sqlite.Open(cfg.GetDSN())
}

func (cfg SQLiteConfig) GetDSN() string {
  return cfg.File
}
