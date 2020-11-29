package apki

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.com/wenerme/tools/pkg/apki/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errDoneIndex = errors.New("DONE")

func IndexCommit(db *gorm.DB, dir string) error {
	batch := 100
	var rows []*models.Commit
	count := 0
	indexIter := func(it object.CommitIter) error {
		var err error
		count = 0
		err = it.ForEach(func(commit *object.Commit) error {
			m := toCommitModel(commit)
			rows = append(rows, m)
			if len(rows) != batch {
				return nil
			}
			r := db.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "hash"}},
					DoNothing: true,
				},
			).Create(&rows)
			rows = nil
			if r.Error != nil {
				return r.Error
			}
			count += int(r.RowsAffected)
			if int(r.RowsAffected) != batch {
				return errDoneIndex
			}
			return nil
		})
		if err == nil && len(rows) > 0 {
			err = db.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "hash"}},
					DoNothing: true,
				},
			).Create(&rows).Error
			rows = nil
		}
		if err == errDoneIndex {
			// done
			err = nil
		}
		return err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	//
	{
		ref, err := repo.Head()
		if err != nil {
			return err
		}
		logrus.WithField("commit", ref.Hash()).Info("start index tail commit")
		c, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		err = indexIter(object.NewCommitPreorderIter(c, map[plumbing.Hash]bool{}, nil))
		if err != nil {
			return err
		}
		logrus.WithField("count", count).Info("indexed tail commit")
	}

	{
		head := models.Commit{}
		db.Order("committer_date asc").Limit(1).First(&head)

		c, err := repo.CommitObject(plumbing.NewHash(head.Hash))
		if err != nil {
			return err
		}
		it := object.NewCommitPreorderIter(c, map[plumbing.Hash]bool{}, nil)
		// skip this one
		_, _ = it.Next()

		logrus.WithField("commit", head.Hash).Info("start index head commit")

		if err := indexIter(it); err != nil {
			return err
		}
		logrus.WithField("count", count).Info("indexed head commit")
	}
	return nil
}
func toCommitSignatureModel(c object.Signature) *models.CommitSignature {
	return &models.CommitSignature{
		Name:  c.Name,
		Email: c.Email,
		Date:  c.When,
	}
}
func toCommitModel(c *object.Commit) *models.Commit {
	p := make([]string, len(c.ParentHashes))
	for i, v := range c.ParentHashes {
		p[i] = v.String()
	}
	return &models.Commit{
		Hash:         c.Hash.String(),
		TreeHash:     c.TreeHash.String(),
		Message:      c.Message,
		ParentHashes: p,

		Author:    toCommitSignatureModel(c.Author),
		Committer: toCommitSignatureModel(c.Committer),
	}
}
