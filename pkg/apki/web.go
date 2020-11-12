package apki

import (
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"

	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
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
	if err := svc.DB.Model(Mirror{}).First(&ent, "host = ?", host).Error; err != nil {
		panic(err)
	}
	_ = res.WriteEntity(ent)
}
func (svc MirrorResource) Mirrors(req *restful.Request, res *restful.Response) {
	var all []MirrorRepresentation
	if err := svc.DB.Model(&Mirror{}).Find(&all).Error; err != nil {
		panic(err)
	}
	_ = res.WriteEntity(all)
}

func (s *IndexerServer) ServeWeb() error {
	c := restful.NewContainer()
	MirrorResource{DB: s.DB}.RegisterTo(c)

	cors := restful.CrossOriginResourceSharing{
		AllowedDomains: []string{"localhost:3000", "alpine.ink"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		Container:      c,
	}
	c.Filter(cors.Filter)
	c.Filter(c.OPTIONSFilter)

	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(http.StripPrefix("/api", c))
	r.HandleFunc("/ping", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("PONG"))
	})
	logrus.Infof("serve %s", s.conf.Web.Addr)
	return http.ListenAndServe(s.conf.Web.Addr, r)
}

type entityJSONAccess struct {
	ContentType string
}

func (e entityJSONAccess) Read(req *restful.Request, v interface{}) error {
	decoder := jsoniter.NewDecoder(req.Request.Body)
	decoder.UseNumber()
	return decoder.Decode(v)
}

func (e entityJSONAccess) Write(resp *restful.Response, status int, v interface{}) error {
	return e.writeJSON(resp, status, e.ContentType, v)
}

func (e entityJSONAccess) writeJSON(resp *restful.Response, status int, contentType string, v interface{}) error {
	if v == nil {
		resp.WriteHeader(status)
		// do not write a nil representation
		return nil
	}
	if true {
		// pretty output must be created and written explicitly
		output, err := jsoniter.MarshalIndent(v, "", " ")
		if err != nil {
			return err
		}
		resp.Header().Set(restful.HEADER_ContentType, contentType)
		resp.WriteHeader(status)
		_, err = resp.Write(output)
		return err
	}
	// not-so-pretty
	resp.Header().Set(restful.HEADER_ContentType, contentType)
	resp.WriteHeader(status)
	return jsoniter.NewEncoder(resp).Encode(v)
}
