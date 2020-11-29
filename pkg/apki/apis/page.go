package apis

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/wenerme/tools/pkg/apki/convert"
	"github.com/wenerme/tools/pkg/apki/models"
	"github.com/wenerme/tools/pkg/apki/structs"
	"gorm.io/gorm"
)

type PackagePageData struct {
	Name     string                  `json:"name"`
	Packages []*structs.PackageIndex `json:"packages"`
}
type PageResource struct {
	DB *gorm.DB
}

func (svc PageResource) RegisterTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/pages").
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/packages/{packageName}").To(svc.GetPackagePageData))

	container.Add(ws)
}

func (svc PageResource) GetPackagePageData(req *restful.Request, res *restful.Response) {
	name := req.PathParameter("packageName")
	data := &PackagePageData{}
	var all []*models.PackageIndex
	throwError(svc.DB.Where("name = ?", name).Order("build_time desc").Find(&all).Error, "load packages failed")
	var r []*structs.PackageIndex
	conv := &convert.Context{}
	conv = conv.With(convert.PackageIndexCommit)
	for _, v := range all {
		s, err := convert.ToPackageIndex(conv, v)
		r = append(r, s)
		throwError(err, "convert struct")
	}
	data.Packages = r
	swallowError(res.WriteEntity(data), "write entity")
}
