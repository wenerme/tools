package convert

import (
	"github.com/wenerme/tools/pkg/apki/models"
	"github.com/wenerme/tools/pkg/apki/structs"
)

func ToCommit(ctx *Context, m *models.Commit) (r *structs.Commit, err error) {
	if m == nil {
		return nil, nil
	}
	r = &structs.Commit{
		Hash:         m.Hash,
		Message:      m.Message,
		ParentHashes: m.ParentHashes,

		Author:    m.Author,
		Committer: m.Committer,
	}

	return
}
