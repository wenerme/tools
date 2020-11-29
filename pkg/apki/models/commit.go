package models

import (
	"time"

	"github.com/lib/pq"
)

type Commit struct {
	Model

	Hash         string `gorm:"uniqueIndex"`
	Message      string
	TreeHash     string
	ParentHashes pq.StringArray `gorm:"index;type:text[]"`

	Author    *CommitSignature `gorm:"embedded;embeddedPrefix:author_"`
	Committer *CommitSignature `gorm:"embedded;embeddedPrefix:committer_"`
}

type CommitSignature struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}
