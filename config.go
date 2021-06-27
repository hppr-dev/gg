package gg

import (
	"gorm.io/gorm"
	"hppr.dev/gg/database"
  "hppr.dev/gg/mdl"
)

// Gingorm Configuration
type Config struct {
	// Gorm configuration
	Gorm *gorm.Config
	//Database configuration
	Database DatabaseConfigurer
	//On database open callback - called on the db anytime it is opened
	OnDBOpen database.DatabaseCallback
	//Model References
	Models []interface{}
	// Output Format
	DefaultOutputFormat OutputFormat
}

type DatabaseConfigurer interface {
	Configure() gorm.Dialector
	GetDSN() string
}

// Utility method to open the gorm database from config
func (cfg Config) OpenDB() (*gorm.DB, error) {
	return database.Open(cfg.Database.Configure(), cfg.OnDBOpen, cfg.Gorm)
}

func (cfg Config) GetModelColumnInfo(model interface{}) mdl.ModelColumnInfo {
  db, _ := cfg.OpenDB()
  return mdl.GetModelColumnInfo(model, db)
}
