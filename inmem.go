package gg

import (
  "gorm.io/gorm"
)

type InMemConfig SQLiteConfig

const inMemDSN = "file::memory:?cache=shared"

func (InMemConfig) Configure() gorm.Dialector {
  return  SQLiteConfig{inMemDSN}.Configure()
}

func (InMemConfig) GetDSN() string{
  return inMemDSN
}
