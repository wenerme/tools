package apki

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

// https://mirrors.alpinelinux.org/mirrors.json
type Mirror struct {
	gorm.Model
	Name                string `gorm:"uniq"`
	Location            string
	Bandwidth           string
	URL                 string
	URLs                pgtype.TextArray `gorm:"type:text[]"`
	LastUpdated         time.Time
	LastError           string
	LastRefreshDuration time.Duration
}
