//go:generate gomodifytags -file dto.go -w -all -add-tags json -transform camelcase
// go get github.com/fatih/gomodifytags
package apis

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
type PackageRepresentation struct {
	Branch      string         `json:"branch"`
	Repo        string         `json:"repo"`
	Arch        string         `json:"arch"`
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Size        int            `json:"size"`
	InstallSize int            `json:"installSize"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	License     string         `json:"license"`
	Origin      string         `json:"origin"`
	BuildTime   time.Time      `json:"buildTime"`
	Commit      string         `json:"commit"`
	Key         string         `json:"key"`
	Path        string         `json:"path"`
	Depends     pq.StringArray `gorm:"type:text[]" json:"depends"`
	Provides    pq.StringArray `gorm:"type:text[]" json:"provides"`
	InstallIf   pq.StringArray `gorm:"type:text[]" json:"installIf"`

	MaintainerName string `json:"maintainerName"`
}
