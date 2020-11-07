package apki

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IndexerServer struct {
	DB   *gorm.DB
	conf IndexerConf
}
type IndexerConf struct {
	Database DatabaseConf
	Web      struct {
		Addr string
	}
}

func NewServer(conf *IndexerConf) (*IndexerServer, error) {
	if conf.Web.Addr == "" {
		conf.Web.Addr = "0.0.0.0:8080"
	}
	logrus.Debug("connecting db")
	db, err := connectDatabase(conf.Database)
	if err != nil {
		return nil, err
	}
	if conf.Database.AutoMigrate {
		logrus.Info("auto migrate")
		if err := db.AutoMigrate(&Mirror{}, &PackageIndex{}, &Setting{}); err != nil {
			return nil, err
		}
	}
	svr := &IndexerServer{
		DB:   db,
		conf: *conf,
	}
	return svr, nil
}

func (s *IndexerServer) Serve() error {
	r := mux.NewRouter()
	r.Use(recoveryMiddleware)
	r.Use(loggingMiddleware)

	pre := r.PathPrefix("/api/service/me.wener.apkindexer/IndexerService").Methods("POST").Subrouter()
	pre.HandleFunc("/RefreshMirror", func(rw http.ResponseWriter, r *http.Request) {
		if err := s.RefreshMirror(); err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(map[string]interface{}{"message": "OK"})
	})
	pre.HandleFunc("/RefreshIndex", func(rw http.ResponseWriter, r *http.Request) {
		if err := s.RefreshIndex(); err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(map[string]interface{}{"message": "OK"})
	})
	pre.HandleFunc("/IndexCoordinates", func(rw http.ResponseWriter, r *http.Request) {
		if v, err := s.IndexCoordinates(); err != nil {
			panic(err)
		} else {
			rw.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(rw).Encode(v)
		}
	})

	logrus.Infof("serve %s", s.conf.Web.Addr)
	return http.ListenAndServe(s.conf.Web.Addr, r)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				e := json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprint(err)})
				if e != nil {
					logrus.WithError(e).Warn("marshal recovery error failed")
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithField("remote", r.RemoteAddr).Infof("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
