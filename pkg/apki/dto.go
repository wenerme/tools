//go:generate gomodifytags -file dto.go -w -all -add-tags json -transform camelcase
// go get github.com/fatih/gomodifytags
package apki

import (
	"time"

	"github.com/lib/pq"
)

type MirrorRepresentation struct {
	Name                string         `json:"name"`
	Location            string         `json:"location"`
	Bandwidth           string         `json:"bandwidth"`
	Host                string         `json:"host"`
	URL                 string         `json:"url"`
	URLs                pq.StringArray `json:"urls" gorm:"type:text[]"`
	LastUpdated         time.Time      `json:"lastUpdated"`
	LastError           string         `json:"lastError"`
	LastRefreshDuration time.Duration  `json:"lastRefreshDuration"`
}
