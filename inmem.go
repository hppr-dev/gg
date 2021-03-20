package gg

import (
  "gorm.io/gorm"
)

type InMemConfig SQLiteConfig

func (i InMemConfig) Configure() gorm.Dialector {
  return  SQLiteConfig{i.GetDSN()}.Configure()
}

func (InMemConfig) GetDSN() string{
  return "file::memory:?cache=shared"
}
