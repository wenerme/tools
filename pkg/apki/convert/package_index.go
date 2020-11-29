package convert

import (
	"github.com/wenerme/tools/pkg/apki/models"
	"github.com/wenerme/tools/pkg/apki/structs"
)

var PackageIndexCommit = Option("PackageIndexCommit")

func ToPackageIndex(ctx *Context, m *models.PackageIndex) (r *structs.PackageIndex, err error) {
	if m == nil {
		return
	}
	r = &structs.PackageIndex{
		Branch:         m.Branch,
		Repo:           m.Repo,
		Arch:           m.Arch,
		Name:           m.Name,
		Version:        m.Version,
		Size:           m.Size,
		InstallSize:    m.InstallSize,
		Description:    m.Description,
		URL:            m.URL,
		License:        m.License,
		Origin:         m.Origin,
		BuildTime:      m.BuildTime,
		Commit:         m.Commit,
		Key:            m.Key,
		Path:           m.Path,
		Depends:        m.Depends,
		Provides:       m.Provides,
		InstallIf:      m.InstallIf,
		MaintainerName: m.MaintainerName,
	}
	if ctx.Has(PackageIndexCommit) {
		v, err := m.GetCommitData()
		if err != nil {
			return nil, err
		}
		r.CommitData, err = ToCommit(ctx, v)
		if err != nil {
			return nil, err
		}
	}
	return
}
