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
	Model
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

	// extra
	CommitData     *Commit         `gorm:"-"`
	DependPackages []*PackageIndex `gorm:"-"`
}

func (m *PackageIndex) BeforeSave(tx *gorm.DB) (err error) {
	if anyEmpty(m.Branch, m.Repo, m.Arch, m.Name, m.Version) {
		return errors.New("invalid package index: contain empty field")
	}
	m.Key = fmt.Sprintf("%v/%v/%v/%v", m.Branch, m.Repo, m.Arch, m.Name)
	m.Path = fmt.Sprintf("%v-%v.apk", m.Key, m.Version)
	if m.Maintainer != "" {
		addr, err := mail.ParseAddress(m.Maintainer)
		if err == nil {
			m.MaintainerName = addr.Name
			m.MaintainerEmail = addr.Address
		}
	}
	return nil
}

func (m *PackageIndex) GetCommitData() (*Commit, error) {
	if m.CommitData != nil {
		return m.CommitData, nil
	}
	if m.Commit == "" {
		return nil, nil
	}
	m.CommitData = &Commit{}
	return m.CommitData, m.DB.Find(m.CommitData, "hash = ?", m.Commit).Error
}

func anyEmpty(a ...string) bool {
	for _, v := range a {
		if v == "" {
			return true
		}
	}
	return false
}
