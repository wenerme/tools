package models

import (
	"time"

	"github.com/lib/pq"
)

// https://mirrors.alpinelinux.org/mirrors.json
type Mirror struct {
	Model
	Name                string
	Location            string
	Bandwidth           string
	Host                string `gorm:"unique"`
	URL                 string
	URLs                pq.StringArray `gorm:"type:text[]"`
	LastUpdated         time.Time
	LastError           string
	LastRefreshDuration time.Duration
}
