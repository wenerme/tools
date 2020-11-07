package apki

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConf struct {
	Type        string
	URL         string
	AutoMigrate bool
}

func connectDatabase(conf DatabaseConf) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	switch conf.Type {
	case "postgresql":
		fallthrough
	case "postgres":
		db, err = gorm.Open(postgres.Open(conf.URL), &gorm.Config{})
	default:
		err = fmt.Errorf("unsupported db type: %q", conf.Type)
	}
	return db, err
}
