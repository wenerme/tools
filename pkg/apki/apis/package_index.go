package apis

import (
	"bytes"
	"net/http"
	"sort"
	"strings"

	"github.com/wenerme/tools/pkg/apki/convert"
	"github.com/wenerme/tools/pkg/apki/structs"

	"github.com/wenerme/tools/pkg/apki/models"

	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PackageIndexResource struct {
	DB *gorm.DB
}

func (svc PackageIndexResource) RegisterTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/packages").
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/names.txt").To(svc.Names).Produces("text/plain"))
	ws.Route(ws.GET("/origins.txt").To(svc.Origins).Produces("text/plain"))
	ws.Route(ws.GET("/{packageName}").To(svc.GetResource))

	container.Add(ws)
}
func (svc PackageIndexResource) GetResource(req *restful.Request, res *restful.Response) {
	name := req.PathParameter("packageName")
	var all []*models.PackageIndex
	throwError(svc.DB.Where("name = ?", name).Order("build_time desc").Find(&all).Error, "load packages failed")
	var r []*structs.PackageIndex
	conv := &convert.Context{}
	for _, v := range all {
		s, err := convert.ToPackageIndex(conv, v)
		r = append(r, s)
		throwError(err, "convert struct")
	}
	swallowError(res.WriteEntity(r), "write entity")
}

func (svc PackageIndexResource) Names(req *restful.Request, res *restful.Response) {
	var mod []struct {
		Name string
	}
	if err := svc.DB.Model(&models.PackageIndex{}).Distinct("name").Select("name").Order("name").Scan(&mod).Error; err != nil {
		logrus.WithError(err).Error("load names failed")
		panic("failed load names")
	}
	var r []string
	for _, v := range mod {
		r = append(r, v.Name)
	}
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(strings.Join(r, "\n")))
}

func (svc PackageIndexResource) AllOrigins() ([][]string, error) {
	var r [][]string

	var mod []struct {
		Name   string
		Origin string
	}
	if err := svc.DB.Model(&models.PackageIndex{}).Distinct("name").Select("name, origin").Order("name").Scan(&mod).Error; err != nil {
		return nil, err
	}
	m := make(map[string][]string)
	for _, v := range mod {
		o := m[v.Origin]
		if len(o) == 0 {
			o = []string{v.Origin}
		}

		if v.Name != v.Origin {
			o = append(o, v.Name)
		}
		m[v.Origin] = o
	}
	for _, v := range m {
		r = append(r, v)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i][0] < r[j][0]
	})
	return r, nil
}
func (svc PackageIndexResource) Origins(req *restful.Request, res *restful.Response) {
	var r [][]string
	if v, ok := _cache.Get("origins.txt"); ok {
		r = v.([][]string)
	} else {
		var err error
		r, err = svc.AllOrigins()
		if err != nil {
			logrus.WithError(err).Error("load origins failed")
			panic("failed load origins")
		}
		_cache.SetDefault("origins.txt", r)
	}

	var buf bytes.Buffer
	for _, v := range r {
		buf.WriteString(strings.Join(v, ","))
		buf.WriteByte('\n')
	}

	res.WriteHeader(http.StatusOK)
	_, _ = buf.WriteTo(res)
}
