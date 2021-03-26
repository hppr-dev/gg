package database

import (
	"gorm.io/gorm"
)

type DatabaseCallback func(*gorm.DB) *gorm.DB

func Open(dialector gorm.Dialector, onDBOpen DatabaseCallback, cfg *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, cfg)
	if onDBOpen != nil {
		onDBOpen(db)
	}
	return db, err
}
