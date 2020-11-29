package structs

import (
	"time"
)

type PackageIndex struct {
	Branch      string    `json:"branch"`
	Repo        string    `json:"repo"`
	Arch        string    `json:"arch"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Size        int       `json:"size"`
	InstallSize int       `json:"installSize"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	License     string    `json:"license"`
	Origin      string    `json:"origin"`
	BuildTime   time.Time `json:"buildTime"`
	Commit      string    `json:"commit"`
	Key         string    `json:"key"`
	Path        string    `json:"path"`
	Depends     []string  `json:"depends"`
	Provides    []string  `json:"provides"`
	InstallIf   []string  `json:"installIf"`

	MaintainerName string `json:"maintainerName"`

	// 额外

	CommitData *Commit `json:"commitData"`
}
