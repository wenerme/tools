package structs

import (
	"github.com/wenerme/tools/pkg/apki/models"
)

type Commit struct {
	Hash         string   `json:"hash"`
	Message      string   `json:"message"`
	ParentHashes []string `json:"parentHashes"`

	Author    *models.CommitSignature `json:"author"`
	Committer *models.CommitSignature `json:"committer"`
}
