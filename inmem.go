package gg

import (
  "gorm.io/gorm"
)

// InMemConfig is a database configuration for an in memory sqlite database
// This is the same as creating an sqlite config with the DSN of "file::memory:?cache=shared"
type InMemConfig SQLiteConfig

func (i InMemConfig) Configure() gorm.Dialector {
  return  SQLiteConfig{i.GetDSN()}.Configure()
}

func (InMemConfig) GetDSN() string{
  return "file::memory:?cache=shared"
}
