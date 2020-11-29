package models

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PackageIndex struct {
	gorm.Model
	Branch      string `gorm:"index"`
	Repo        string `gorm:"index"`
	Arch        string `gorm:"index"`
	Name        string `gorm:"index"`
	Version     string
	Size        int
	InstallSize int
	Description string
	URL         string
	License     string `gorm:"index"`
	Maintainer  string
	Origin      string `gorm:"index"`
	BuildTime   time.Time
	Commit      string
	Depends     pq.StringArray `gorm:"type:text[]"`
	Provides    pq.StringArray `gorm:"type:text[]"`
	InstallIf   pq.StringArray `gorm:"type:text[]"`

	// derived
	MaintainerName  string
	MaintainerEmail string
	Path            string `gorm:"uniqueIndex"`
	Key             string `gorm:"uniqueIndex"` // $BRANCH/$REPO/$ARCH/$NAME
}

func (p *PackageIndex) BeforeSave(tx *gorm.DB) (err error) {
	if anyEmpty(p.Branch, p.Repo, p.Arch, p.Name, p.Version) {
		return errors.New("invalid package index: contain empty field")
	}
	p.Key = fmt.Sprintf("%v/%v/%v/%v", p.Branch, p.Repo, p.Arch, p.Name)
	p.Path = fmt.Sprintf("%v-%v.apk", p.Key, p.Version)
	if p.Maintainer != "" {
		addr, err := mail.ParseAddress(p.Maintainer)
		if err == nil {
			p.MaintainerName = addr.Name
			p.MaintainerEmail = addr.Address
		}
	}
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
