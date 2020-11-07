package apki

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// https://mirrors.alpinelinux.org/mirrors.json
type Mirror struct {
	gorm.Model
	Name                string `gorm:"unique"`
	Location            string
	Bandwidth           string
	URL                 string
	URLs                pgtype.TextArray `gorm:"type:text[]"`
	LastUpdated         time.Time
	LastError           string
	LastRefreshDuration time.Duration
}

type PackageIndex struct {
	gorm.Model
	Branch      string `grom:"index"`
	Repo        string `grom:"index"`
	Arch        string `grom:"index"`
	Name        string `grom:"index"`
	Version     string
	Size        int
	InstallSize int
	Description string
	URL         string
	License     string `grom:"index"`
	Maintainer  string `grom:"index"`
	Origin      string `grom:"index"`
	BuildTime   time.Time
	Commit      string
	Key         string           `gorm:"uniqueIndex"` // $BRANCH/$REPO/$ARCH/$NAME
	Path        string           `gorm:"uniqueIndex"`
	Depends     pgtype.TextArray `gorm:"type:text[]"`
	Provides    pgtype.TextArray `gorm:"type:text[]"`
	InstallIf   pgtype.TextArray `gorm:"type:text[]"`
}

func (p *PackageIndex) BeforeSave(tx *gorm.DB) (err error) {
	if anyEmpty(p.Branch, p.Repo, p.Arch, p.Name, p.Version) {
		return errors.New("invalid package index: contain empty field")
	}
	p.Key = fmt.Sprintf("%v/%v/%v/%v", p.Branch, p.Repo, p.Arch, p.Name)
	p.Path = fmt.Sprintf("%v-%v.apk", p.Key, p.Version)
	return nil
}

func anyEmpty(a ...string) bool {
	for _, v := range a {
		if v == "" {
			return true
		}
	}
	return false
}

type Setting struct {
	gorm.Model
	Name    string `gorm:"unique"`
	Value   datatypes.JSON
	Version string
}
