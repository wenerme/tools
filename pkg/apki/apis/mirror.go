package apis

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"github.com/wenerme/tools/pkg/apki/models"
	"gorm.io/gorm"
)

type MirrorResource struct {
	DB *gorm.DB
}

func (svc MirrorResource) RegisterTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/mirrors").
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(svc.Mirrors).Doc("All mirrors").Writes([]MirrorRepresentation{}))
	ws.Route(ws.GET("/{mirror-host}").To(svc.Mirror))

	container.Add(ws)
}
func (svc MirrorResource) Mirror(req *restful.Request, res *restful.Response) {
	host := req.PathParameter("mirror-host")
	var ent MirrorRepresentation
	if err := svc.DB.Model(models.Mirror{}).First(&ent, "host = ?", host).Error; err != nil {
		logrus.WithError(err).Error("load mirror failed")
		panic("failed load data")
	}
	_ = res.WriteEntity(ent)
}

func (svc MirrorResource) Mirrors(req *restful.Request, res *restful.Response) {
	var all []MirrorRepresentation
	if v, ok := _cache.Get("mirrors"); ok {
		all = v.([]MirrorRepresentation)
	} else if err := svc.DB.Model(&models.Mirror{}).Find(&all).Error; err != nil {
		logrus.WithError(err).Error("load mirrors failed")
		panic("failed load data")
	}

	_ = res.WriteEntity(all)
}
